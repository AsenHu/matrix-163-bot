package netease

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"matrix-163-bot/internal/config"

	"github.com/rs/zerolog/log"

	"github.com/XiaoMengXinX/Music163Api-Go/api"
	"github.com/XiaoMengXinX/Music163Api-Go/utils"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

type Song struct {
	ID          int             // GetSongInfo 生成 (网易云侧准备)
	Name        string          // GetSongInfo 生成 (网易云侧准备)
	Artist      []string        // GetSongInfo 生成 (网易云侧准备)
	Album       string          // GetSongInfo 生成 (网易云侧准备)
	SendContent sendMusiContent // !! 二阶段 Matrix 侧读操作要用，警惕隔壁网易的 Music 的写操作 !!
	Music       *bytes.Buffer   // DownloadSong 生成 !!!! 二阶段网易写操作 !!!!
	MusicMxc    id.ContentURI   // GenSongMxc 生成 (Matrix 侧准备)
	Mutex       sync.RWMutex
}

type sendMusiContent struct {
	Body     string   `json:"body"`     // EarlyDownload 生成 (网易云侧准备)
	Filename string   `json:"filename"` // EarlyDownload 生成 (网易云侧准备)
	Info     SongInfo `json:"info"`
	MsgType  string   `json:"msgtype"` // GetSongInfo 生成 (网易云侧准备)
	URL      string   `json:"url"`     // GenSongMxc 生成 (Matrix 侧准备)
}

type SongInfo struct {
	Duration int    `json:"duration"` // GetSongInfo 生成 (网易云侧准备)
	MimeType string `json:"mimetype"` // EarlyDownload 生成 (网易云侧准备)
	Size     int    `json:"size"`     // EarlyDownload 生成 (网易云侧准备)
}

type EarlyDownloadInfo struct {
	DownloadURL string
	MD5         string
}

// -------- 准备信息阶段 --------

// ---- 网易云侧准备

func (s *Song) GetSongInfoByName(name string, cfg config.Netease) (err error) {
	data := utils.RequestData{
		Cookies: []*http.Cookie{
			{
				Name:  "MUSIC_U",
				Value: cfg.MUSIC_U,
			},
		},
	}
	// 搜索歌曲
	searchConfig := api.SearchSongConfig{
		Keyword: name,
		Limit:   1,
	}
	searchResult, err := api.SearchSong(data, searchConfig)
	if err != nil {
		return
	}

	// 检查搜索结果是否为空
	if len(searchResult.Result.Songs) == 0 {
		err = fmt.Errorf("error: search result is empty")
		return
	}

	// 整理艺术家信息
	artists := make([]string, len(searchResult.Result.Songs[0].Artists))
	for i, artist := range searchResult.Result.Songs[0].Artists {
		artists[i] = artist.Name
	}
	// 整理歌曲信息
	s.Mutex.Lock()
	s.ID = searchResult.Result.Songs[0].Id
	s.Name = searchResult.Result.Songs[0].Name
	s.Artist = artists
	s.Album = searchResult.Result.Songs[0].Album.Name
	s.SendContent.Info.Duration = searchResult.Result.Songs[0].Duration
	s.SendContent.MsgType = "m.audio"
	s.Mutex.Unlock()
	return
}

func (s *Song) GetSongInfoById(id int, cfg config.Netease) (err error) {
	data := utils.RequestData{
		Cookies: []*http.Cookie{
			{
				Name:  "MUSIC_U",
				Value: cfg.MUSIC_U,
			},
		},
	}
	// 获取歌曲信息
	songInfo, err := api.GetSongDetail(data, []int{id})
	if err != nil {
		return
	}

	// 整理艺术家信息
	artists := make([]string, len(songInfo.Songs[0].Ar))
	for i, artist := range songInfo.Songs[0].Ar {
		artists[i] = artist.Name
	}
	// 整理歌曲信息
	s.Mutex.Lock()
	s.ID = songInfo.Songs[0].Id
	s.Name = songInfo.Songs[0].Name
	s.Artist = artists
	s.Album = songInfo.Songs[0].Al.Name
	s.SendContent.Info.Duration = songInfo.Songs[0].Dt
	s.SendContent.MsgType = "m.audio"
	s.Mutex.Unlock()
	return
}

func (song *Song) EarlyDownload(earlyDownloadInfo chan EarlyDownloadInfo, cfg config.Netease) (err error) {
	data := utils.RequestData{
		Cookies: []*http.Cookie{
			{
				Name:  "MUSIC_U",
				Value: cfg.MUSIC_U,
			},
		},
	}
	// 检查 song id 是否为空
	song.Mutex.RLock()
	if song.ID == 0 {
		err = fmt.Errorf("error: song id is empty")
		song.Mutex.RUnlock()
		return
	}

	// 准备下载配置
	songCfg := api.SongURLConfig{
		Ids:   []int{song.ID},
		Level: cfg.Quailty,
	}
	song.Mutex.RUnlock()
	// 获取歌曲下载信息
	downloadInfo, err := api.GetSongURL(data, songCfg)
	if err != nil {
		return
	}

	// 检查歌曲下载链接是否为空
	if downloadInfo.Data[0].Url == "" {
		err = fmt.Errorf("error: download url is empty, NetEase may not have the copyright to this song")
		return
	}

	// 生成歌曲提示文本
	// 歌名 作者 专辑 码率(千比特每秒)
	song.Mutex.Lock()
	artists := strings.Join(song.Artist, ", ")
	text := fmt.Sprintf("%s\n%s\n%s\n%d Kilobit per sec",
		song.Name, artists, song.Album, (downloadInfo.Data[0].Br+500)/1000) // 四舍五入到整数
	song.SendContent.Body = text

	// 生成歌曲文件名
	song.SendContent.Filename = fmt.Sprintf("%s - %s.%s", song.Name, artists, downloadInfo.Data[0].Type)

	// 设置 mime type
	// 网易云的 Content-Type 是不可靠的
	if downloadInfo.Data[0].Type == "mp3" {
		song.SendContent.Info.MimeType = "audio/mpeg"
	} else if downloadInfo.Data[0].Type == "flac" {
		song.SendContent.Info.MimeType = "audio/flac"
	} else {
		song.SendContent.Info.MimeType = "audio/%s" + downloadInfo.Data[0].Type
		log.Warn().Msg(fmt.Sprintf("unknown mime type: %s", downloadInfo.Data[0].Type))
		log.Warn().Msg("please report this issue to the developer")
	}

	// 设置大小
	song.SendContent.Info.Size = downloadInfo.Data[0].Size
	song.Mutex.Unlock()

	// 生成下载信息
	var earlyInfo EarlyDownloadInfo
	earlyInfo.DownloadURL = downloadInfo.Data[0].Url
	if cfg.CheckMd5 {
		earlyInfo.MD5 = downloadInfo.Data[0].Md5
	}
	earlyDownloadInfo <- earlyInfo
	close(earlyDownloadInfo)
	return
}

// ---- Matrix 侧准备

func (s *Song) GenSongMxc(client *mautrix.Client) (err error) {
	inline := func(s *Song) (err error) { // 为了方便 defer
		s.Mutex.RLock()
		defer s.Mutex.RUnlock()
		// 检查歌曲是否已存在
		if s.MusicMxc.String() != "" {
			s.SendContent.URL = s.MusicMxc.String()
			return
		}
		if s.SendContent.URL != "" {
			err = fmt.Errorf("why you delete ContentURI but keep URL?")
			return
		}
		return
	}
	err = inline(s)
	if err != nil {
		return err
	}

	// 生成 mxc
	resp, err := client.CreateMXC(context.Background())
	if err != nil {
		return
	}

	// 设置 mxc
	s.Mutex.Lock()
	s.MusicMxc = resp.ContentURI
	s.SendContent.URL = resp.ContentURI.String()
	log.Info().Msg(fmt.Sprintf("Media ID: %s", s.SendContent.URL))
	s.Mutex.Unlock()

	return
}

// -------- 上传下载 和 发送阶段 --------

// ---- 网易云 --> Matrix 下载上传

func (s *Song) DownloadSong(downloadInfo EarlyDownloadInfo) (err error) { // 这里的写入可能影响到隔壁 Matrix 侧的读取，这是这个对象最后一次写入，所以这个函数之后的向 Matrix 上传不需要读锁，但 Matrix 侧发送消息需要读锁，因为**发消息和下载是并发的，而下载和上传是串行的**
	// 检查歌曲下载链接是否为空
	if downloadInfo.DownloadURL == "" {
		err = fmt.Errorf("error: download url is empty, NetEase may not have the copyright to this song")
		return
	}
	// 下载歌曲
	resp, err := http.Get(downloadInfo.DownloadURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	// 同时计算 md5
	var tee io.Reader = resp.Body
	md5Calc := md5.New()
	if downloadInfo.MD5 != "" {
		tee = io.TeeReader(resp.Body, md5Calc)
	}
	// 读取歌曲
	s.Mutex.Lock()
	s.Music = bytes.NewBuffer(nil)
	_, err = s.Music.ReadFrom(tee)
	s.Mutex.Unlock()
	if err != nil {
		return
	}
	// 检查 md5
	if downloadInfo.MD5 != "" {
		md5sum := fmt.Sprintf("%x", md5Calc.Sum(nil))
		if md5sum == downloadInfo.MD5 {
			log.Info().Msg(fmt.Sprintf("md5 check success: %s", md5sum))
		} else {
			err = fmt.Errorf("md5 check failed: want %s, got %s", downloadInfo.MD5, md5sum)
			s.Mutex.Lock()
			s.Music = nil
			s.Mutex.Unlock()
			return
		}
	} else {
		log.Warn().Msg("md5 check disabled, skip check")
	}
	return
}

func (s *Song) UploadSong(client mautrix.Client) (err error) { // 这个函数执行时，这个对象不会再被写入，所以不需要锁。因此未来修改时，需要检查是否有写入操作，如果有，需要加锁
	// 检查上传信息是否齐全
	/*
		需要
		s.Music
		s.SendContent.Info.Size
		s.SendContent.Info.MimeType
		s.SendContent.Filename
		s.MusicMxc
	*/
	if s.Music == nil || s.SendContent.Info.Size == 0 || s.SendContent.Info.MimeType == "" || s.SendContent.Filename == "" || s.MusicMxc.String() == "" {
		err = fmt.Errorf("error: upload information is not complete")
		return
	}

	// 准备上传数据
	media := mautrix.ReqUploadMedia{
		Content:       s.Music,
		ContentLength: int64(s.SendContent.Info.Size),
		ContentType:   s.SendContent.Info.MimeType,
		FileName:      s.SendContent.Filename,
		MXC:           s.MusicMxc,
	}

	// 上传歌曲
	_, err = client.UploadMedia(context.Background(), media)
	return
}

// ---- Matrix 侧发送

func (s *Song) SendSong(client mautrix.Client, roomId id.RoomID) (err error) { // 这个函数执行时，song 可能正在被写入，所以需要读锁
	// 检查发送信息是否齐全
	/*
		需要 SendContent 全部信息
		s.SendContent.Body
		s.SendContent.Filename
		s.SendContent.Info.Duration
		s.SendContent.Info.MimeType
		s.SendContent.Info.Size
		s.SendContent.MsgType
		s.SendContent.URL
	*/
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()
	if s.SendContent.Body == "" || s.SendContent.Filename == "" || s.SendContent.Info.Duration == 0 || s.SendContent.Info.MimeType == "" || s.SendContent.Info.Size == 0 || s.SendContent.MsgType == "" || s.SendContent.URL == "" {
		marshaledContent, _ := json.Marshal(s)
		log.Error().Msg(fmt.Sprintf("error: send information is not complete: %s", marshaledContent))
		err = fmt.Errorf("error: send information is not complete")
		return
	}

	// 发送消息
	_, err = client.SendMessageEvent(context.Background(), roomId, event.EventMessage, s.SendContent)
	return
}

// -------- 同步上传发送 --------
// 用于支持不支持异步的 Matrix Homeserver

func (s *Song) SyncSendSong(client mautrix.Client, roomId id.RoomID) (err error) {
	// 同步上传 --------
	// 检查上传信息是否齐全
	/*
		需要
		s.Music
		s.SendContent.Info.Size
		s.SendContent.Info.MimeType
		s.SendContent.Filename
	*/
	if s.Music == nil || s.SendContent.Info.Size == 0 || s.SendContent.Info.MimeType == "" || s.SendContent.Filename == "" {
		err = fmt.Errorf("error: upload information is not complete")
		return
	}

	// 准备上传数据
	media := mautrix.ReqUploadMedia{
		Content:       s.Music,
		ContentLength: int64(s.SendContent.Info.Size),
		ContentType:   s.SendContent.Info.MimeType,
		FileName:      s.SendContent.Filename,
	}

	// 上传歌曲
	resp, err := client.UploadMedia(context.Background(), media)
	if err != nil {
		return
	}

	// 生成 mxc
	s.MusicMxc = resp.ContentURI
	s.SendContent.URL = resp.ContentURI.String()

	// 同步发送 --------
	// 检查发送信息是否齐全
	/*
		需要 SendContent 全部信息
		s.SendContent.Body
		s.SendContent.Filename
		s.SendContent.Info.Duration
		s.SendContent.Info.MimeType
		s.SendContent.Info.Size
		s.SendContent.MsgType
		s.SendContent.URL
	*/
	if s.SendContent.Body == "" || s.SendContent.Filename == "" || s.SendContent.Info.Duration == 0 || s.SendContent.Info.MimeType == "" || s.SendContent.Info.Size == 0 || s.SendContent.MsgType == "" || s.SendContent.URL == "" {
		marshaledContent, _ := json.Marshal(s)
		log.Error().Msg(fmt.Sprintf("error: send information is not complete: %s", marshaledContent))
		err = fmt.Errorf("error: send information is not complete")
		return
	}

	// 发送消息
	_, err = client.SendMessageEvent(context.Background(), roomId, event.EventMessage, s.SendContent)
	return
}

/*
func SearchSong(name string) (songs []Song) {
	data := utils.RequestData{
		Cookies: []*http.Cookie{
			{
				Name:  "MUSIC_U",
				Value: "",
			},
		},
	}
	// 搜索配置
	config := api.SearchSongConfig{
		Keyword: name,
		Limit:   5,
	}

	// 搜索歌曲
	searchResult, err := api.SearchSong(data, config)
	if err != nil {
		return
	}

	// 整理搜索结果
	songs = make([]Song, len(searchResult.Result.Songs))
	for i, song := range searchResult.Result.Songs {
		// 整理艺术家信息
		artists := make([]string, len(song.Artists))
		for j, artist := range song.Artists {
			artists[j] = artist.Name
		}
		// 整理歌曲信息
		songs[i] = Song{
			ID:     song.Id,
			Name:   song.Name,
			Artist: artists,
		}
	}
	return
}
*/

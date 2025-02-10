package netease

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/XiaoMengXinX/Music163Api-Go/api"
	"github.com/XiaoMengXinX/Music163Api-Go/utils"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/id"
)

type Song struct {
	ID          int
	Name        string
	Artist      []string
	SendContent sendMusiContent
	MusicByte   []byte
	MusicMxc    id.ContentURI
}

type sendMusiContent struct {
	Body     string   `json:"body"`
	Filename string   `json:"filename"`
	Info     SongInfo `json:"info"`
	MsgType  string   `json:"msgtype"`
	URL      string   `json:"url"`
}

type SongInfo struct {
	Duration int    `json:"duration"`
	MimeType string `json:"mimetype"`
	Size     int    `json:"size"`
}

func GetSongInfoByName(name string) (song Song, err error) {
	data := utils.RequestData{
		Cookies: []*http.Cookie{
			{
				Name:  "MUSIC_U",
				Value: "",
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
	// 整理艺术家信息
	artists := make([]string, len(searchResult.Result.Songs[0].Artists))
	for i, artist := range searchResult.Result.Songs[0].Artists {
		artists[i] = artist.Name
	}
	// 整理歌曲信息
	song = Song{
		ID:     searchResult.Result.Songs[0].Id,
		Name:   searchResult.Result.Songs[0].Name,
		Artist: artists,
	}
	song.Info.Duration = searchResult.Result.Songs[0].Duration
	return
}

func GetSongInfoById(id int) (song Song, err error) {
	data := utils.RequestData{
		Cookies: []*http.Cookie{
			{
				Name:  "MUSIC_U",
				Value: "",
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
	song = Song{
		ID:     songInfo.Songs[0].Id,
		Name:   songInfo.Songs[0].Name,
		Artist: artists,
	}
	song.Info.Duration = songInfo.Songs[0].Dt
	return
}

func GenSongMxc(client *mautrix.Client, song *Song) (err error) {
	data := utils.RequestData{
		Cookies: []*http.Cookie{
			{
				Name:  "MUSIC_U",
				Value: "",
			},
		},
	}
	// 检查 song id 是否为空
	if song.ID == 0 {
		err = fmt.Errorf("error: song id is empty")
		return
	}

	// 获取歌曲下载信息
	downloadInfo, err := api.GetSongDownloadURL(data, song.ID)
	if err != nil {
		return
	}

	// 检查歌曲下载链接是否为空
	if downloadInfo.Data.Url == "" {
		err = fmt.Errorf("error: download url is empty, NetEase may not have the copyright to this song")
		return
	}

	// 准备流式上传歌曲
	// 准备输入流
	resp, err := http.Get(downloadInfo.Data.Url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	// 准备输出流
	media, err := client.UploadMedia(context.Background(), mautrix.ReqUploadMedia{
		Content:       resp.Body,
		ContentLength: resp.ContentLength,
		ContentType:   resp.Header.Get("Content-Type"),
		FileName:      fmt.Sprintf("%s - %s.mp3", song.Name, song.Artist),
	})
	if err != nil {
		return
	}

	// 设置 mime type 和大小
	song.Info.MimeType = resp.Header.Get("Content-Type")
	size, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return err
	}
	song.Info.Size = size

	// 生成 mxc
	song.MusicMxc = media.ContentURI
	return
}

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

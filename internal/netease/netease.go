package netease

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"

	"github.com/XiaoMengXinX/Music163Api-Go/api"
	"github.com/XiaoMengXinX/Music163Api-Go/utils"
	"github.com/rs/zerolog/log"
)

func DownloadSong(song_id int) (reader io.Reader, err error) {
	log.Printf("Downloading song %d", song_id)
	// 获取歌曲下载信息
	download_info, err := api.GetSongDownloadURL(utils.RequestData{}, song_id)
	if err != nil {
		log.Print(err)
		return
	}

	// 下载歌曲
	if download_info.Data.Url == "" {
		err = fmt.Errorf("error: download URL is empty")
		log.Print(err)
		return
	}
	log.Printf("Downloading: %s", download_info.Data.Url)
	resp, err := http.Get(download_info.Data.Url)
	if err != nil {
		log.Print(err)
		return
	}
	// 验证文件
	if download_info.Data.Md5 == "" {
		log.Warn().Msg("md5 is empty, skip verification")
	} else {
		hash := md5.New()
		if _, err := io.Copy(hash, resp.Body); err != nil {
			log.Print(err)
			return nil, err
		}
		md5sum := fmt.Sprintf("%x", hash.Sum(nil))
		if md5sum != download_info.Data.Md5 {
			err = fmt.Errorf("error: md5 mismatch, expected %s, got %s", download_info.Data.Md5, md5sum)
			log.Print(err)
			return
		}
		log.Printf("md5 verification passed: %s", md5sum)
	}

	return resp.Body, nil
}

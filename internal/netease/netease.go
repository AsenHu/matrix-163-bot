package netease

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/XiaoMengXinX/Music163Api-Go/api"
	"github.com/XiaoMengXinX/Music163Api-Go/types"
	"github.com/XiaoMengXinX/Music163Api-Go/utils"
)

func Search_Song(keyword string) (result types.SearchSongData, err error) {
	data := utils.RequestData{}
	Search_Config := api.SearchSongConfig{
		Keyword: keyword,
		Limit:   5,
	}

	Search_Result, err := api.SearchSong(data, Search_Config)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	return Search_Result, nil
}

func Download_Song(song_id int, path string) (full_path string, err error) {
	data := utils.RequestData{}

	download_info, err := api.GetSongDownloadURL(data, song_id)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 下载歌曲
	resp, err := http.Get(download_info.Data.Url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
	err = os.Mkdir("cache", 0755)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// 下载路径是 cache/song_id-md5
	full_path = fmt.Sprintf("cache/%d-%s", song_id, download_info.Data.Md5)
	out, err := os.Create(full_path)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)

	// 验证文件
	file, err := os.Open(full_path)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()
	hash := md5.New()
	if _, err = io.Copy(hash, file); err != nil {
		fmt.Println("Error:", err)
		return
	}
	hashInBytes := hash.Sum(nil)
	md5_str := fmt.Sprintf("%x", hashInBytes)
	if md5_str != download_info.Data.Md5 {
		fmt.Println("Error: MD5 not match.")
		return
	}

	return full_path, nil
}

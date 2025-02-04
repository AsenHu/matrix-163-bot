package netease

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/XiaoMengXinX/Music163Api-Go/api"
	"github.com/XiaoMengXinX/Music163Api-Go/types"
	"github.com/XiaoMengXinX/Music163Api-Go/utils"
)

var data = utils.RequestData{}

func Search_Song(keyword string) (result types.SearchSongData, err error) {
	// 搜索设置
	Search_Config := api.SearchSongConfig{
		Keyword: keyword,
		Limit:   5,
	}

	// 搜索歌曲
	Search_Result, err := api.SearchSong(data, Search_Config)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	return Search_Result, nil
}

func Download_Song(song_id int, path string) (full_path string, err error) {
	// 检查歌曲是否在缓存中
	// 匹配文件名为 song_id-md5
	full_path = fmt.Sprintf("%s/%d-*", path, song_id)
	files, err := filepath.Glob(full_path)
 	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if len(files) > 0 {
		fmt.Println("Song already exists in cache.")
	}
	// 获取歌曲下载信息
	download_info, err := api.GetSongDownloadURL(data, song_id)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 下载歌曲
	if download_info.Data.Url == "" {
		err = fmt.Errorf("error: download URL is empty")
		return
	}
	fmt.Println("Downloading:", download_info.Data.Url)
	resp, err := http.Get(download_info.Data.Url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
	err = os.Mkdir(path, 0755)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// 检查 md5 是否存在
	if download_info.Data.Md5 == "" {
		fmt.Println("Warning: MD5 is empty.")
		download_info.Data.Md5 = "empty"
	}
	// 下载路径是 path/song_id-md5
	full_path = fmt.Sprintf("%s/%d-%s", path, song_id, download_info.Data.Md5)
	out, err := os.Create(full_path)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 验证文件
	if download_info.Data.Md5 == "empty" {
		fmt.Println("Warning: MD5 is empty, skip verification.")
		return full_path, nil
	}
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
		err = fmt.Errorf("error: MD5 is %s, should be %s", md5_str, download_info.Data.Md5)
		return
	}

	return full_path, nil
}

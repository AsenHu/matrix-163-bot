package main

import (
	"fmt"
	"matrix-163-bot/internal/netease"
)

func main() {
	// 搜索歌曲
	search_result, err := netease.Search_Song("warma")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(search_result.Result.Songs[0].Id)
	// 下载歌曲
	full_path, err := netease.Download_Song(search_result.Result.Songs[0].Id, "cache")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Downloaded:", full_path)
}

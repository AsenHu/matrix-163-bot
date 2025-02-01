package main

import (
	"fmt"
	"strings"

	"github.com/XiaoMengXinX/Music163Api-Go/api"
	"github.com/XiaoMengXinX/Music163Api-Go/utils"
)

func main() {
	data := utils.RequestData{}
	Search_Config := api.SearchSongConfig{
		Keyword: "114514",
		Limit:   5,
	}

	Search_Result, err := api.SearchSong(data, Search_Config)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, song := range Search_Result.Result.Songs {
		var artists_list []string
		for _, artist := range song.Artists {
			artists_list = append(artists_list, artist.Name)
		}
		artists := strings.Join(artists_list, ", ")
		fmt.Printf("歌曲名称: %s 鸽手: %s 歌曲 ID: %d\n", song.Name, artists, song.Id)
	}
}

package netease

import (
	"fmt"

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

func Get_Download_URL(song_id int) (result types.SongDownloadURLData, err error) {
	data := utils.RequestData{}

	download_URL, err := api.GetSongDownloadURL(data, song_id)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	return download_URL, nil
}

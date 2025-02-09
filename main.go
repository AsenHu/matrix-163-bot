package main

import (
	"fmt"
	"net/http"

	"github.com/XiaoMengXinX/Music163Api-Go/api"
	"github.com/XiaoMengXinX/Music163Api-Go/utils"
)

func main() {
	data := utils.RequestData{
		Cookies: []*http.Cookie{
			{
				Name:  "MUSIC_U",
				Value: "",
			},
		},
	}
	songCfg := api.SongURLConfig{
		Ids: []int{2629003017},
	}
	result, err := api.GetSongURL(data, songCfg)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(result.RawJson)
}

/*
func main() {
	data := utils.RequestData{
		Cookies: []*http.Cookie{
			{
				Name:  "MUSIC_U",
				Value: "",
			},
		},
	}
	result, err := api.GetLoginStatus(data)
	if err != nil {
		println(err)
	}
	println(result.RawJson)
}

/*
func main() {
	var data = utils.RequestData{}
	result, err := api.GetQrUnikey(data)
	if err != nil {
		println(err)
	}
	println(result.Unikey)
}
*/
/*
import (
	"fmt"
	"strings"

	"matrix-163-bot/internal/config"
	"matrix-163-bot/internal/matrix"
	"matrix-163-bot/internal/netease"

	"github.com/rs/zerolog/log"
)

func main() {
	// 加载配置
	cfg := config.NewConfig("config.json")
	err := cfg.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	// 打印账户信息
	log.Info().Msg("Account info by config")
	cfg.Print()

	// 开始登陆
	log.Info().Msg("Try to login...")
	client, err := matrix.Login(&cfg.Content)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to login")
	}
	log.Info().Msg("Login success")
	cfg.Print()

	songId := 2601622999
	song, err := netease.GetSongInfoById(songId)
	if err != nil {
		fmt.Println(err)
	}
	err = netease.GenSongMxc(client, &song)
	if err != nil {
		fmt.Println(err)
	}
	println(song.Name)
	println(strings.Join(song.Artist, ", "))
	println(song.ID)
	println(song.Info.Duration)
	println(song.Info.MimeType)
	println(song.Info.Size)
	println(song.MusicMxc.String())
}
*/

package main

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

	songName := "泛泛人类不会祈祷 人声本家"
	song, err := netease.GetSongInfoByName(songName)
	if err != nil {
		fmt.Println(err)
		return
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

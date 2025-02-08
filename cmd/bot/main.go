package main

import (
	"matrix-163-bot/internal/config"
	"matrix-163-bot/internal/matrix"

	"github.com/rs/zerolog/log"
)

const CONFIG_PATH = "config.json"

func main() {
	// 加载配置
	cfg := config.NewConfig(CONFIG_PATH)
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

	// 保存配置
	err = cfg.Save()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to save config")
	}

	// 设置回调函数
	matrix.SetSyncer(client)

	// 启动同步
	log.Info().Msg("Start sync...")
	if err := client.Sync(); err != nil {
		log.Fatal().Err(err).Msg("Failed to sync")
	}
}

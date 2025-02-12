package main

import (
	"errors"
	"flag"
	"matrix-163-bot/internal/config"
	"matrix-163-bot/internal/limiter"
	"matrix-163-bot/internal/worker"

	"github.com/rs/zerolog/log"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/id"
)

var CONFIG_PATH string

func init() {
	flag.StringVar(&CONFIG_PATH, "c", "config.json", "Path to the configuration file")
}

func main() {
	flag.Parse()
	// 加载配置
	cfg := config.NewConfig(CONFIG_PATH)
	err := cfg.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	// 准备 client
	client, err := mautrix.NewClient(cfg.Content.Matrix.BaseURL, id.UserID(cfg.Content.Matrix.Username), cfg.Content.Matrix.Token)
	if err != nil {
		return
	}

	// 登陆
	log.Info().Msg("Try to login...")
	if cfg.Content.Matrix.Token == "" {
		if err := login(client, &cfg.Content.Matrix); err != nil {
			log.Fatal().Err(err).Msg("Failed to login")
		}
		log.Info().Msg("Login success")
		cfg.Save()
	}

	for {
		// 设置回调函数
		limiter := worker.Limiter{
			SearchLimiter:   limiter.NewLimiter(cfg.Content.Speed.Search.Rate, cfg.Content.Speed.Search.Burst),
			DownloadLimiter: limiter.NewLimiter(cfg.Content.Speed.Download.Rate, cfg.Content.Speed.Download.Burst),
		}
		worker.Init(client, &cfg, &limiter)

		// 启动同步
		log.Info().Msg("Start sync...")
		if err := client.Sync(); err != nil {
			// 处理错误
			if errors.Is(err, mautrix.MUnknownToken) {
				log.Warn().Msg("Token expired, relogin...")
				client.AccessToken = ""
				if err := login(client, &cfg.Content.Matrix); err != nil {
					log.Fatal().Err(err).Msg("Failed to relogin")
				}
				cfg.Save()
				continue
			}
			if errors.Is(err, mautrix.MInvalidParam) {
				log.Warn().Msg("Username format error, relogin...")
				if err := login(client, &cfg.Content.Matrix); err != nil {
					log.Fatal().Err(err).Msg("Failed to relogin")
				}
				cfg.Save()
				continue
			}
			log.Fatal().Err(err).Msg("Sync failed")
		}
	}
}

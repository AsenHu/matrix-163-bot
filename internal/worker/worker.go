package worker

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"matrix-163-bot/internal/config"
	"matrix-163-bot/internal/limiter"
	"matrix-163-bot/internal/netease"

	"github.com/rs/zerolog/log"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

type Worker struct {
	Cfg             *config.Config
	Client          *mautrix.Client
	SearchLimiter   *limiter.Limiter
	DownloadLimiter *limiter.Limiter
}

func SetCallBack(client *mautrix.Client, cfg *config.Config) {
	worker := Worker{
		Client:          client,
		Cfg:             cfg,
		SearchLimiter:   limiter.NewLimiter(cfg.Content.Speed.Search.Rate, cfg.Content.Speed.Search.Burst),
		DownloadLimiter: limiter.NewLimiter(cfg.Content.Speed.Download.Rate, cfg.Content.Speed.Download.Burst),
	}
	syncer := mautrix.NewDefaultSyncer()
	syncer.OnEventType(event.EventMessage, func(ctx context.Context, ev *event.Event) {
		go func() {
			// log.Info().Msg("Receive message: " + ev.Content.AsMessage().Body)
			err := worker.ProcessMessage(ctx, ev)
			if err != nil {
				log.Error().Err(err).Msg("Failed to process message")
			}
		}()
	})
	worker.Client.Syncer = syncer
}

func (w *Worker) ProcessMessage(ctx context.Context, evt *event.Event) (err error) {
	message := evt.Content.AsMessage()
	if message == nil {
		return
	}
	strSlice := strings.Split(message.Body, " ")
	cmd := strSlice[0]
	song := strings.Join(strSlice[1:], " ")

	switch cmd {
	case "!help":
		log.Info().Msg(fmt.Sprintf("User %s in room %s request help", evt.Sender, evt.RoomID))
		err = w.help(evt.RoomID)
	case "!music":
		log.Info().Msg(fmt.Sprintf("User %s in room %s request music", evt.Sender, evt.RoomID))
		if w.Cfg.Content.Matrix.AsyncUpload {
			err = w.asyncMusic(song, evt.RoomID)
		} else {
			err = w.syncMusic(song, evt.RoomID)
		}
	case "!search":
		log.Info().Msg(fmt.Sprintf("User %s in room %s request search", evt.Sender, evt.RoomID))
		err = w.search(evt.RoomID)
	}
	return
}

func (w Worker) help(roomId id.RoomID) (err error) {
	text := "Commands:\n" +
		"!help: Show help\n" +
		"!music <song name or id>: Request music\n" +
		"!search <song name>: Search music"

	_, err = w.Client.SendText(context.Background(), roomId, text)
	return
}

func (w Worker) syncMusic(input string, roomId id.RoomID) (err error) {
	// 检查 song 是否为空
	if input == "" {
		err = w.help(roomId)
		return
	}

	// 速率限制
	err = w.DownloadLimiter.Take()
	if err != nil {
		return
	}

	var song netease.Song
	earlyDownloadInfo := make(chan netease.EarlyDownloadInfo, 1)

	// GetSongInfo
	var inputId int
	// 检查 song 是否可以转换为 int, 并获取歌曲信息
	if inputId, err = strconv.Atoi(input); err == nil {
		err = song.GetSongInfoById(inputId, w.Cfg.Content.Netease)
	} else {
		err = song.GetSongInfoByName(input, w.Cfg.Content.Netease)
	}
	if err != nil {
		return
	}

	log.Info().Msg("Get song info success")
	// EarlyDownload
	err = song.EarlyDownload(earlyDownloadInfo, w.Cfg.Content.Netease)
	if err != nil {
		return
	}

	log.Info().Msg("Early download success")

	// DownloadSong
	earlyInfo, ok := <-earlyDownloadInfo
	if !ok {
		err = fmt.Errorf("earlyDownloadInfo is closed")
		return
	}
	err = song.DownloadSong(earlyInfo)
	if err != nil {
		return
	}
	log.Info().Msg("Download song success")

	err = song.SyncSendSong(*w.Client, roomId)
	return
}

func (w Worker) asyncMusic(input string, roomId id.RoomID) (err error) {
	// 检查 song 是否为空
	if input == "" {
		err = w.help(roomId)
		return
	}

	// 速率限制
	err = w.DownloadLimiter.Take()
	if err != nil {
		return
	}

	var song netease.Song
	earlyDownloadInfo := make(chan netease.EarlyDownloadInfo, 1)
	startForUploadSong := make(chan bool, 2)
	startForSendSong := make(chan bool, 2)
	// 第一阶段: 获取歌曲信息 和 准备 mxc

	/*
		获取歌曲信息
		1. GetSongInfo
		2. EarlyDownload
	*/
	go func() {
		// GetSongInfo
		var inputId int
		// 检查 song 是否可以转换为 int, 并获取歌曲信息
		if inputId, err = strconv.Atoi(input); err == nil {
			err = song.GetSongInfoById(inputId, w.Cfg.Content.Netease)
		} else {
			err = song.GetSongInfoByName(input, w.Cfg.Content.Netease)
		}
		if err != nil {
			startForUploadSong <- false
			startForSendSong <- false
			log.Error().Err(err).Msg("Failed to get song info")
			return
		}

		log.Info().Msg("Get song info success")
		// EarlyDownload
		err = song.EarlyDownload(earlyDownloadInfo, w.Cfg.Content.Netease)
		if err != nil {
			startForUploadSong <- false
			startForSendSong <- false
			log.Error().Err(err).Msg("Failed to early download")
			return
		}

		log.Info().Msg("Early download success")
		startForUploadSong <- true
		startForSendSong <- true
	}()

	/*
		准备 mxc
		1. GenSongMxc
	*/
	go func() {
		err = song.GenSongMxc(w.Client)
		if err != nil {
			startForUploadSong <- false
			startForSendSong <- false
			if errors.Is(err, mautrix.MUnrecognized) {
				log.Warn().Msg("Your home server may not support AsyncUpload.")
				log.Warn().Msg("Please set matrix.asyncUpload to false in your config file.")
				log.Fatal().Err(err).Msg("Failed to gen song mxc")
			}
			log.Error().Err(err).Msg("Failed to gen song mxc")
			return
		}

		startForUploadSong <- true
		startForSendSong <- true
	}()

	// 一阶段完成，开始二阶段 下载上传歌曲 和 发送消息到 Matrix

	/*
		下载上传歌曲
		1. DownloadSong
		2. UploadSong
	*/
	go func() {
		// DownloadSong
		earlyInfo, ok := <-earlyDownloadInfo
		if !ok {
			err = fmt.Errorf("earlyDownloadInfo is closed")
			log.Error().Err(err).Msg("Failed to get early download info")
			return
		}
		err = song.DownloadSong(earlyInfo)
		if err != nil {
			log.Error().Err(err).Msg("Failed to download song")
			return
		}
		log.Info().Msg("Download song success")

		// UploadSong
		a := <-startForUploadSong
		b := <-startForUploadSong
		if !a || !b {
			log.Error().Msg("Some error happened in the first stage, cancel the upload song stage")
			return
		}
		err = song.UploadSong(*w.Client)
		if err != nil {
			log.Error().Err(err).Msg("Failed to upload song")
			return
		}
		log.Info().Msg("Upload song success")
	}()

	/*
		发送消息到 Matrix
		1. SendSong
	*/
	go func() {
		a := <-startForSendSong
		b := <-startForSendSong
		if !a || !b {
			log.Error().Msg("Some error happened in the first stage, cancel the send song stage")
			return
		}
		err = song.SendSong(*w.Client, roomId)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send song")
			return
		}
		log.Info().Msg("Send song success")
	}()
	log.Info().Msg("All stages are started")
	return
}

func (w Worker) search(roomId id.RoomID) (err error) {
	text := "Search function is not implemented yet"
	_, err = w.Client.SendText(context.Background(), roomId, text)
	return
}

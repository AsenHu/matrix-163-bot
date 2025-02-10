package worker

import (
	"context"
	"strconv"
	"strings"

	"matrix-163-bot/internal/netease"

	"github.com/rs/zerolog/log"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

type Worker struct {
	Client *mautrix.Client
}

type sendMusiContent struct {
	Body     string `json:"body"`
	Filename string `json:"filename"`
	Info     struct {
		Duration int    `json:"duration"`
		MimeType string `json:"mimetype"`
		Size     int    `json:"size"`
	} `json:"info"`
	MsgType string `json:"msgtype"`
	URL     string `json:"url"`
}

func SetCallBack(client *mautrix.Client) {
	worker := Worker{
		Client: client,
	}
	syncer := mautrix.NewDefaultSyncer()
	syncer.OnEventType(event.EventMessage, func(ctx context.Context, ev *event.Event) {
		go func() {
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
		err = w.help(evt.RoomID)
	case "!music":
		err = w.music(song, evt.RoomID)
	case "!search":
		err = w.search(evt.RoomID)
	}
	return
}

func sendMusic(client *mautrix.Client, song netease.Song, roomId id.RoomID) (err error) {
	// 准备 Content
	content := sendMusiContent{
		Body:     song.Name + " - " + strings.Join(song.Artist, ", "),
		Filename: song.Name + ".mp3",
		Info: struct {
			Duration int    `json:"duration"`
			MimeType string `json:"mimetype"`
			Size     int    `json:"size"`
		}{
			Duration: song.Info.Duration,
			MimeType: song.Info.MimeType,
			Size:     song.Info.Size,
		},
		MsgType: "m.audio",
		URL:     song.MusicMxc.String(),
	}

	// 发送消息
	_, err = client.SendMessageEvent(context.Background(), roomId, event.EventMessage, content)
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

func (w Worker) music(input string, roomId id.RoomID) (err error) {
	// 检查 song 是否为空
	if input == "" {
		err = w.help(roomId)
		return
	}

	var song netease.Song
	var inputId int
	// 检查 song 是否可以转换为 int, 并获取歌曲信息
	if inputId, err = strconv.Atoi(input); err == nil {
		song, err = netease.GetSongInfoById(inputId)
	} else {
		song, err = netease.GetSongInfoByName(input)
	}
	if err != nil {
		return
	}

	// 生成歌曲 mxc
	err = netease.GenSongMxc(w.Client, &song)
	if err != nil {
		return
	}

	// 发送歌曲
	err = sendMusic(w.Client, song, roomId)
	return
}

func (w Worker) search(roomId id.RoomID) (err error) {
	text := "Search function is not implemented yet"
	_, err = w.Client.SendText(context.Background(), roomId, text)
	return
}

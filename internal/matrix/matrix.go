package matrix

import (
	"context"
	"matrix-163-bot/internal/config"
	"strings"
	"time"

	"matrix-163-bot/internal/netease"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

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

func Login(cfg *config.Account) (client *mautrix.Client, err error) {
	// 准备登陆
	client, err = mautrix.NewClient(cfg.BaseURL, id.UserID(cfg.Username), cfg.Token)
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	// 开始登陆
	resp, err := client.Login(ctx, &mautrix.ReqLogin{
		Type: "m.login.password",
		Identifier: mautrix.UserIdentifier{
			Type: "m.id.user",
			User: cfg.Username,
		},
		Password:                 cfg.Password,
		Token:                    cfg.Token,
		DeviceID:                 id.DeviceID(cfg.DeviceID),
		InitialDeviceDisplayName: "163 Music Bot",

		StoreCredentials:   true,
		StoreHomeserverURL: true,
	})
	if err != nil {
		return nil, err
	}
	cfg.BaseURL = resp.WellKnown.Homeserver.BaseURL
	cfg.Username = resp.UserID.String()
	cfg.DeviceID = resp.DeviceID.String()
	cfg.Token = resp.AccessToken
	return
}

func SendMusic(client *mautrix.Client, song netease.Song, roomId id.RoomID) (err error) {
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

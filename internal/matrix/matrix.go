package matrix

import (
	"context"
	"matrix-163-bot/internal/config"
	"time"

	"matrix-163-bot/internal/worker"

	"github.com/rs/zerolog/log"
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

func SetSyncer(client *mautrix.Client, worker *worker.Worker) {
	syncer := mautrix.NewDefaultSyncer()
	syncer.OnEventType(event.EventMessage, func(ctx context.Context, ev *event.Event) {
		go func() {
			err := worker.ProcessMessage(ctx, ev)
			if err != nil {
				log.Error().Err(err).Msg("Failed to process message")
			}
		}()
	})
	client.Syncer = syncer
}

/*
func SendMusicWithText(client *mautrix.Client, song Song) (err error) {
	// 准备 Content
	content := sendMusiContent{
		Body:     song.Name - strings.Join(song.Artist, ", "),
		Filename: song.Name + ".mp3",
		
}
		*/
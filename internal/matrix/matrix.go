package matrix

import (
	"context"
	"matrix-163-bot/internal/config"
	"time"

	"github.com/rs/zerolog/log"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

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

func SetSyncer(client *mautrix.Client) {
	syncer := mautrix.NewDefaultSyncer()
	syncer.OnEventType(event.EventMessage, func(ctx context.Context, ev *event.Event) {
		go func() {
			switch {
			case <-ctx.Done():
				return
			default:
			}
			log.Info().Str("Message", ev.Content.AsMessage().Body).Msg("Received Message")
		}()
	})
	client.Syncer = syncer
}

func SendMessage(client *mautrix.Client, roomId string, message string) (err error) {
	_, err = client.SendMessageEvent(roomId, "m.room.message", mautrix.Content{
		MsgType: "m.text",
		Body:    message,
	})
	return
}

func EditMessage(client *mautrix.Client, roomId string, message string, eventId string) (err error) {
	_, err = client.SendMessageEvent(roomId, "m.room.message", mautrix.Content{
		MsgType: "m.text",
		Body:    message,
		Format:  "org.matrix.custom.html",
		RelatesTo: &mautrix.RelatesTo{
			RelType: "m.replace",
			EventID: eventId,
		},
	})
	return
}

func ReplyMessage(client *mautrix.Client, roomId string, message string, eventId string) (err error) {
	_, err = client.SendMessageEvent(roomId, "m.room.message", mautrix.Content{
		MsgType: "m.text",
		Body:    message,
		Format:  "org.matrix.custom.html",
		RelatesTo: &mautrix.RelatesTo{
			RelType: "m.in_reply_to",
			EventID: eventId,
		},
	})
	return
}

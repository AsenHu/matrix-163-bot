package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/chzyer/readline"
	"github.com/rs/zerolog/log"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

func main() {
	baseURL := "matrix.tsubasa.moe"
	username := "neko"
	password := ""

	client, err := mautrix.NewClient(baseURL, "@neko:tsubasa.moe", "")
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.Login(ctx, &mautrix.ReqLogin{
		Type:     "m.login.password",
		Password: password,
		DeviceID: "POGWDWIIDN",
		Identifier: mautrix.UserIdentifier{
			Type: "m.id.user",
			User: username,
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("AccessToken " + resp.AccessToken)
	fmt.Println("DeviceID " + resp.DeviceID)
	fmt.Println("UserID " + resp.UserID)
	fmt.Println("HomeServer " + resp.WellKnown.Homeserver.BaseURL)
	fmt.Println("IdentityServer " + resp.WellKnown.IdentityServer.BaseURL)

	log.Info().Msg("Now running")
	syncCtx, cancelSync := context.WithCancel(context.Background())
	var syncStopWait sync.WaitGroup
	syncStopWait.Add(1)

	// 消息钩子
	syncer := client.Syncer.(*mautrix.DefaultSyncer)
	syncer.OnEventType(event.EventMessage, func(ctx context.Context, ev *event.Event) {
		text := "房间ID: " + ev.RoomID.String() + " 发送者: " + ev.Sender.String() + " 消息ID: " + ev.ID.String() + " 内容: " + ev.Content.AsMessage().Body
		if ev.Sender != "@neko:tsubasa.moe" {
			_, err := client.SendText(context.TODO(), ev.RoomID, text)
			if err != nil {
				fmt.Println(err)
			}
		}
	})

	// 启动同步协程
	go func() {
		err = client.SyncWithContext(syncCtx)
		defer syncStopWait.Done()
		if err != nil && !errors.Is(err, context.Canceled) {
			panic(err)
		}
	}()

	rl, err := readline.New("[no room]> ")
	for {
		rl, _ := rl.Readline()
		if rl == "exit" {
			cancelSync()
			syncStopWait.Wait()
			break
		}
	}
}

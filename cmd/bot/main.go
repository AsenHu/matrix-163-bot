package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

type Account struct {
	Username string `json:"username"`
	Password string `json:"password"`
	BaseURL  string `json:"baseURL"`
	DeviceID string `json:"deviceID"`
	Token    string `json:"token"`
}

const CONFIG_PATH = "config.json"

func main() {
	// 检查本地是否有账户状态
	// 读取文件
	buffer, err := os.ReadFile(CONFIG_PATH)
	if err != nil {
		log.Fatal(err)
	}
	// 解析文件
	var account Account
	if err := json.Unmarshal(buffer, &account); err != nil {
		log.Fatal(err)
	}
	// 打印账户信息
	log.Printf("Username: %s", account.Username)
	log.Printf("Password: %s", account.Password)
	log.Printf("BaseURL: %s", account.BaseURL)
	log.Printf("DeviceID: %s", account.DeviceID)
	log.Printf("Token: %s", account.Token)

	// 开始登陆
	log.Printf("Try to login...")
	client, err := mautrix.NewClient(account.BaseURL, id.UserID(account.Username), account.Token)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	resp, err := client.Login(ctx, &mautrix.ReqLogin{
		Type: "m.login.password",
		Identifier: mautrix.UserIdentifier{
			Type: "m.id.user",
			User: account.Username,
		},
		Password:                 account.Password,
		Token:                    account.Token,
		DeviceID:                 id.DeviceID(account.DeviceID),
		InitialDeviceDisplayName: "163 Music Bot",

		StoreCredentials:   true,
		StoreHomeserverURL: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	account.Username = resp.UserID.String()
	account.DeviceID = resp.DeviceID.String()
	account.Token = resp.AccessToken
	account.BaseURL = resp.WellKnown.Homeserver.BaseURL
	log.Printf("Login success")
	log.Printf("Username: %s", account.Username)
	log.Printf("DeviceID: %s", account.DeviceID)
	log.Printf("Token: %s", account.Token)
	log.Printf("BaseURL: %s", account.BaseURL)
	// 保存账户状态
	buffer, err = json.MarshalIndent(account, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(CONFIG_PATH, buffer, 0600); err != nil {
		log.Fatal(err)
	}
	log.Printf("Save account status success")

	// 设置回调函数
	syncer := mautrix.NewDefaultSyncer()
	syncer.OnEventType(event.EventMessage, func(ctx context.Context, ev *event.Event) {
		go func() {
			switch {
			case <-ctx.Done():
				return
			default:
			}
			log.Printf("Message: %s", ev.Content.AsMessage().Body)
		}()
	})
	client.Syncer = syncer

	// 启动同步
	log.Printf("Start sync...")
	if err := client.Sync(); err != nil {
		log.Fatal(err)
	}
}

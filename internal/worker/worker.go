package worker

import (
	"context"
	//"matrix-163-bot/internal/netease"
	"strings"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

type Worker struct {
	client *mautrix.Client
}


func (w *Worker) ProcessMessage(ctx context.Context, evt *event.Event) (err error) {
	message := evt.Content.AsMessage()
	if message == nil {
		return
	}
	strSlice := strings.Split(message.Body, " ")
	cmd := strSlice[0]
	//song := strings.Join(strSlice[1:], " ")

	switch cmd {
	case "!help":
		//err = w.help()
	case "!music":
		//err = w.music(song)
	case "!search":
		//err = w.search(song)
	}
	return
}

/*
func (w Worker) help() (err error) {
	text := "Commands:\n" +
		"!help: Show help\n" +
		"!music <song name or id>: Request music\n" +
		"!search <song name>: Search music"

	_, err = w.client.SendText(context.Background(), "@163-bot:matrix.org", text)
	return
}

func (w Worker) music(input string) (err error) {
	// 检查 song 是否为空
	if input == "" {
		err = w.help()
		return
	}

	var song Song
	// 检查 song 是否可以转换为 int, 并获取歌曲信息
	if _, err = strconv.Atoi(input); err == nil {
		song = netease.GetSongInfoById(input)
	} else {
		song, err = netease.GetSongInfoByName(input)
	}
	if err != nil {
		return
	}

	// 生成歌曲 mxc
	err = netease.GenSongMxc(w.client, &song)

	// 准备歌曲信息
	text := song.Name + " - " + strings.Join(song.Artist, " ")
	fileName := text + ".mp3"
	*/
package worker

import (
	"context"
	"strings"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"matrix-163-bot/internal/netease"
	"github.com/rs/zerolog/log"
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
	song := strings.Join(strSlice, " ")
	err = errors.New("TODO")
	return
}

func useSongId(client *mautrix.Client, songId int) (message string, err error) {
	// 下载歌曲
	music, err := netease.DownloadSong(songId)
	if err != nil {
		return
	}

	// 上传歌曲
	resp, err := client.UploadMedia(&contaxt.Contaxt{}, mautrix.ReqUploadMedia{
		Content: music,
	})
}

func useSongName(songName string) {
	// do something
}

func searchSong(songName string) {
	// do something
}
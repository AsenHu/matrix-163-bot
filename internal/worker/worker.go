package worker

import (
	"context"
	"matrix-163-bot/internal/config"
	"matrix-163-bot/internal/limiter"
	"strconv"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

type Worker struct {
	WellKnownInfo       WellKnownInfo // DO NOT EDIT THIS STRUCT, EVEN AFTER COPY IN PRIVATE VARIABLE
	GetSongDetail       GetSongDetail
	GetSongDetailByName GetSongDetailByName
	GetSongURL          GetSongURL
	DownloadSong        DownloadSong
	CompressSong        CompressSong
	UploadMedia         UploadMedia
	CreateMXC           CreateMXC
	SendMessageEvent    SendMessageEvent
	Limiter             *Limiter
}

type WellKnownInfo struct {
	SongId     *int
	SearchName *string
	Config     *config.Content
	Client     *mautrix.Client
	Event      *event.Event
}

type Limiter struct {
	SearchLimiter   *limiter.Limiter
	DownloadLimiter *limiter.Limiter
}

type Status struct {
	Mutex  sync.Mutex
	IsDone bool
}

func Init(client *mautrix.Client, cfg *config.Config, limiter *Limiter) {
	syncer := mautrix.NewDefaultSyncer()
	syncer.OnEventType(event.EventMessage, func(ctx context.Context, ev *event.Event) {
		go func() {
			NewWorker(client, cfg, limiter).ProcessMessage(ctx, ev)
		}()
	})
	syncer.OnEventType(event.StateMember, func(ctx context.Context, evt *event.Event) {
		go func() {
			if evt.GetStateKey() == client.UserID.String() && evt.Content.AsMember().Membership == event.MembershipInvite {
				_, err := client.JoinRoomByID(ctx, evt.RoomID)
				if err == nil {
					log.Info().
						Str("room_id", evt.RoomID.String()).
						Str("inviter", evt.Sender.String()).
						Msg("Joined room after invite")
				} else {
					log.Error().Err(err).
						Str("room_id", evt.RoomID.String()).
						Str("inviter", evt.Sender.String()).
						Msg("Failed to join room after invite")
				}
			}
		}()
	})
	client.Syncer = syncer
}

func NewWorker(client *mautrix.Client, cfg *config.Config, limiter *Limiter) *Worker {
	return &Worker{
		WellKnownInfo: WellKnownInfo{
			Config: &cfg.Content,
			Client: client,
		},
		Limiter: limiter,
	}
}

func (w *Worker) ProcessMessage(ctx context.Context, evt *event.Event) (err error) {
	// 速率限制
	err = w.Limiter.DownloadLimiter.Take()
	if err != nil {
		log.Warn().Err(err).Msg("Download rate limit")
		return
	}
	w.WellKnownInfo.Event = evt
	message := evt.Content.AsMessage()
	if message == nil {
		return
	}
	strSlice := strings.Split(message.Body, " ")
	cmd := strSlice[0]
	song := strings.Join(strSlice[1:], " ")

	switch cmd {
	case "!help":
		err = w.help(evt)
	case "!music":
		err = w.music(evt, song)
	case "!search":
		err = w.search(evt)
	}
	return
}

func (w *Worker) help(evt *event.Event) (err error) {
	text := "Commands:\n" +
		"!help: Show help\n" +
		"!music <song name or id>: Request music\n" +
		"!search <song name>: Search music"

	_, err = w.WellKnownInfo.Client.SendText(context.Background(), evt.RoomID, text)
	return
}

func (w *Worker) music(evt *event.Event, input string) (err error) {
	// 检查 song 是否为空
	if input == "" {
		err = w.help(evt)
		return
	}

	// 尝试把 input 转换为数字
	inputId, err := strconv.Atoi(input)
	if err == nil {
		// 说明 input 是数字
		w.WellKnownInfo.SongId = &inputId
	} else {
		// 说明 input 是字符串
		w.WellKnownInfo.SearchName = &input
	}

	log.Info().Msg("Input: " + input)

	// 启动！
	w.StartSendMessageEvent()
	return
}

func (w *Worker) search(evt *event.Event) (err error) {
	text := "Search function is not implemented yet"
	_, err = w.WellKnownInfo.Client.SendText(context.Background(), evt.RoomID, text)
	return
}

package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

type SendMessageEvent struct {
	SendMessageEventReturn
	Status Status
}

type SendMessageEventReturn struct {
	Resp *mautrix.RespSendEvent
	Err  error
}

func (w *Worker) StartSendMessageEvent() (ch chan SendMessageEventReturn) {
	ch = make(chan SendMessageEventReturn, 1)

	go func() {
		w.SendMessageEvent.Status.Mutex.Lock()

		if w.SendMessageEvent.Status.IsDone {
			w.SendMessageEvent.Status.Mutex.Unlock()
			ch <- w.SendMessageEvent.SendMessageEventReturn
			close(ch)
			return
		}

		// 开始准备数据
		// 终于是最后一个对象了
		// 这个对象也是个巨难准备的对象
		var sgsdbn chan GetSongDetailByNameReturn
		var sgsd chan GetSongDetailReturn
		if w.WellKnownInfo.SongId == nil {
			sgsdbn = w.StartGetSongDetailByName()
		} else {
			sgsd = w.StartGetSongDetail()
		}
		var scmxc chan CreateMXCReturn
		var sum chan UploadMediaReturn
		if w.WellKnownInfo.Config.Matrix.AsyncUpload {
			scmxc = w.StartCreateMXC()
		} else {
			sum = w.StartUploadMedia()
		}
		sgsurl := w.StartGetSongURL()

		// 等待所有对象完成
		var duration int
		var songName string
		var artists []string
		var album string
		if w.WellKnownInfo.SongId == nil {
			resp := <-sgsdbn
			if resp.Err != nil {
				w.SendMessageEvent.Err = resp.Err
				w.SendMessageEvent.Status.IsDone = true
				w.SendMessageEvent.Status.Mutex.Unlock()

				ch <- SendMessageEventReturn{
					Err: resp.Err,
				}
				close(ch)
				return
			}
			duration = resp.Resp.Result.Songs[0].Duration
			for _, artist := range resp.Resp.Result.Songs[0].Artists {
				artists = append(artists, artist.Name)
			}
			songName = resp.Resp.Result.Songs[0].Name
			album = resp.Resp.Result.Songs[0].Album.Name
		} else {
			resp := <-sgsd
			if resp.Err != nil {
				w.SendMessageEvent.Err = resp.Err
				w.SendMessageEvent.Status.IsDone = true
				w.SendMessageEvent.Status.Mutex.Unlock()

				ch <- SendMessageEventReturn{
					Err: resp.Err,
				}
				close(ch)
				return
			}
			duration = resp.Resp.Songs[0].Dt
			for _, artist := range resp.Resp.Songs[0].Ar {
				artists = append(artists, artist.Name)
			}
			songName = resp.Resp.Songs[0].Name
			album = resp.Resp.Songs[0].Al.Name
		}
		var songMXC id.ContentURI
		if w.WellKnownInfo.Config.Matrix.AsyncUpload {
			resp := <-scmxc
			if resp.Err != nil {
				w.SendMessageEvent.Err = resp.Err
				w.SendMessageEvent.Status.IsDone = true
				w.SendMessageEvent.Status.Mutex.Unlock()

				ch <- SendMessageEventReturn{
					Err: resp.Err,
				}
				close(ch)
				return
			}
			songMXC = resp.Resp.ContentURI
		} else {
			resp := <-sum
			if resp.Err != nil {
				w.SendMessageEvent.Err = resp.Err
				w.SendMessageEvent.Status.IsDone = true
				w.SendMessageEvent.Status.Mutex.Unlock()

				ch <- SendMessageEventReturn{
					Err: resp.Err,
				}
				close(ch)
				return
			}
			songMXC = resp.Resp.ContentURI
		}
		getURLResp := <-sgsurl
		if getURLResp.Err != nil {
			w.SendMessageEvent.Err = getURLResp.Err
			w.SendMessageEvent.Status.IsDone = true
			w.SendMessageEvent.Status.Mutex.Unlock()

			ch <- SendMessageEventReturn{
				Err: getURLResp.Err,
			}
			close(ch)
			return
		}

		// 准备发送数据
		/* 文本格式
		🎵 **Cream Puff**
		🎤 米虾Fomiki | 💿 Cream Puff
		🔊 Bitrate: 128 kbps
		*/
		text := fmt.Sprintf("🎵 **%s**\n🎤 *%s* | 💿 *%s*", songName, strings.Join(artists, ", "), album)
		if !w.WellKnownInfo.Config.Ffmpeg.Enable {
			text += fmt.Sprintf("\n🔊 **Bitrate:** *%d kbps*", (getURLResp.Resp.Data[0].Br+500)/1000)
		}
		fmtBody := fmt.Sprintf("🎵 <strong>%s</strong><br>🎤 <em>%s</em> | 💿 <em>%s</em>", songName, strings.Join(artists, ", "), album)
		if !w.WellKnownInfo.Config.Ffmpeg.Enable {
			fmtBody += fmt.Sprintf("<br>🔊 <strong>Bitrate:</strong> <em>%d kbps</em>", (getURLResp.Resp.Data[0].Br+500)/1000)
		}

		content := event.MessageEventContent{
			Body: text,
			Mentions: &event.Mentions{
				UserIDs: []id.UserID{
					w.WellKnownInfo.Event.Sender,
				},
			},
			RelatesTo: &event.RelatesTo{
				InReplyTo: &event.InReplyTo{
					EventID: w.WellKnownInfo.Event.ID,
				},
			},
			Format:        "org.matrix.custom.html",
			FormattedBody: fmtBody,
			FileName:      fmt.Sprintf("%s - %s.%s", songName, strings.Join(artists, ", "), getURLResp.Resp.Data[0].Type),
			Info: &event.FileInfo{
				Duration: duration,
			},
			MsgType: "m.audio",
			URL:     id.ContentURIString(songMXC.String()),
		}
		// 确认 Mime Type
		if getURLResp.Resp.Data[0].Type == "mp3" {
			content.Info.MimeType = "audio/mpeg"
		} else if getURLResp.Resp.Data[0].Type == "flac" {
			content.Info.MimeType = "audio/flac"
		} else {
			content.Info.MimeType = "audio/%s" + getURLResp.Resp.Data[0].Type
			log.Warn().Msg(fmt.Sprintf("unknown mime type: %s", getURLResp.Resp.Data[0].Type))
			log.Warn().Msg("please report this issue to the developer")
		}
		// 如果没有压缩，则使用原始长度
		if !w.WellKnownInfo.Config.Ffmpeg.Enable {
			content.Info.Size = getURLResp.Resp.Data[0].Size
		}

		// 发送数据
		json, _ := json.Marshal(content)
		log.Info().Msg(string(json))
		resp, err := w.WellKnownInfo.Client.SendMessageEvent(context.Background(), w.WellKnownInfo.Event.RoomID, event.EventMessage, content)

		// 返回数据
		w.SendMessageEvent.Resp = resp
		w.SendMessageEvent.Err = err
		w.SendMessageEvent.Status.IsDone = true
		w.SendMessageEvent.Status.Mutex.Unlock()

		ch <- SendMessageEventReturn{
			Resp: resp,
			Err:  err,
		}
		close(ch)
	}()
	return
}

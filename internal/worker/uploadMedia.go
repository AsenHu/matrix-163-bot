package worker

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"maunium.net/go/mautrix"
)

type UploadMedia struct {
	UploadMediaReturn
	Status Status
}

type UploadMediaReturn struct {
	Resp mautrix.RespMediaUpload
	Err  error
}

func (w *Worker) StartUploadMedia() (ch chan UploadMediaReturn) {
	ch = make(chan UploadMediaReturn, 1)

	go func() {
		w.UploadMedia.Status.Mutex.Lock()

		if w.UploadMedia.Status.IsDone {
			w.UploadMedia.Status.Mutex.Unlock()
			ch <- w.UploadMedia.UploadMediaReturn
			close(ch)
			return
		}

		// 开始准备数据
		// 这其实是一个很复杂的对象，因为它需要的数据很多
		// 最复杂的东西了

		// 艺术家信息
		var sgsdbn chan GetSongDetailByNameReturn
		var sgsd chan GetSongDetailReturn
		if w.WellKnownInfo.SongId == nil {
			sgsdbn = w.StartGetSongDetailByName()
		} else {
			sgsd = w.StartGetSongDetail()
		}

		// MXC 信息，如果开启了异步上传
		var scmxc chan CreateMXCReturn
		if w.WellKnownInfo.Config.Matrix.AsyncUpload {
			scmxc = w.StartCreateMXC()
		}

		// 歌曲的 Mime Type
		// 明明规范里没说要这个东西，但这个 sdk 里面有，氦
		sgsurl := w.StartGetSongURL()

		// 要上传的歌曲
		scs := w.StartCompressSong()

		// 艺术家
		var songName string
		var artists []string
		if w.WellKnownInfo.SongId == nil {
			resp := <-sgsdbn
			if resp.Err != nil {
				w.UploadMedia.Err = resp.Err
				w.UploadMedia.Status.IsDone = true
				w.UploadMedia.Status.Mutex.Unlock()

				ch <- UploadMediaReturn{
					Err: resp.Err,
				}
				close(ch)
				return
			}
			for _, artist := range resp.Resp.Result.Songs[0].Artists {
				artists = append(artists, artist.Name)
			}
			songName = resp.Resp.Result.Songs[0].Name
		} else {
			resp := <-sgsd
			if resp.Err != nil {
				w.UploadMedia.Err = resp.Err
				w.UploadMedia.Status.IsDone = true
				w.UploadMedia.Status.Mutex.Unlock()

				ch <- UploadMediaReturn{
					Err: resp.Err,
				}
				close(ch)
				return
			}
			for _, artist := range resp.Resp.Songs[0].Ar {
				artists = append(artists, artist.Name)
			}
			songName = resp.Resp.Songs[0].Name
		}

		// MXC 信息
		var mxcResp mautrix.RespCreateMXC
		if w.WellKnownInfo.Config.Matrix.AsyncUpload {
			mxcResp := <-scmxc
			if mxcResp.Err != nil {
				w.UploadMedia.Err = mxcResp.Err
				w.UploadMedia.Status.IsDone = true
				w.UploadMedia.Status.Mutex.Unlock()

				ch <- UploadMediaReturn{
					Err: mxcResp.Err,
				}
				close(ch)
				return
			}
		}

		// Mime Type
		getURLResp := <-sgsurl
		if getURLResp.Err != nil {
			w.UploadMedia.Err = getURLResp.Err
			w.UploadMedia.Status.IsDone = true
			w.UploadMedia.Status.Mutex.Unlock()

			ch <- UploadMediaReturn{
				Err: getURLResp.Err,
			}
			close(ch)
			return
		}

		// 歌曲
		compressSongResp := <-scs // 它是流式的，所以不需要等待
		if compressSongResp.Err != nil {
			w.UploadMedia.Err = compressSongResp.Err
			w.UploadMedia.Status.IsDone = true
			w.UploadMedia.Status.Mutex.Unlock()

			ch <- UploadMediaReturn{
				Err: compressSongResp.Err,
			}
			close(ch)
			return
		}

		// 准备上传数据
		// 上传数据是一个很复杂的过程，因为它需要很多数据
		// 而且网易云它 header 里面的数据是不对的，很讨厌

		req := mautrix.ReqUploadMedia{
			Content:  compressSongResp.Resp,
			FileName: fmt.Sprintf("%s - %s.%s", songName, strings.Join(artists, ", "), getURLResp.Resp.Data[0].Type),
		}
		// 如果没有压缩，则使用原始长度
		if !w.WellKnownInfo.Config.Ffmpeg.Enable {
			req.ContentLength = int64(getURLResp.Resp.Data[0].Size)
		}
		// 确认 Mime Type
		if getURLResp.Resp.Data[0].Type == "mp3" {
			req.ContentType = "audio/mpeg"
		} else if getURLResp.Resp.Data[0].Type == "flac" {
			req.ContentType = "audio/flac"
		} else {
			req.ContentType = "audio/%s" + getURLResp.Resp.Data[0].Type
			log.Warn().Msg(fmt.Sprintf("unknown mime type: %s", getURLResp.Resp.Data[0].Type))
			log.Warn().Msg("please report this issue to the developer")
		}
		// 如果有异步上传，那么就需要 MXC 信息
		if w.WellKnownInfo.Config.Matrix.AsyncUpload {
			req.MXC = mxcResp.ContentURI
		}

		// 开始上传
		resp, err := w.WellKnownInfo.Client.UploadMedia(context.Background(), req)

		// 返回数据
		w.UploadMedia.Resp = *resp
		w.UploadMedia.Err = err
		w.UploadMedia.Status.IsDone = true
		w.UploadMedia.Status.Mutex.Unlock()

		ch <- UploadMediaReturn{
			Resp: *resp,
			Err:  err,
		}
		close(ch)
	}()
	return
}

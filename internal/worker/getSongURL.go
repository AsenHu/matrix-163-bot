package worker

import (
	"net/http"

	"github.com/XiaoMengXinX/Music163Api-Go/api"
	"github.com/XiaoMengXinX/Music163Api-Go/types"
	"github.com/XiaoMengXinX/Music163Api-Go/utils"
	"github.com/rs/zerolog/log"
)

type GetSongURL struct {
	GetSongURLReturn
	Status Status
}

type GetSongURLReturn struct {
	Resp types.SongsURLData
	Err  error
}

func (w *Worker) StartGetSongURL() (ch chan GetSongURLReturn) {
	ch = make(chan GetSongURLReturn, 1)

	go func() {
		w.GetSongURL.Status.Mutex.Lock()

		if w.GetSongURL.Status.IsDone {
			w.GetSongURL.Status.Mutex.Unlock()
			ch <- w.GetSongURL.GetSongURLReturn
			close(ch)
			return
		}

		// 开始准备数据
		// 准备 songId
		var songId int
		if w.WellKnownInfo.SongId == nil {
			resp := <-w.StartGetSongDetailByName()
			if resp.Err != nil {
				w.GetSongURL.Err = resp.Err
				w.GetSongURL.Status.IsDone = true
				w.GetSongURL.Status.Mutex.Unlock()

				ch <- GetSongURLReturn{
					Err: w.GetSongURL.Err,
				}
				close(ch)
				return
			}
			songId = resp.Resp.Result.Songs[0].Id
		} else {
			songId = *w.WellKnownInfo.SongId
		}

		// 检查 MUSIC_U
		if w.WellKnownInfo.Config.Netease.MUSIC_U == "" {
			log.Warn().Msg("MUSIC_U is empty, Some songs may not be able to download")
		}

		// 准备数据
		data := utils.RequestData{
			Cookies: []*http.Cookie{
				{
					Name:  "MUSIC_U",
					Value: w.WellKnownInfo.Config.Netease.MUSIC_U,
				},
			},
		}
		config := api.SongURLConfig{
			Level: w.WellKnownInfo.Config.Netease.Quailty,
			Ids:   []int{songId},
		}

		// 发送请求
		resp, err := api.GetSongURL(data, config)

		// 返回数据
		w.GetSongURL.Resp = resp
		w.GetSongURL.Err = err
		w.GetSongURL.Status.IsDone = true
		w.GetSongURL.Status.Mutex.Unlock()

		ch <- GetSongURLReturn{
			Resp: resp,
			Err:  err,
		}
		close(ch)
	}()
	return
}

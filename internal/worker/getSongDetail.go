package worker

import (
	"net/http"

	"github.com/XiaoMengXinX/Music163Api-Go/api"
	"github.com/XiaoMengXinX/Music163Api-Go/types"
	"github.com/XiaoMengXinX/Music163Api-Go/utils"
	"github.com/rs/zerolog/log"
)

type GetSongDetail struct {
	GetSongDetailReturn
	Status Status
}

type GetSongDetailReturn struct {
	Resp types.SongsDetailData
	Err  error
}

func (w *Worker) StartGetSongDetail() (ch chan GetSongDetailReturn) {
	ch = make(chan GetSongDetailReturn, 1)

	go func() {
		w.GetSongDetail.Status.Mutex.Lock()

		if w.GetSongDetail.Status.IsDone {
			w.GetSongDetail.Status.Mutex.Unlock()
			ch <- w.GetSongDetail.GetSongDetailReturn
			close(ch)
			return
		}

		// 开始准备数据
		// 准备 songId
		var songId int
		if w.WellKnownInfo.SongId == nil {
			resp := <-w.StartGetSongDetailByName()
			if resp.Err != nil {
				// 爆
				w.GetSongDetail.Err = resp.Err
				w.GetSongDetail.Status.IsDone = true
				w.GetSongDetail.Status.Mutex.Unlock()

				ch <- GetSongDetailReturn{
					Err: w.GetSongDetail.Err,
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

		// 获取歌曲详情
		resp, err := api.GetSongDetail(data, []int{songId})

		// 准备返回数据
		w.GetSongDetail.Resp = resp
		w.GetSongDetail.Err = err
		w.GetSongDetail.Status.IsDone = true
		w.GetSongDetail.Status.Mutex.Unlock()

		ch <- GetSongDetailReturn{
			Resp: resp,
			Err:  err,
		}
		close(ch)
	}()
	return
}

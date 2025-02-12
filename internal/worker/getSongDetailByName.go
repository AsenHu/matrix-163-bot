package worker

import (
	"errors"
	"net/http"

	"github.com/XiaoMengXinX/Music163Api-Go/api"
	"github.com/XiaoMengXinX/Music163Api-Go/types"
	"github.com/XiaoMengXinX/Music163Api-Go/utils"
	"github.com/rs/zerolog/log"
)

type GetSongDetailByName struct {
	GetSongDetailByNameReturn
	Status Status
}

type GetSongDetailByNameReturn struct {
	Resp types.SearchSongData
	Err  error
}

func (w *Worker) StartGetSongDetailByName() (ch chan GetSongDetailByNameReturn) {
	ch = make(chan GetSongDetailByNameReturn, 1)

	go func() {
		w.GetSongDetailByName.Status.Mutex.Lock()

		if w.GetSongDetailByName.Status.IsDone {
			w.GetSongDetailByName.Status.Mutex.Unlock()
			ch <- w.GetSongDetailByName.GetSongDetailByNameReturn
			close(ch)
			return
		}

		// 检查数据是否准备好
		if *w.WellKnownInfo.SearchName == "" {
			// 爆
			w.GetSongDetailByName.Err = errors.New("SearchName is empty")
			w.GetSongDetailByName.Status.IsDone = true
			w.GetSongDetailByName.Status.Mutex.Unlock()
			ch <- GetSongDetailByNameReturn{
				Err: w.GetSongDetailByName.Err,
			}
			close(ch)
			return
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
		config := api.SearchSongConfig{
			Keyword: *w.WellKnownInfo.SearchName,
		}

		// 获取歌曲信息
		resp, err := api.SearchSong(data, config)

		// 返回数据
		w.GetSongDetailByName.Resp = resp
		w.GetSongDetailByName.Err = err
		w.GetSongDetailByName.Status.IsDone = true
		w.GetSongDetailByName.Status.Mutex.Unlock()

		ch <- GetSongDetailByNameReturn{
			Resp: resp,
			Err:  err,
		}
	}()
	return
}

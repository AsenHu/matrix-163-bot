package worker

import (
	"crypto/md5"
	"errors"
	"hash"
	"io"
	"net/http"
)

type DownloadSong struct {
	DownloadSongReturn
	Status Status
}

type DownloadSongReturn struct {
	Resp      http.Response
	TeeReader io.Reader
	Md5Calc   hash.Hash
	Err       error
}

func (w *Worker) StartDownloadSong() (ch chan DownloadSongReturn) {
	ch = make(chan DownloadSongReturn, 1)

	go func() {
		w.DownloadSong.Status.Mutex.Lock()

		if w.DownloadSong.Status.IsDone {
			w.DownloadSong.Status.Mutex.Unlock()
			ch <- DownloadSongReturn{
				Err: errors.New("You can't download the song again"),
			}
			close(ch)
			return
		}

		// 因为它不能被执行两次，所以它可以提前解锁，以便其他线程可以直接报错
		w.DownloadSong.Status.IsDone = true
		w.DownloadSong.Status.Mutex.Unlock()

		// 开始准备数据
		// 它的数据是从 GetSongURL 中获取的，不接受来自 WellKnownInfo 的数据
		getSongUrlResp := <-w.StartGetSongURL()
		if getSongUrlResp.Err != nil {
			// 爆
			w.DownloadSong.Err = getSongUrlResp.Err
			ch <- DownloadSongReturn{
				Err: getSongUrlResp.Err,
			}
			close(ch)
			return
		}
		url := getSongUrlResp.Resp.Data[0].Url

		// 开始下载
		resp, err := http.Get(url)

		// 准备 md5 计算器
		var tee io.Reader
		md5Calc := md5.New()
		if w.WellKnownInfo.Config.Netease.CheckMd5 {
			tee = io.TeeReader(resp.Body, md5Calc)
		} else {
			tee = resp.Body
		}

		// 准备返回数据
		w.DownloadSong.Resp = *resp
		w.DownloadSong.TeeReader = tee
		w.DownloadSong.Md5Calc = md5Calc
		w.DownloadSong.Err = err

		// 把数据放入 channel，并关闭 channel
		ch <- DownloadSongReturn{
			Resp:      *resp,
			TeeReader: tee,
			Md5Calc:   md5Calc,
			Err:       err,
		}
		close(ch)
	}()
	return
}

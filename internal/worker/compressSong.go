package worker

import (
	"io"
	"matrix-163-bot/internal/ffmpeg"
)

type CompressSong struct {
	CompressSongReturn
	Status Status
}

type CompressSongReturn struct {
	Resp io.Reader
	Err  error
}

func (w *Worker) StartCompressSong() (ch chan CompressSongReturn) {
	// 准备 channel
	ch = make(chan CompressSongReturn, 1)

	go func() {
		// 锁门
		w.CompressSong.Status.Mutex.Lock()

		// 检查是否运行过
		if w.CompressSong.Status.IsDone {
			w.CompressSong.Status.Mutex.Unlock()
			// 如果运行过，就把数据放入 channel，并关闭 channel
			ch <- w.CompressSong.CompressSongReturn
			close(ch)
			return
		}

		// 开始准备数据
		// 这是一个很特别的函数，因为它就是一个管子
		downSongResp := <-w.StartDownloadSong()
		if downSongResp.Err != nil {
			w.CompressSong.Err = downSongResp.Err
			w.CompressSong.Status.IsDone = true
			w.CompressSong.Status.Mutex.Unlock()

			ch <- CompressSongReturn{
				Err: downSongResp.Err,
			}
			close(ch)
			return
		}
		reader := downSongResp.TeeReader

		// 开始压缩
		// 我就是导管大师！
		resp, err := ffmpeg.Compress(reader, &w.WellKnownInfo.Config.Ffmpeg)

		// 返回数据
		w.CompressSong.Resp = resp
		w.CompressSong.Err = err
		w.CompressSong.Status.IsDone = true
		w.CompressSong.Status.Mutex.Unlock()

		// 把数据放入 channel，并关闭 channel
		ch <- w.CompressSong.CompressSongReturn
		close(ch)
	}()
	return
}

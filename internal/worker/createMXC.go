package worker

import (
	"context"

	"maunium.net/go/mautrix"
)

type CreateMXC struct {
	CreateMXCReturn
	Status Status
}

type CreateMXCReturn struct {
	Resp mautrix.RespCreateMXC
	Err  error
}

func (w *Worker) StartCreateMXC() (ch chan CreateMXCReturn) {
	// 准备 channel
	ch = make(chan CreateMXCReturn, 1)

	go func() {
		// 锁门
		w.CreateMXC.Status.Mutex.Lock()

		// 检查是否运行过
		if w.CreateMXC.Status.IsDone {
			w.CreateMXC.Status.Mutex.Unlock()
			// 如果运行过，就把数据放入 channel，并关闭 channel
			ch <- w.CreateMXC.CreateMXCReturn
			close(ch)
			return
		}

		// 说明没有运行过，那么就开始运行
		// 开始准备数据
		resp, err := w.WellKnownInfo.Client.CreateMXC(context.Background())

		// 准备数据完毕，填写数据，开门
		w.CreateMXC.Resp = *resp
		w.CreateMXC.Err = err
		w.CreateMXC.Status.IsDone = true
		w.CreateMXC.Status.Mutex.Unlock()

		// 把数据放入 channel，并关闭 channel
		ch <- w.CreateMXC.CreateMXCReturn
		close(ch)
	}()
	return
}

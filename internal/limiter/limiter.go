package limiter

import (
	"errors"
	"time"
)

type Limiter struct {
	bucket chan struct{}
}

func NewLimiter(rate int, burst int) *Limiter {
	l := &Limiter{
		bucket: make(chan struct{}, burst),
	}
	go func() {
		for {
			l.bucket <- struct{}{}
			time.Sleep(time.Second / time.Duration(rate))
		}
	}()
	return l
}

func (l *Limiter) Take() (err error) {
	select {
	case <-l.bucket:
	default:
		err = errors.New("error: rate limit exceeded")
	}
	return
}

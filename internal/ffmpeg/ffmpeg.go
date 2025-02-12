package ffmpeg

import (
	"errors"
	"io"
	"matrix-163-bot/internal/config"
)

func Compress(input io.Reader, config *config.Ffmpeg) (output io.Reader, err error) {
	if !config.Enable {
		return input, nil
	}
	err = errors.New("ffmpeg compress not implemented")
	return nil, err
}

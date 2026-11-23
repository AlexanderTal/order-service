package processor

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

func WatchForShutdown(ctx context.Context, wg *sync.WaitGroup, closer io.Closer) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		if err := closer.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close component")
		}
	}()
}

func Wrap(wg *sync.WaitGroup, fn func() error) io.Closer {
	return CloserFunc(func() error {
		wg.Add(1)
		defer wg.Done()
		return fn()
	})
}

type CloserFunc func() error

func (f CloserFunc) Close() error {
	return f()
}

func NewCloserContextFunc(timeout time.Duration, shutdown func(context.Context) error) io.Closer {
	return CloserFunc(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		return shutdown(ctx)
	})
}

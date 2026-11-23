package processor

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type CloserFunc func() error

func (f CloserFunc) Close() error {
	return f()
}

func NewCloserContextFunc(f func(ctx context.Context) error, ctx context.Context, timeout time.Duration) io.Closer {
	return CloserFunc(func() error {
		shutdownCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		return f(shutdownCtx)
	})
}

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

func Wrap(wg *sync.WaitGroup, f func() error) io.Closer {
	return CloserFunc(func() error {
		wg.Add(1)
		defer wg.Done()
		return f()
	})
}

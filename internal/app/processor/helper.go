package processor

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"
)

type Closer interface {
	Close() error
}

type CloserFunc func() error

func (f CloserFunc) Close() error {
	return f()
}

func NewCloserContextFunc(cancel context.CancelFunc) Closer {
	return CloserFunc(func() error {
		cancel()
		return nil
	})
}

func WatchForShutdown(ctx context.Context, wg *sync.WaitGroup, closer Closer) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		if err := closer.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close component")
		}
	}()
}

func Wrap(wg *sync.WaitGroup, fn func() error) Closer {
	return CloserFunc(func() error {
		wg.Add(1)
		defer wg.Done()
		return fn()
	})
}

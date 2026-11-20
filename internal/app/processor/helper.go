package processor

import (
	"context"

	"github.com/rs/zerolog/log"
)

type CloserFunc func() error

func (f CloserFunc) Close() error {
	return f()
}

func NewCloserContextFunc(ctx context.Context, cancel context.CancelFunc) CloserFunc {
	return func() error {
		cancel()
		return nil
	}
}

func WatchForShutdown(ctx context.Context, closer interface{ Close() error }) {
	<-ctx.Done()
	if err := closer.Close(); err != nil {
		log.Error().Err(err).Msg("failed to close component")
	}
}

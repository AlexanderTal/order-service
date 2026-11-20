package builder

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/AlexanderTal/order-service/internal/app/config"
	rhandler "github.com/AlexanderTal/order-service/internal/app/handler/http"
	rhealth "github.com/AlexanderTal/order-service/internal/app/handler/http/health"
	"github.com/AlexanderTal/order-service/internal/app/processor"
	rprocessor "github.com/AlexanderTal/order-service/internal/app/processor/http"
	rcpostgres "github.com/AlexanderTal/order-service/internal/app/repository/conn/postgres"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

type Builder struct {
	cCtx   *cli.Context
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	err    error
	cfg    config.Config

	chErrors chan error

	connPostgres *rcpostgres.Client

	healthHandler rhandler.Health

	processors []processor.Processor
}

func NewBuilder(cCtx *cli.Context) *Builder {
	ctx, cancel := context.WithCancel(context.Background())
	b := &Builder{
		cCtx:          cCtx,
		ctx:           ctx,
		cancel:        cancel,
		chErrors:      make(chan error, 1),
		healthHandler: rhealth.NewHandler(),
	}

	go b.waitForSignal()
	go b.printErrors()

	return b
}

func (b *Builder) BuildConfig() *Builder {
	if b.err != nil {
		return b
	}

	config.Load(config.LoadArgs{
		Output:          b.cCtx.App.Writer,
		EnableSimpleLog: b.cCtx.Bool("no-json"),
	})
	b.cfg = config.Root
	return b
}

func (b *Builder) BuildRepoConnPostgres() *Builder {
	if b.err != nil {
		return b
	}

	b.connPostgres, b.err = rcpostgres.NewClient(b.ctx, b.cfg.Repository.Postgres)
	return b
}

func (b *Builder) BuildProcHttp() *Builder {
	if b.err != nil {
		return b
	}

	b.processors = append(b.processors, rprocessor.NewHTTP(b.healthHandler, b.cfg.Processor.WebServer))
	return b
}

func (b *Builder) Run() error {
	if b.err != nil {
		log.Fatal().Err(b.err).Msg("Failed to initialize application")
		return b.err
	}

	log.Info().Msg("Application initialized")

	b.exec()

	return b.err
}

func (b *Builder) waitForSignal() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-b.chErrors:
		b.err = err
		b.cancel()
	case sig := <-sigCh:
		log.Info().Str("signal", sig.String()).Msg("Shutdown is requested")
		b.cancel()
	case <-b.ctx.Done():
		return
	}
}

func (b *Builder) printErrors() {
	select {
	case err := <-b.chErrors:
		b.err = err
		b.cancel()
	case <-b.ctx.Done():
		return
	}
}

func (b *Builder) exec() {
	for _, p := range b.processors {
		p.StartAsync(b.ctx, &b.wg)
	}

	b.wg.Wait()

	if b.connPostgres != nil {
		if err := b.connPostgres.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close database connection")
		}
	}

	log.Info().Msg("Application completed")
}

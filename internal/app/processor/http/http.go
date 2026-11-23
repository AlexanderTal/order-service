package rprocessor

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/AlexanderTal/order-service/internal/app/config/section"
	rhandler "github.com/AlexanderTal/order-service/internal/app/handler/http"
	"github.com/AlexanderTal/order-service/internal/app/processor"
	"github.com/AlexanderTal/order-service/internal/app/util"
	"github.com/AlexanderTal/order-service/internal/pkg/http/httph"
	"github.com/AlexanderTal/order-service/internal/pkg/http/mzerolog"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type httpProc struct {
	server http.Server
	addr   string
}

func NewHTTP(hHealth rhandler.Health, cfg section.ProcessorWebServer) processor.Processor {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(adaptRequestMiddleware(httph.NewErrorMiddleware()))
	router.Use(mzerolog.NewMiddleware(mzerolog.WithSkipper(util.IsFilteredHttpRoute)))
	router.Use(gin.Recovery())

	router.NoRoute(handleNotFound)
	vGenericRegHealthCheck(router, hHealth)

	for _, route := range router.Routes() {
		log.Info().Str("method", route.Method).Str("path", route.Path).Msg("Route registered")
	}

	p := &httpProc{addr: fmt.Sprintf(":%d", cfg.ListenPort)}
	p.server.Handler = router

	return p
}

func (p *httpProc) StartAsync(ctx context.Context, wg *sync.WaitGroup) {
	lc := net.ListenConfig{}
	l, err := lc.Listen(ctx, "tcp", p.addr)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}

	go p.serve(l)
	log.Info().Str("listen_addr", p.addr).Msg("Listening of TCP addr for HTTP server has been started")

	listenerCloser := processor.CloserFunc(func() error {
		if err := l.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
			return err
		}
		return nil
	})
	processor.WatchForShutdown(ctx, wg, listenerCloser)

	shutdownCloser := processor.NewCloserContextFunc(func(ctx context.Context) error {
		err := p.server.Shutdown(ctx)
		if err != nil && !errors.Is(err, net.ErrClosed) {
			return err
		}
		return nil
	}, ctx, 5*time.Second)
	processor.WatchForShutdown(ctx, wg, shutdownCloser)
}

func (p *httpProc) serve(l net.Listener) {
	_ = p.server.Serve(l)
}

func adaptRequestMiddleware(m httph.Middleware) gin.HandlerFunc {
	return func(c *gin.Context) {
		inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Request = r
			c.Next()
		})
		m(inner).ServeHTTP(c.Writer, c.Request)
	}
}

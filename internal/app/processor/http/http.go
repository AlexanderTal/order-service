package rprocessor

import (
	"context"
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
	addr   string
	server *http.Server
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

	p := httpProc{addr: fmt.Sprintf(":%d", cfg.ListenPort)}
	p.server = &http.Server{
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &p
}

func (p *httpProc) StartAsync(ctx context.Context, wg *sync.WaitGroup) {
	var lc net.ListenConfig
	l, err := lc.Listen(ctx, "tcp", p.addr)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start listening TCP addr")
	}

	log.Info().Str("listen_addr", p.addr).Msg("Listening of TCP addr for HTTP server has been started")

	go p.serve(l)

	processor.WatchForShutdown(ctx, wg, processor.CloserFunc(l.Close))
	processor.WatchForShutdown(ctx, wg, processor.NewCloserContextFunc(
		p.server.Shutdown,
		context.Background(), 5*time.Second,
	))
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

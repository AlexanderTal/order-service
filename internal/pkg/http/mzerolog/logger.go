package mzerolog

import (
	"net/http"
	"strings"
	"time"

	"github.com/AlexanderTal/order-service/internal/pkg/http/httph"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type middleware struct {
	log zerolog.Logger

	fromOptions struct {
		skipper func(r *http.Request) bool
	}
}

func (m *middleware) Callback(c *gin.Context) {
	start := time.Now()

	c.Next()

	err := httph.ErrorGet(c.Request)
	if m.fromOptions.skipper != nil && m.fromOptions.skipper(c.Request) {
		return
	}

	duration := time.Since(start)

	var sb strings.Builder
	sb.Grow(48 + len(c.Request.RequestURI))
	sb.WriteString(c.Request.Method)
	sb.WriteByte(' ')
	sb.WriteString(c.Request.RequestURI)
	if err == nil {
		sb.WriteString(" finished with no error")
	} else {
		sb.WriteString(" finished (or aborted) with error")
	}

	var event *zerolog.Event
	if err != nil {
		event = m.log.Error().Err(err)
	} else {
		event = m.log.Debug()
	}

	event.
		Ctx(c.Request.Context()).
		Dur("exec_time", duration).
		Str("client_ip", c.ClientIP()).
		Msg(sb.String())
}

func NewMiddleware(opts ...Option) gin.HandlerFunc {
	m := &middleware{
		log: log.Logger,
		fromOptions: struct {
			skipper func(r *http.Request) bool
		}{
			skipper: defaultSkipper,
		},
	}

	for _, opt := range opts {
		opt(m)
	}

	return m.Callback
}

func defaultSkipper(_ *http.Request) bool {
	return false
}

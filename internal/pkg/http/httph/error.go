package httph

import (
	"context"
	"net/http"
)

type errorCtxKey struct{}

type Middleware = func(http.Handler) http.Handler

type errorContext struct {
	err        error
	statusCode int
}

func ErrorPrepare(r *http.Request) *http.Request {
	if r.Context().Value(errorCtxKey{}) != nil {
		return r
	}
	ctx := context.WithValue(r.Context(), errorCtxKey{}, &errorContext{})
	return r.WithContext(ctx)
}

func ErrorApply(r *http.Request, err error) {
	ctx, ok := r.Context().Value(errorCtxKey{}).(*errorContext)
	if ok {
		ctx.err = err
	}
}

func ErrorGet(r *http.Request) error {
	ctx, ok := r.Context().Value(errorCtxKey{}).(*errorContext)
	if ok {
		return ctx.err
	}
	return nil
}

func ErrorApplyStatusCode(r *http.Request, code int) {
	ctx, ok := r.Context().Value(errorCtxKey{}).(*errorContext)
	if ok {
		ctx.statusCode = code
	}
}

func ErrorGetStatusCode(r *http.Request) int {
	ctx, ok := r.Context().Value(errorCtxKey{}).(*errorContext)
	if ok && ctx.statusCode != 0 {
		return ctx.statusCode
	}
	return http.StatusOK
}

func NewErrorMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, ErrorPrepare(r))
		})
	}
}

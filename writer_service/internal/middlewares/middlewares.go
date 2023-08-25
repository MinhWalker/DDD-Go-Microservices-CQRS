package middlewares

import (
	"fmt"
	"github.com/minhwalker/cqrs-microservices/core/pkg/logger"
	"github.com/minhwalker/cqrs-microservices/writer_service/config"
	"github.com/valyala/fasthttp"
	"strings"
	"time"
)

const (
	maxHeaderBytes = 1 << 20
	bodyLimit      = 2 * 1024 * 1024 // 2 MB
	readTimeout    = 15 * time.Second
	writeTimeout   = 15 * time.Second
	gzipLevel      = 5
)

type MiddlewareManager interface {
	RequestLoggerMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler
	RecoverMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler
	RequestIDMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler
	GzipMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler
	BodyLimitMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler
	ApplyMiddlewares(handlers ...fasthttp.RequestHandler) fasthttp.RequestHandler
}

type middlewareManager struct {
	log logger.Logger
	cfg *config.Config
}

func NewMiddlewareManager(log logger.Logger, cfg *config.Config) *middlewareManager {
	return &middlewareManager{log: log, cfg: cfg}
}

func (mw *middlewareManager) RequestLoggerMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		// Your logging logic here
		next(ctx)
	}
}

func (mw *middlewareManager) RecoverMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if r := recover(); r != nil {
				// Your recovery logic here
				errMsg := fmt.Sprintf("Panic: %v", r)
				mw.log.Error(errMsg)

				// You can also write a custom error response to ctx here
				ctx.SetStatusCode(fasthttp.StatusInternalServerError)
				ctx.SetBody([]byte("Internal Server Error"))
			}
		}()

		next(ctx)
	}
}

func (mw *middlewareManager) RequestIDMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		// Your request ID logic here
		next(ctx)
	}
}

func (mw *middlewareManager) GzipMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if strings.Contains(string(ctx.RequestURI()), "swagger") {
			next(ctx) // Skip Gzip compression
			return
		}

		// Check if client accepts Gzip
		if strings.Contains(string(ctx.Request.Header.Peek("Accept-Encoding")), "gzip") {
			// Set Gzip content encoding
			ctx.Response.Header.Add("Content-Encoding", "gzip")
			// Apply Gzip compression to response body
			fasthttp.WriteGzip(ctx.Response.BodyWriter(), ctx.Response.Body())
		} else {
			next(ctx) // Proceed without Gzip compression
		}
	}
}

func (mw *middlewareManager) BodyLimitMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		// Your body limit logic here
		limit := bodyLimit
		if ctx.Request.Header.ContentLength() > limit {
			ctx.Response.SetStatusCode(fasthttp.StatusRequestEntityTooLarge)
			ctx.Response.SetBody([]byte("Request Entity Too Large"))
			return
		}
		next(ctx)
	}
}

func (mw *middlewareManager) ApplyMiddlewares(handlers ...fasthttp.RequestHandler) fasthttp.RequestHandler {
	// Chain the provided handlers into a single handler function
	return func(ctx *fasthttp.RequestCtx) {
		for _, handler := range handlers {
			handler(ctx)
		}
	}
}

package server

import (
	"fmt"
	"github.com/minhwalker/cqrs-microservices/docs"
	"github.com/valyala/fasthttp"
)

func (s *server) runHttpServer() error {
	return fasthttp.ListenAndServe(fmt.Sprintf("%s", s.cfg.Http.Port), s.mapRoutes)
}

func (s *server) mapRoutes(ctx *fasthttp.RequestCtx) {
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Title = "API Gateway"
	docs.SwaggerInfo.Description = "API Gateway CQRS microservices."
	docs.SwaggerInfo.BasePath = "/api/v1"

	if string(ctx.Path()) == "/swagger/*" {
		ctx.NotFound()
		return
	}

	//// Apply middleware
	//handler := s.applyMiddlewares(
	//	s.mw.RequestLoggerMiddleware,
	//	s.mw.RecoverMiddleware,
	//	s.mw.GzipMiddleware,
	//	s.mw.BodyLimitMiddleware,
	//)
	//
	//// Call the handler chain
	//handler(ctx)
}

func (s *server) applyMiddlewares(middlewares ...func(fasthttp.RequestHandler) fasthttp.RequestHandler) fasthttp.RequestHandler {
	finalHandler := func(ctx *fasthttp.RequestCtx) {
		// Final handler to be executed
	}

	// Apply middlewares in reverse order
	for i := len(middlewares) - 1; i >= 0; i-- {
		finalHandler = middlewares[i](finalHandler)
	}

	return finalHandler
}

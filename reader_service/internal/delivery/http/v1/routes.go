package v1

import (
	"github.com/valyala/fasthttp"
)

func (h *httpService) MapRoutes() {
	path := "/api/v1/products"
	h.routes.GET(path+"/search", h.SearchProduct())
	h.routes.GET(path+"/get/:id", h.GetProductByID())
	h.routes.OPTIONS(path+"/health", func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.WriteString("OK")
	})
}

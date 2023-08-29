package v1

import (
	"github.com/valyala/fasthttp"
)

func (h *httpService) MapRoutes() {
	path := "/api/v1/products"
	h.routes.POST(path+"", h.CreateProduct())
	h.routes.PUT(path+"/:id", h.UpdateProduct())
	h.routes.DELETE(path+"/:id", h.DeleteProduct())
	h.routes.OPTIONS(path+"/health", func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.WriteString("OK")
	})
}

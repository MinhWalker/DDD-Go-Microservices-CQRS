package v1

import (
	"github.com/valyala/fasthttp"
	"net/http"
)

func (h *httpService) MapRoutes() {
	path := "/api/v1/products"
	h.routes.POST(path+"", h.CreateProduct())
	h.routes.PUT(path+"/:id", h.UpdateProduct())
	h.routes.DELETE(path+"/:id", h.DeleteProduct())
	h.routes.OPTIONS("/health", func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(http.StatusOK)
		ctx.WriteString("OK")
	})
}

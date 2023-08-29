package delivery

import (
	"github.com/valyala/fasthttp"
)

type HttpDeliveryProduct interface {
	GetProductByID() fasthttp.RequestHandler
	SearchProduct() fasthttp.RequestHandler
}

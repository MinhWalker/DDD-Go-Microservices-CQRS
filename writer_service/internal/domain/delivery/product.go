package delivery

import (
	"github.com/valyala/fasthttp"
)

type HttpDeliveryProduct interface {
	CreateProduct() fasthttp.RequestHandler
	UpdateProduct() fasthttp.RequestHandler
	DeleteProduct() fasthttp.RequestHandler
}

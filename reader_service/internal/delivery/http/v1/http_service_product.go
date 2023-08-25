package v1

import (
	"encoding/json"
	"github.com/minhwalker/cqrs-microservices/core/pkg/constants"
	httpErrors "github.com/minhwalker/cqrs-microservices/core/pkg/http_errors"
	"github.com/minhwalker/cqrs-microservices/core/pkg/logger"
	"github.com/minhwalker/cqrs-microservices/core/pkg/tracing"
	"github.com/minhwalker/cqrs-microservices/core/pkg/utils"
	"github.com/minhwalker/cqrs-microservices/reader_service/config"
	"github.com/minhwalker/cqrs-microservices/reader_service/internal/domain/usecase"
	"github.com/minhwalker/cqrs-microservices/reader_service/internal/dto"
	product2 "github.com/minhwalker/cqrs-microservices/reader_service/internal/metrics"
	uuid "github.com/satori/go.uuid"
	"github.com/valyala/fasthttp"
	"net/http"

	"github.com/buaazp/fasthttprouter"
	"github.com/go-playground/validator"
	"github.com/opentracing/opentracing-go"
)

type httpService struct {
	routes  *fasthttprouter.Router
	log     logger.Logger
	cfg     *config.Config
	v       *validator.Validate
	ps      usecase.IProductUsecase
	metrics *product2.ReaderServiceMetrics
}

func NewReaderHttpService(routes *fasthttprouter.Router, log logger.Logger, cfg *config.Config, v *validator.Validate, ps usecase.IProductUsecase, metrics *product2.ReaderServiceMetrics) *httpService {
	return &httpService{routes: routes, log: log, cfg: cfg, v: v, ps: ps, metrics: metrics}
}

// SearchProduct
// @Tags Products
// @Summary Search product
// @Description Get product by name with pagination
// @Accept json
// @Produce json
// @Param search query string false "search text"
// @Param page query string false "page number"
// @Param size query string false "number of elements"
// @Success 200 {object} dto.ProductsListResponse
// @Router /products/search [get]
func (h *httpService) SearchProduct() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		h.metrics.SearchProductHttpRequests.Inc()

		_, span := tracing.StartHttpServerTracerSpan(ctx, "productsHandlers.SearchProduct")
		defer span.Finish()

		pq := utils.NewPaginationFromQueryParams(string(ctx.QueryArgs().Peek(constants.Size)), string(ctx.QueryArgs().Peek(constants.Page)))

		query := dto.NewSearchProductQuery(string(ctx.QueryArgs().Peek(constants.Search)), pq)
		response, err := h.ps.Search(ctx, query)
		if err != nil {
			h.log.WarnMsg("SearchProduct", err)
			h.metrics.ErrorHttpRequests.Inc()
			httpErrors.ErrorCtxResponse(ctx, err, h.cfg.Http.DebugErrorsResponse)
			return
		}

		h.metrics.SuccessHttpRequests.Inc()
		// Marshal the response data to JSON
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			h.log.WarnMsg("JSON Marshal", err)
			httpErrors.ErrorCtxResponse(ctx, err, h.cfg.Http.DebugErrorsResponse)
			return
		}

		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusOK)
		_, _ = ctx.Write(jsonResponse)
	}
}

// GetProductByID
// @Tags Products
// @Summary Get product
// @Description Get product by id
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dto.ProductResponse
// @Router /products/{id} [get]
func (h *httpService) GetProductByID() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		h.metrics.GetProductByIdHttpRequests.Inc()

		_, span := tracing.StartHttpServerTracerSpan(ctx, "productsHandlers.GetProductByID")
		defer span.Finish()

		productUUIDStr := string(ctx.UserValue(constants.ID).(string))
		productUUID, err := uuid.FromString(productUUIDStr)
		if err != nil {
			h.log.WarnMsg("uuid.FromString", err)
			h.traceErr(span, err)
			httpErrors.ErrorCtxResponse(ctx, err, h.cfg.Http.DebugErrorsResponse)
			return
		}

		query := dto.NewGetProductByIdQuery(productUUID)
		response, err := h.ps.GetProductById(ctx, query)
		if err != nil {
			h.log.WarnMsg("GetProductById", err)
			h.metrics.ErrorHttpRequests.Inc()
			httpErrors.ErrorCtxResponse(ctx, err, h.cfg.Http.DebugErrorsResponse)
			return
		}

		h.metrics.SuccessHttpRequests.Inc()

		// Marshal the response data to JSON
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			h.log.WarnMsg("JSON Marshal", err)
			httpErrors.ErrorCtxResponse(ctx, err, h.cfg.Http.DebugErrorsResponse)
			return
		}

		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusOK)
		ctx.Write(jsonResponse)
	}
}

func (h *httpService) traceErr(span opentracing.Span, err error) {
	span.SetTag("error", true)
	span.LogKV("error_code", err.Error())
	h.metrics.ErrorHttpRequests.Inc()
}

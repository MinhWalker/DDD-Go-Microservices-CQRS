package v1

import (
	"encoding/json"

	"github.com/minhwalker/cqrs-microservices/core/pkg/constants"
	httpErrors "github.com/minhwalker/cqrs-microservices/core/pkg/http_errors"
	"github.com/minhwalker/cqrs-microservices/core/pkg/logger"
	"github.com/minhwalker/cqrs-microservices/core/pkg/tracing"
	"github.com/minhwalker/cqrs-microservices/writer_service/config"
	"github.com/minhwalker/cqrs-microservices/writer_service/internal/domain/models"
	"github.com/minhwalker/cqrs-microservices/writer_service/internal/domain/usecase"
	"github.com/minhwalker/cqrs-microservices/writer_service/internal/dto"
	product2 "github.com/minhwalker/cqrs-microservices/writer_service/internal/metrics"

	"github.com/buaazp/fasthttprouter"
	"github.com/go-playground/validator"
	"github.com/opentracing/opentracing-go"
	uuid "github.com/satori/go.uuid"
	"github.com/valyala/fasthttp"
)

type httpService struct {
	routes  *fasthttprouter.Router
	log     logger.Logger
	cfg     *config.Config
	v       *validator.Validate
	ps      usecase.IProductUsecase
	metrics *product2.WriterServiceMetrics
}

func NewWriterHttpService(routes *fasthttprouter.Router, log logger.Logger, cfg *config.Config, v *validator.Validate, ps usecase.IProductUsecase, metrics *product2.WriterServiceMetrics) *httpService {
	return &httpService{routes: routes, log: log, cfg: cfg, v: v, ps: ps, metrics: metrics}
}

func (h *httpService) CreateProduct() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		h.metrics.CreateProductHttpRequests.Inc()

		tracingCtx, span := tracing.StartHttpServerTracerSpan(ctx, "productsHandlers.CreateProduct")
		defer span.Finish()

		var createDto models.CreateProductRequest
		if err := json.Unmarshal(ctx.Request.Body(), &createDto); err != nil {
			h.log.WarnMsg("Bind", err)
			h.traceErr(span, err)
			httpErrors.ErrorCtxResponse(ctx, err, h.cfg.Http.DebugErrorsResponse)
			return
		}

		createDto.ProductID = uuid.NewV4() // Assuming uuid.NewV4() generates a UUID
		if err := h.v.StructCtx(tracingCtx, &createDto); err != nil {
			h.log.WarnMsg("validate", err)
			h.traceErr(span, err)
			httpErrors.ErrorCtxResponse(ctx, err, h.cfg.Http.DebugErrorsResponse)
			return
		}

		if err := h.ps.Create(tracingCtx, dto.NewCreateProductCommand(createDto.ProductID, createDto.Name, createDto.Description, createDto.Price)); err != nil {
			h.log.WarnMsg("CreateProduct", err)
			h.metrics.ErrorHttpRequests.Inc()
			httpErrors.ErrorCtxResponse(ctx, err, h.cfg.Http.DebugErrorsResponse)
			return
		}

		h.metrics.SuccessHttpRequests.Inc()
		ctx.SetStatusCode(fasthttp.StatusCreated)
		responseDto := models.CreateProductResponseDto{ProductID: createDto.ProductID}
		responseJSON, _ := json.Marshal(responseDto)
		ctx.SetContentType("application/json")
		ctx.Write(responseJSON)
	}
}

func (h *httpService) UpdateProduct() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		h.metrics.UpdateProductHttpRequests.Inc()

		tracingCtx, span := tracing.StartHttpServerTracerSpan(ctx, "productsHandlers.UpdateProduct")
		defer span.Finish()

		productUUID, err := uuid.FromString(string(ctx.UserValue(constants.ID).(string)))
		if err != nil {
			h.log.WarnMsg("uuid.FromString", err)
			h.traceErr(span, err)
			httpErrors.ErrorCtxResponse(ctx, err, h.cfg.Http.DebugErrorsResponse)
			return
		}

		var updateDto models.UpdateProductRequest
		if err := json.Unmarshal(ctx.Request.Body(), &updateDto); err != nil {
			h.log.WarnMsg("Bind", err)
			h.traceErr(span, err)
			httpErrors.ErrorCtxResponse(ctx, err, h.cfg.Http.DebugErrorsResponse)
			return
		}
		updateDto.ProductID = productUUID

		if err := h.v.StructCtx(tracingCtx, &updateDto); err != nil {
			h.log.WarnMsg("validate", err)
			h.traceErr(span, err)
			httpErrors.ErrorCtxResponse(ctx, err, h.cfg.Http.DebugErrorsResponse)
			return
		}

		if err := h.ps.Update(tracingCtx, dto.NewUpdateProductCommand(updateDto.ProductID, updateDto.Name, updateDto.Description, updateDto.Price)); err != nil {
			h.log.WarnMsg("UpdateProduct", err)
			h.metrics.ErrorHttpRequests.Inc()
			httpErrors.ErrorCtxResponse(ctx, err, h.cfg.Http.DebugErrorsResponse)
			return
		}

		h.metrics.SuccessHttpRequests.Inc()
		ctx.SetStatusCode(fasthttp.StatusOK)
		responseJSON, _ := json.Marshal(updateDto)
		ctx.SetContentType("application/json")
		ctx.Write(responseJSON)
	}
}

func (h *httpService) DeleteProduct() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		h.metrics.DeleteProductHttpRequests.Inc()

		tracingCtx, span := tracing.StartHttpServerTracerSpan(ctx, "productsHandlers.DeleteProduct")
		defer span.Finish()

		productUUID, err := uuid.FromString(string(ctx.UserValue(constants.ID).(string)))
		if err != nil {
			h.log.WarnMsg("uuid.FromString", err)
			h.traceErr(span, err)
			httpErrors.ErrorCtxResponse(ctx, err, h.cfg.Http.DebugErrorsResponse)
			return
		}

		if err := h.ps.Delete(tracingCtx, dto.NewDeleteProductCommand(productUUID)); err != nil {
			h.log.WarnMsg("DeleteProduct", err)
			h.metrics.ErrorHttpRequests.Inc()
			httpErrors.ErrorCtxResponse(ctx, err, h.cfg.Http.DebugErrorsResponse)
			return
		}

		h.metrics.SuccessHttpRequests.Inc()
		ctx.SetStatusCode(fasthttp.StatusOK)
	}
}

func (h *httpService) traceErr(span opentracing.Span, err error) {
	span.SetTag("error", true)
	span.LogKV("error_code", err.Error())
	h.metrics.ErrorHttpRequests.Inc()
}

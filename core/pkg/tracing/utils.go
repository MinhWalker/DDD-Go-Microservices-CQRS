package tracing

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/segmentio/kafka-go"
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc/metadata"
)

type FastHTTPRequestCarrier struct {
	ctx *fasthttp.RequestCtx
}

func (c *FastHTTPRequestCarrier) Set(key, val string) {
	c.ctx.Request.Header.Set(key, val)
}

func (c *FastHTTPRequestCarrier) ForeachKey(handler func(key, val string) error) error {
	c.ctx.Request.Header.VisitAll(func(k, v []byte) {
		if err := handler(string(k), string(v)); err != nil {
			// Handle error if needed
		}
	})
	return nil
}

func StartHttpServerTracerSpan(ctx *fasthttp.RequestCtx, operationName string) (context.Context, opentracing.Span) {
	carrier := &FastHTTPRequestCarrier{ctx: ctx}

	spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap, carrier)
	if err != nil {
		serverSpan := opentracing.GlobalTracer().StartSpan(operationName)
		ctxWithSpan := opentracing.ContextWithSpan(ctx, serverSpan)
		return ctxWithSpan, serverSpan
	}

	serverSpan := opentracing.GlobalTracer().StartSpan(operationName, ext.RPCServerOption(spanCtx))
	ctxWithSpan := opentracing.ContextWithSpan(ctx, serverSpan)

	return ctxWithSpan, serverSpan
}

func GetTextMapCarrierFromMetaData(ctx context.Context) opentracing.TextMapCarrier {
	metadataMap := make(opentracing.TextMapCarrier)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		for key := range md.Copy() {
			metadataMap.Set(key, md.Get(key)[0])
		}
	}
	return metadataMap
}

func StartGrpcServerTracerSpan(ctx context.Context, operationName string) (context.Context, opentracing.Span) {
	textMapCarrierFromMetaData := GetTextMapCarrierFromMetaData(ctx)

	span, err := opentracing.GlobalTracer().Extract(opentracing.TextMap, textMapCarrierFromMetaData)
	if err != nil {
		serverSpan := opentracing.GlobalTracer().StartSpan(operationName)
		ctx = opentracing.ContextWithSpan(ctx, serverSpan)
		return ctx, serverSpan
	}

	serverSpan := opentracing.GlobalTracer().StartSpan(operationName, ext.RPCServerOption(span))
	ctx = opentracing.ContextWithSpan(ctx, serverSpan)

	return ctx, serverSpan
}

func StartKafkaConsumerTracerSpan(ctx context.Context, headers []kafka.Header, operationName string) (context.Context, opentracing.Span) {
	carrierFromKafkaHeaders := TextMapCarrierFromKafkaMessageHeaders(headers)

	spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap, carrierFromKafkaHeaders)
	if err != nil {
		serverSpan := opentracing.GlobalTracer().StartSpan(operationName)
		ctx = opentracing.ContextWithSpan(ctx, serverSpan)
		return ctx, serverSpan
	}

	serverSpan := opentracing.GlobalTracer().StartSpan(operationName, ext.RPCServerOption(spanCtx))
	ctx = opentracing.ContextWithSpan(ctx, serverSpan)

	return ctx, serverSpan
}

func TextMapCarrierToKafkaMessageHeaders(textMap opentracing.TextMapCarrier) []kafka.Header {
	headers := make([]kafka.Header, 0, len(textMap))

	if err := textMap.ForeachKey(func(key, val string) error {
		headers = append(headers, kafka.Header{
			Key:   key,
			Value: []byte(val),
		})
		return nil
	}); err != nil {
		return headers
	}

	return headers
}

func TextMapCarrierFromKafkaMessageHeaders(headers []kafka.Header) opentracing.TextMapCarrier {
	textMap := make(map[string]string, len(headers))
	for _, header := range headers {
		textMap[header.Key] = string(header.Value)
	}
	return opentracing.TextMapCarrier(textMap)
}

func InjectTextMapCarrier(spanCtx opentracing.SpanContext) (opentracing.TextMapCarrier, error) {
	m := make(opentracing.TextMapCarrier)
	if err := opentracing.GlobalTracer().Inject(spanCtx, opentracing.TextMap, m); err != nil {
		return nil, err
	}
	return m, nil
}

func InjectTextMapCarrierToGrpcMetaData(ctx context.Context, spanCtx opentracing.SpanContext) context.Context {
	if textMapCarrier, err := InjectTextMapCarrier(spanCtx); err == nil {
		md := metadata.New(textMapCarrier)
		ctx = metadata.NewOutgoingContext(ctx, md)
	}
	return ctx
}

func GetKafkaTracingHeadersFromSpanCtx(spanCtx opentracing.SpanContext) []kafka.Header {
	textMapCarrier, err := InjectTextMapCarrier(spanCtx)
	if err != nil {
		return []kafka.Header{}
	}

	kafkaMessageHeaders := TextMapCarrierToKafkaMessageHeaders(textMapCarrier)
	return kafkaMessageHeaders
}

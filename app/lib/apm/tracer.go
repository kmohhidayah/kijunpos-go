package apm

import (
	"context"
	"fmt"
	"github/kijunpos/app/lib/logger"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type Option struct {
	ServiceName  string
	CollectorURL string
	ApiKey       string
	Environment  string
	Insecure     bool
}

// InitTracer will set global otel with our config...
func InitTracer(ctx context.Context, opt Option) func(context.Context) error {
	exporter, err := GetTraceExporterHttp(ctx, opt)
	if err != nil {
		panic(err)
	}
	resources, err := getResource(ctx, opt)
	if err != nil {
		panic(err)
	}

	otel.SetTracerProvider(
		sdkTrace.NewTracerProvider(
			sdkTrace.WithSampler(sdkTrace.AlwaysSample()),
			sdkTrace.WithBatcher(exporter),
			sdkTrace.WithResource(resources),
		),
	)
	return exporter.Shutdown
}

func GetTraceExporterHttp(ctx context.Context, opt Option) (*otlptrace.Exporter, error) {
	secureOption := otlptracehttp.WithInsecure()
	return otlptrace.New(
		ctx,
		otlptracehttp.NewClient(
			secureOption,
			otlptracehttp.WithEndpoint(opt.CollectorURL),
			otlptracehttp.WithHeaders(
				map[string]string{
					"api-key": opt.ApiKey,
				},
			),
		),
	)
}

func getResource(ctx context.Context, opt Option) (*resource.Resource, error) {
	return resource.New(
		ctx,
		resource.WithAttributes(
			attribute.String("library.language", "go"),
			attribute.String("service.name", opt.ServiceName),
			attribute.String("service.environment", opt.Environment),
		),
	)
}

// init global tracer
var otelTracer = otel.Tracer("github/kijunpos")

// GetTracer get
func GetTracer() trace.Tracer {
	return otelTracer
}

// GetTraceIDByCtx get trace id from context
func GetTraceIDByCtx(ctx context.Context) (traceID string) {
	if span := trace.SpanFromContext(ctx); span != nil {
		traceID = span.SpanContext().TraceID().String()
	}

	return traceID
}

func TraceError(ctx context.Context, err error) {
	logger.GetLogger().Error(err)
	if span := trace.SpanFromContext(ctx); span != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(attribute.String("err", err.Error()))
	}
}

func TraceErrorf(ctx context.Context, format string, path string, err error) {
	logger.GetLogger().Errorf(format, path, err)
	if span := trace.SpanFromContext(ctx); span != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(attribute.String("err", fmt.Sprintf(format, path, err)))
	}
}

func TraceWarnf(ctx context.Context, format string, path string, err error) {
	logger.GetLogger().Warnf(format, path, err)
	if span := trace.SpanFromContext(ctx); span != nil {
		span.SetAttributes(attribute.String("warn", fmt.Sprintf(format, path, err)))
	}
}

func TraceInfof(ctx context.Context, format string, path string, err error) {
	logger.GetLogger().Infof(format, path, err)
	if span := trace.SpanFromContext(ctx); span != nil {
		span.SetAttributes(attribute.String("info", fmt.Sprintf(format, path, err)))
	}
}

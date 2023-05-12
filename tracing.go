package otelxorm

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"go.opentelemetry.io/otel/trace"
	"xorm.io/xorm"
	"xorm.io/xorm/contexts"
)

const (
	tracerName = "github.com/jenbonzhang/otelxorm"
)

type OpenTelemetryHook struct {
	config *config
}

func Hook(opts ...Option) contexts.Hook {
	cfg := &config{}
	for _, opt := range opts {
		opt.apply(cfg)
	}
	if cfg.tracerProvider == nil {
		cfg.tracerProvider = otel.GetTracerProvider()
	}

	return &OpenTelemetryHook{
		config: cfg,
	}
}

func WrapEngine(e *xorm.Engine, opts ...Option) {
	e.AddHook(Hook(opts...))
}

func WrapEngineGroup(eg *xorm.EngineGroup, opts ...Option) {
	eg.AddHook(Hook(opts...))
}

func (h *OpenTelemetryHook) BeforeProcess(c *contexts.ContextHook) (context.Context, error) {
	ctx, _ := h.config.tracer.Start(c.Ctx,
		"xorm-db",
		trace.WithSpanKind(trace.SpanKindClient),
	)
	return ctx, nil
}

func (h *OpenTelemetryHook) AfterProcess(c *contexts.ContextHook) error {
	span := trace.SpanFromContext(c.Ctx)
	attrs := make([]attribute.KeyValue, 0)
	defer span.End()

	attrs = append(attrs, h.config.attrs...)
	attrs = append(attrs, attribute.Key("go.orm").String("xorm"))
	attrs = append(attrs, semconv.DBStatement(fmt.Sprintf("%v %v", c.SQL, c.Args)))

	if c.Err != nil {
		span.RecordError(c.Err)
		span.SetStatus(codes.Error, c.Err.Error())
	}
	span.SetAttributes(attrs...)
	return nil
}

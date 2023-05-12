package otelxorm

import (
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"go.opentelemetry.io/otel/trace"
)

// Option specifies instrumentation configuration options.
type Option interface {
	apply(*config)
}
type optionFunc func(*config)

func (o optionFunc) apply(c *config) {
	o(c)
}

type config struct {
	tracerProvider trace.TracerProvider
	tracer         trace.Tracer
	attrs          []attribute.KeyValue
}

func newConfig() *config {
	return &config{
		tracerProvider: trace.NewNoopTracerProvider(),
		tracer: trace.NewNoopTracerProvider().Tracer(
			tracerName,
			trace.WithInstrumentationVersion(SemVersion()),
		),
	}
}

// WithTracerProvider with tracer provider.
func WithTracerProvider(provider trace.TracerProvider) Option {
	return optionFunc(func(cfg *config) {
		cfg.tracerProvider = provider
	})
}

// WithDBSystem configures a db.system attribute. You should prefer using
// WithAttributes and semconv, for example, `otelsql.WithAttributes(semconv.DBSystemSqlite)`.
func WithDBSystem(system string) Option {
	return optionFunc(func(c *config) {
		c.attrs = append(c.attrs, semconv.DBSystemKey.String(system))
	})
}

// WithDBName configures a db.name attribute.
func WithDBName(name string) Option {
	return optionFunc(func(c *config) {
		c.attrs = append(c.attrs, semconv.DBName(name))
	})
}

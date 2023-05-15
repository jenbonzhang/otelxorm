package otelxorm

import (
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"go.opentelemetry.io/otel/trace"
	"reflect"
	"strings"
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
	dbName         string
	tracerProvider trace.TracerProvider
	tracer         trace.Tracer
	attrs          []attribute.KeyValue
	formatSQL      func(sql string, args []interface{}) string
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

func WithFormatSQL(formatSQL func(sql string, args []interface{}) string) Option {
	return optionFunc(func(c *config) {
		c.formatSQL = formatSQL
	})
}

func WithFormatSQLReplace() Option {
	return WithFormatSQL(formatSQLReplace)
}

func defaultFormatSQL(sql string, args []interface{}) string {
	argsStr := fmt.Sprintf("%v", args)
	m, err := json.Marshal(args)
	if err == nil {
		argsStr = string(m)
	}
	return fmt.Sprintf("%v %v", sql, argsStr)
}

func formatSQLReplace(sql string, args []interface{}) string {
	for i, arg := range args {
		if reflect.TypeOf(arg).Kind() == reflect.Ptr {
			sql = strings.Replace(sql, fmt.Sprintf("$%d", i+1), fmt.Sprintf("'%v'", reflect.ValueOf(arg).Elem().Interface()), -1)
		} else {
			sql = strings.Replace(sql, fmt.Sprintf("$%d", i+1), fmt.Sprintf("'%v'", arg), -1)
		}
	}
	return sql
}

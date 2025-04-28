package telemetry

import (
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracerOnce sync.Once
	tracer     trace.Tracer
)

func Tracer() trace.Tracer {

	tracerOnce.Do(func() {
		tracer = otel.Tracer(schemaName)
	})

	return tracer
}

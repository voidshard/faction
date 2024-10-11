package log

import (
	"context"
	"fmt"
	"reflect"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Span struct {
	root trace.Span
}

func (s *Span) End() {
	s.root.End()
}

func (s *Span) SetAttributes(attrs ...map[string]interface{}) {
	s.root.SetAttributes(toOtelAttrs(attrs...)...)
}

// NewSpan creates a new span with the given name and attributes from the default tracer
func NewSpan(ctx context.Context, name string, attrs ...map[string]interface{}) *Span {
	return newSpan(tracer, ctx, name, attrs...)
}

func newSpan(tr trace.Tracer, ctx context.Context, name string, attrs ...map[string]interface{}) *Span {
	Debug().Str("span", name).Msg("Starting span")
	_, span := tr.Start(ctx, name, trace.WithAttributes(toOtelAttrs(attrs...)...))
	return &Span{root: span}
}

func (s *Span) Err(err error) error {
	if err != nil {
		s.root.RecordError(err)
	}
	return err
}

func toOtelAttrs(attrs ...map[string]interface{}) []attribute.KeyValue {
	otelAttrs := []attribute.KeyValue{}
	for _, attrset := range attrs {
		for k, v := range attrset {
			var val reflect.Value
			if reflect.TypeOf(v).Kind() == reflect.Ptr {
				// if the value is a pointer, dereference it (otherwise we get the pointer address)
				val = reflect.ValueOf(v).Elem()
			} else {
				val = reflect.ValueOf(v)
			}

			err, ok := v.(error)
			if ok {
				// the only "struct" we will accept, should be set with Err() (above)
				otelAttrs = append(otelAttrs, attribute.String(k, err.Error()))
				continue
			}

			// do our best to convert the value to an OTel attribute
			switch val.Kind() {
			case reflect.String:
				otelAttrs = append(otelAttrs, attribute.String(k, val.String()))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				otelAttrs = append(otelAttrs, attribute.Int64(k, val.Int()))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				otelAttrs = append(otelAttrs, attribute.Int64(k, int64(val.Uint())))
			case reflect.Float32, reflect.Float64:
				otelAttrs = append(otelAttrs, attribute.Float64(k, val.Float()))
			case reflect.Bool:
				otelAttrs = append(otelAttrs, attribute.Bool(k, val.Bool()))
			default:
				// reject complex types
				Warn().Str("key", k).Str("kind", fmt.Sprintf("%v", val.Kind())).Msg("Unsupported type for span attribute, dropping")
			}
		}
	}
	return otelAttrs
}

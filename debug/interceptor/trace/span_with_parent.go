package trace

import (
	"github.com/opentracing/opentracing-go"
)

type spanWithParent struct {
	span       opentracing.Span
	parentSpan *spanWithParent
}

func newSpanWithParent(span opentracing.Span, parentSpan *spanWithParent) *spanWithParent {
	return &spanWithParent{
		span:       span,
		parentSpan: parentSpan,
	}
}

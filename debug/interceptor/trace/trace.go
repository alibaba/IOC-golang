package trace

import (
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type trace struct {
	grID           int64
	rootSpan       opentracing.Span
	currentSpan    *spanWithParent
	entranceMethod string
}

func newTraceWithClientSpanContext(grID int64, entranceMethod string, clientSpanContext opentracing.SpanContext) *trace {
	rootSpan := GetGlobalTracer().StartSpan(entranceMethod, ext.RPCServerOption(clientSpanContext))
	return &trace{
		grID:           grID,
		currentSpan:    newSpanWithParent(rootSpan, nil),
		rootSpan:       rootSpan,
		entranceMethod: entranceMethod,
	}
}

func newTrace(grID int64, entranceMethod string) *trace {
	rootSpan := GetGlobalTracer().StartSpan(entranceMethod)
	return &trace{
		grID:           grID,
		currentSpan:    newSpanWithParent(rootSpan, nil),
		rootSpan:       rootSpan,
		entranceMethod: entranceMethod,
	}
}

func (t *trace) addChildSpan(name string) {
	func1Span := GetGlobalTracer().StartSpan(name, opentracing.ChildOf(t.currentSpan.span.Context()), opentracing.StartTime(time.Now()))
	innerChildSpan := newSpanWithParent(func1Span, t.currentSpan)
	t.currentSpan = innerChildSpan
}

func (t *trace) returnSpan() {
	t.currentSpan.span.Finish()
	t.currentSpan = t.currentSpan.parentSpan
}

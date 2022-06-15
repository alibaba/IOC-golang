package trace

import (
	"sync"

	"github.com/opentracing/opentracing-go"

	"github.com/alibaba/ioc-golang/debug/api/ioc_golang/debug"
	"github.com/alibaba/ioc-golang/debug/interceptor/common"
)

type Context struct {
	SDID         string
	MethodName   string
	Ch           chan *debug.TraceResponse
	FieldMatcher *common.FieldMatcher

	// trace all gr
	tracesMap sync.Map // goroutine-id -> *trace

	// trace single gr
	targetGR          bool
	grID              int64
	ClientSpanContext opentracing.SpanContext
}

func (t *Context) finish(grID int64) {
	t.tracesMap.Delete(grID)
}

func (t *Context) getTrace(grID int64) *trace {
	val, ok := t.tracesMap.Load(grID)
	if !ok {
		// todo handle not ok
		return nil
	}
	return val.(*trace)
}

func (t *Context) createTrace(grID int64, entranceMethod string) {
	if !t.targetGR {
		t.tracesMap.Store(grID, newTrace(grID, entranceMethod))
		return
	}
	t.tracesMap.Store(grID, newTraceWithClientSpanContext(grID, entranceMethod, t.ClientSpanContext))
}

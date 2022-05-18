package debug

import (
	"sort"

	"github.com/alibaba/ioc-golang/debug/interceptor"
)

var registeredInterceptors = make(interceptor.PriorityInterceptors, 0)

type InterceptorPriorityImpl struct {
	interceptor.Interceptor
	priority int
}

func (i *InterceptorPriorityImpl) Priority() int {
	return i.priority
}

func init() {
	RegisterInterceptor(interceptor.GetWatchInterceptor())
	RegisterInterceptor(interceptor.GetEditInterceptor())
}

func RegisterInterceptor(inter interceptor.Interceptor) {
	if priorityInterceptor, ok := inter.(interceptor.PriorityInterceptor); ok {
		RegisterPriorityInterceptor(priorityInterceptor)
		return
	}
	RegisterPriorityInterceptor(&InterceptorPriorityImpl{
		priority:    0,
		Interceptor: inter,
	})
}

func RegisterPriorityInterceptor(priorityInterceptor interceptor.PriorityInterceptor) {
	registeredInterceptors = append(registeredInterceptors, priorityInterceptor)
	sort.Sort(registeredInterceptors)
}

package common

import "runtime"

const (
	ProxyMethod = "github.com/alibaba/ioc-golang/debug.makeProxyFunction.func1"
)

func CurrentCallingMethodName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(4, pc)
	return runtime.FuncForPC(pc[0]).Name()
}

func TraceLevel(entranceName string) int64 {
	pc := make([]uintptr, 100)
	n := runtime.Callers(0, pc)
	foundEntrance := false
	level := int64(0)

	for i := n - 1; i >= 0; i-- {
		fName := runtime.FuncForPC(pc[i]).Name()
		if foundEntrance {
			if fName == ProxyMethod {
				level++
			}
			continue
		}
		if fName == entranceName {
			foundEntrance = true
		}
	}

	return level - 1
}

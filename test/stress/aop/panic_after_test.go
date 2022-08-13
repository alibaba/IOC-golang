package aop

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang"
	"github.com/alibaba/ioc-golang/aop"
)

const Name = "afterIsCalledAfterPanicTestInterceptor"

func init() {
	// register a custom aop interceptor
	aop.RegisterAOP(aop.AOP{
		Name: Name,
		InterceptorFactory: func() aop.Interceptor {
			i, _ := GetafterIsCalledAfterPanicTestInterceptorIOCInterfaceSingleton()
			return i
		},
	})
}

func TestAOPAfterWithPanic(t *testing.T) {
	assert.Nil(t, ioc.Load())
	panicApp, err := GetPanicAfterCalledTestAppIOCInterfaceSingleton()
	assert.Nil(t, err)
	result := panicApp.RunWithPanic("test")
	assert.Equal(t, "test", result)

	i, _ := GetafterIsCalledAfterPanicTestInterceptorIOCInterfaceSingleton()
	// the invocation have two layers, the internal layer is called after panic
	// assert the internal layer is not skipped because of panic.
	assert.Equal(t, 2, i.GetAfterIsCalledNum())
}

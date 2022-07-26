package main

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang"
	"github.com/alibaba/ioc-golang/test/iocli_command"
)

func TestObservabilityList(t *testing.T) {
	assert.Nil(t, ioc.Load())
	app, err := GetAppIOCInterfaceSingleton()
	assert.Nil(t, err)
	go func() {
		app.Run()
	}()
	time.Sleep(time.Millisecond * 500)
	output, err := iocli_command.Run([]string{"list"}, time.Second)
	assert.Nil(t, err)
	assert.Equal(t, `github.com/alibaba/ioc-golang/example/aop/observability.App
[Run]

github.com/alibaba/ioc-golang/example/aop/observability.ServiceImpl1
[GetHelloString]

github.com/alibaba/ioc-golang/example/aop/observability.ServiceImpl2
[GetHelloString]

`, output)
}

func TestObservabilityMonitor(t *testing.T) {
	assert.Nil(t, ioc.Load())
	app, err := GetAppIOCInterfaceSingleton()
	assert.Nil(t, err)
	go func() {
		app.Run()
	}()
	time.Sleep(time.Millisecond * 500)
	output, err := iocli_command.Run([]string{"monitor", "-i", "3"}, time.Second*4)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(output, `github.com/alibaba/ioc-golang/example/aop/observability.ServiceImpl1.GetHelloString()
Total: 1, Success: 1, Fail: 0, AvgRT: `))
	assert.True(t, strings.Contains(output, `us, FailRate: 0.00%
github.com/alibaba/ioc-golang/example/aop/observability.ServiceImpl2.GetHelloString()
Total: 1, Success: 1, Fail: 0, AvgRT: `))
}

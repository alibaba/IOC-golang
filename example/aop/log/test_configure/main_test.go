package test_default

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang"
	"github.com/alibaba/ioc-golang/aop/common"
	"github.com/alibaba/ioc-golang/config"
	"github.com/alibaba/ioc-golang/example/aop/log/app"
	aopLog "github.com/alibaba/ioc-golang/extension/aop/log"
)

func TestLogAOPLevelConfiguration(t *testing.T) {
	/*
		set global logrus level to warning,
		the properties have the same effect as config.yaml:
		```yaml
		ioc-golang:
		  aop:
		    log:
		      level: warning # global logrus log level
		      invocation-aop-log:
		        level: info # invocation aop log level
		```
	*/
	assert.Nil(t, ioc.Load(
		config.AddProperty(common.IOCGolangAOPConfigPrefix+"."+aopLog.Name+".level", "warning"),
		config.AddProperty(common.IOCGolangAOPConfigPrefix+"."+aopLog.Name+".invocation-aop-log.level", "info"),
	))
	_, err := app.GetAppIOCInterfaceSingleton()
	assert.Nil(t, err)
	assert.Equal(t, logrus.WarnLevel, logrus.GetLevel())
	logInterceptorSingleton, _ := aopLog.GetlogInterceptorIOCInterfaceSingleton(nil)
	assert.Equal(t, logrus.InfoLevel, logInterceptorSingleton.GetInvocationCtxLogger().GetLevel())
}

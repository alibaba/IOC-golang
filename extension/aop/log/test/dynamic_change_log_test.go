package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang"
	"github.com/alibaba/ioc-golang/test/iocli_command"
)

func TestABCSetLogLevel(t *testing.T) {
	assert.Equal(t, logrus.InfoLevel, logrus.GetLevel())
	assert.Nil(t, ioc.Load())

	// set log level to error
	resull, err := iocli_command.Run([]string{"call", "singleton", "github.com/alibaba/ioc-golang/extension/aop/log.GlobalLogrusIOCCtxHook", "SetLogLevel", "--params", "[2]"}, time.Second)
	fmt.Println(resull)
	assert.Nil(t, err)
	assert.Equal(t, logrus.ErrorLevel, logrus.GetLevel())

	// set log level to debug
	_, err = iocli_command.Run([]string{"call", "singleton", "github.com/alibaba/ioc-golang/extension/aop/log.GlobalLogrusIOCCtxHook", "SetLogLevel", "--params", "[5]"}, time.Second)
	assert.Nil(t, err)
	assert.Equal(t, logrus.DebugLevel, logrus.GetLevel())
}

package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang"
	"github.com/alibaba/ioc-golang/config"
)

func (a *App) TestRun(t *testing.T) {
	assert.Equal(t, "myValue", a.ConfigValue.Value())
}

func TestSetConfigTypeByOption(t *testing.T) {
	// start
	if err := ioc.Load(config.WithConfigType("yml")); err != nil {
		panic(err)
	}

	app, err := GetAppSingleton()
	if err != nil {
		panic(err)
	}
	app.TestRun(t)
}

func TestSetConfigTypeByEnv(t *testing.T) {
	os.Setenv("IOC_GOLANG_CONFIG_TYPE", "yml")
	if err := ioc.Load(); err != nil {
		panic(err)
	}

	app, err := GetAppSingleton()
	if err != nil {
		panic(err)
	}
	app.TestRun(t)
}

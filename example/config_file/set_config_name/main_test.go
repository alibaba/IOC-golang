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

func TestSetConfigNameByOption(t *testing.T) {
	// start
	if err := ioc.Load(config.WithConfigName("custom_config")); err != nil {
		panic(err)
	}

	app, err := GetAppSingleton()
	if err != nil {
		panic(err)
	}
	app.TestRun(t)
}

func TestSetConfigNameByEnv(t *testing.T) {
	// start
	os.Setenv("IOC_GOALNG_CONFIG_NAME", "custom_config")
	if err := ioc.Load(); err != nil {
		panic(err)
	}

	app, err := GetAppSingleton()
	if err != nil {
		panic(err)
	}
	app.TestRun(t)
}

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

func TestSetConfigSearchPathByOption(t *testing.T) {
	// start
	if err := ioc.Load(config.WithSearchPath("./custom_config_path")); err != nil {
		panic(err)
	}

	app, err := GetAppSingleton()
	if err != nil {
		panic(err)
	}
	app.TestRun(t)
}

func TestSetConfigSearchPathByEnv(t *testing.T) {
	// start
	os.Setenv("IOC_GOLANG_CONFIG_SEARCH_PATH", "./custom_config_path")
	if err := ioc.Load(); err != nil {
		panic(err)
	}

	app, err := GetAppSingleton()
	if err != nil {
		panic(err)
	}
	app.TestRun(t)
}

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

func TestSetConfigActiveProfileByOption(t *testing.T) {
	// start
	if err := ioc.Load(config.WithProfilesActive("dev")); err != nil {
		panic(err)
	}

	app, err := GetAppSingleton()
	if err != nil {
		panic(err)
	}
	app.TestRun(t)
}

func TestSetConfigActiveProfileByEnv(t *testing.T) {
	os.Setenv("IOC_GOLANG_CONFIG_ACTIVE_PROFILE", "dev")
	if err := ioc.Load(); err != nil {
		panic(err)
	}

	app, err := GetAppSingleton()
	if err != nil {
		panic(err)
	}
	app.TestRun(t)
}

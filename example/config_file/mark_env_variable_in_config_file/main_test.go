package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang"
)

func (a *App) TestRun(t *testing.T) {
	assert.Equal(t, "myValue", a.ConfigValue.Value())
	assert.Equal(t, "myEnvValue", a.ConfigValueFromEnv.Value())
}

func TestSetConfigName(t *testing.T) {
	// start
	os.Setenv("MY_CONFIG_ENV_KEY", "myEnvValue")
	if err := ioc.Load(); err != nil {
		panic(err)
	}

	app, err := GetAppSingleton()
	if err != nil {
		panic(err)
	}
	app.TestRun(t)
}

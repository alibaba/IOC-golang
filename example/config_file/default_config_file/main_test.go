package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang"
)

func (a *App) TestRun(t *testing.T) {
	assert.Equal(t, "myValue", a.ConfigValue.Value())
}

func TestSetConfigName(t *testing.T) {
	// start
	if err := ioc.Load(); err != nil {
		panic(err)
	}

	app, err := GetAppSingleton()
	if err != nil {
		panic(err)
	}
	app.TestRun(t)
}

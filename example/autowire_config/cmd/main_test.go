package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang"
	"github.com/alibaba/ioc-golang/autowire/singleton"
)

func TestConfig(t *testing.T) {
	if err := ioc.Load(); err != nil {
		panic(err)
	}
	appInterface, err := singleton.GetImpl("App-App")
	if err != nil {
		panic(err)
	}
	app := appInterface.(*App)
	assert.Equal(t, "stringValue", app.DemoConfigString.Value())
	assert.Equal(t, 123, app.DemoConfigInt.Value())
	assert.Equal(t, "map[key1:value1 key2:value2 key3:value3 obj:map[objkey1:objvalue1 objkey2:objvalue2 objkeyslice:objslicevalue]]", fmt.Sprint(app.DemoConfigMap.Value()))
	assert.Equal(t, "[sliceValue1 sliceValue2 sliceValue3 sliceValue4]", fmt.Sprint(app.DemoConfigSlice.Value()))
}

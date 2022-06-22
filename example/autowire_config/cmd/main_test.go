/*
 * Copyright (c) 2022, Alibaba Group;
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang"
)

func TestConfig(t *testing.T) {
	if err := ioc.Load(); err != nil {
		panic(err)
	}
	app, err := GetAppSingleton()
	if err != nil {
		panic(err)
	}
	assert.Equal(t, "stringValue", app.DemoConfigString.Value())
	assert.Equal(t, 123, app.DemoConfigInt.Value())
	assert.Equal(t, "map[key1:value1 key2:value2 key3:value3 obj:map[objkey1:objvalue1 objkey2:objvalue2 objkeyslice:objslicevalue]]", fmt.Sprint(app.DemoConfigMap.Value()))
	assert.Equal(t, "[sliceValue1 sliceValue2 sliceValue3 sliceValue4]", fmt.Sprint(app.DemoConfigSlice.Value()))
}

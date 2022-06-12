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

package debug

import (
	"github.com/fatih/color"

	"github.com/alibaba/ioc-golang/config"

	"github.com/alibaba/ioc-golang/debug/common"
	"github.com/alibaba/ioc-golang/debug/interceptor"
)

var paramInterceptors = make([]interceptor.Interceptor, 0)
var responseInterceptors = make([]interceptor.Interceptor, 0)

const (
	defaultDebugPort = "1999"
)

func init() {
	paramInterceptors = append(paramInterceptors, interceptor.GetWatchInterceptor())
	paramInterceptors = append(paramInterceptors, interceptor.GetEditInterceptor())

	responseInterceptors = append(responseInterceptors, interceptor.GetWatchInterceptor())
	responseInterceptors = append(responseInterceptors, interceptor.GetEditInterceptor())
}

var debugMetadata = make(map[string]*common.StructMetadata)

func Load() error {
	bootConfig := &Config{}
	_ = config.LoadConfigByPrefix("debug", bootConfig)
	if bootConfig.Port == "" {
		color.Blue("[Debug] Debug port is set to default :%s", defaultDebugPort)
		bootConfig.Port = defaultDebugPort
	}
	if err := interceptor.Start(bootConfig.Port, debugMetadata); err != nil {
		color.Red("[Debug] Start debug server error = %s", err)
		return err
	}
	return nil
}

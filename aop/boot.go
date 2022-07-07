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

package aop

import (
	"github.com/fatih/color"

	"github.com/alibaba/ioc-golang/aop/common"
	"github.com/alibaba/ioc-golang/config"
)

const (
	defaultDebugPort = "1999"
)

var debugMetadata = make(common.AllInterfaceMetadata)

func Load() error {
	debugConfig := &common.Config{}
	_ = config.LoadConfigByPrefix("debug", debugConfig)
	if debugConfig.Port == "" {
		color.Blue("[Debug] Debug port is set to default :%s", defaultDebugPort)
		debugConfig.Port = defaultDebugPort
	}
	if err := start(debugConfig); err != nil {
		color.Red("[Debug] Start debug server error = %s", err)
		return err
	}
	return nil
}

func GetAllInterfaceMetadata() common.AllInterfaceMetadata {
	return debugMetadata
}

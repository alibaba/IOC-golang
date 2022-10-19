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
	"sync"

	"github.com/alibaba/ioc-golang/logger"

	"github.com/alibaba/ioc-golang/aop/common"
	"github.com/alibaba/ioc-golang/config"
)

const (
	defaultDebugPort = "1999"
)

var debugMetadata = make(common.AllInterfaceMetadata)
var debugMetadataLock = sync.Mutex{}

func Load() error {
	aopConfig := &common.Config{}
	_ = config.LoadConfigByPrefix(common.IOCGolangAOPConfigPrefix, aopConfig)
	if aopConfig.DebugServer.Disable {
		logger.Blue("[Debug] Debug server is disabled")
		return nil
	}
	if aopConfig.DebugServer.Port == "" {
		logger.Blue("[Debug] Debug server port is set to default :%s", defaultDebugPort)
		aopConfig.DebugServer.Port = defaultDebugPort
	}
	if err := start(aopConfig); err != nil {
		logger.Red("[Debug] Start debug server error = %s", err)
		return err
	}
	return nil
}

func GetAllInterfaceMetadata() common.AllInterfaceMetadata {
	return debugMetadata
}

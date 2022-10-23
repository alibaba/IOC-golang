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
	// 1.1 check AOP config
	aopConfig := &common.Config{}
	_ = config.LoadConfigByPrefix(common.IOCGolangAOPConfigPrefix, aopConfig)
	if aopConfig.Disable {
		logger.Blue("[AOP] AOP is disabled")
		return nil
	}

	// 1.2 enable AOP
	enableAOP(aopConfig)

	// 2.1 check Debug Server config
	if aopConfig.DebugServer.Disable {
		// aop is enabled but debug server is disabled, just return
		// on such condition, all aop interceptors also works, but debug service would not be registered
		logger.Blue("[AOP] Debug server is disabled")
		return nil
	} else if aopConfig.DebugServer.Port == "" {
		logger.Blue("[AOP] Debug server port is set to default :%s", defaultDebugPort)
		aopConfig.DebugServer.Port = defaultDebugPort
	}

	// 2.2 start debug server
	if err := startDebugServer(aopConfig); err != nil {
		logger.Red("[AOP] Start debug server error = %s", err)
		return err
	}
	return nil
}

func GetAllInterfaceMetadata() common.AllInterfaceMetadata {
	return debugMetadata
}

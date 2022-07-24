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

package autowire

import (
	"fmt"

	"github.com/alibaba/ioc-golang/logger"
)

var allWrapperAutowires = make(map[string]WrapperAutowire)

func printAutowireRegisteredStructDescriptor() {
	for autowireType, aw := range allWrapperAutowires {
		logger.Blue("[Autowire Type] Found registered autowire type %s", autowireType)
		for sdID := range aw.GetAllStructDescriptors() {
			logger.Blue("[Autowire Struct Descriptor] Found type %s registered SD %s", autowireType, sdID)
		}
	}
}

func Load() error {
	// get all autowires
	allWrapperAutowires = GetAllWrapperAutowires()

	printAutowireRegisteredStructDescriptor()

	// autowire all struct that can be entrance
	for _, aw := range allWrapperAutowires {
		for sdID := range aw.GetAllStructDescriptors() {
			if aw.CanBeEntrance() {
				sd := GetStructDescriptor(sdID)
				if sd == nil {
					continue
				}
				_, err := aw.ImplWithoutParam(sdID, !sd.DisableProxy)
				if err != nil {
					return fmt.Errorf("[Autowire] Impl sd %s failed, reason is %s", sdID, err)
				}
			}
		}
	}
	return nil
}

func Impl(autowireType, key string, param interface{}) (interface{}, error) {
	return impl(autowireType, key, param, false)
}

func ImplWithProxy(autowireType, key string, param interface{}) (interface{}, error) {
	return impl(autowireType, key, param, true)
}

func impl(autowireType, key string, param interface{}, withProxy bool) (interface{}, error) {
	targetSDID := GetSDIDByAliasIfNecessary(key)

	for _, wrapperAutowire := range allWrapperAutowires {
		if wrapperAutowire.TagKey() == autowireType {
			return wrapperAutowire.ImplWithParam(targetSDID, param, withProxy)
		}
	}
	logger.Red("[Autowire] SDID %s with autowire type %s not found in all autowires", key, autowireType)
	return nil, fmt.Errorf("[Autowire] SDID %s with autowire type %s not found in all autowires", key, autowireType)
}

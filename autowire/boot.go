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

func printAutowireRegisteredStructDescriptor() {
	for autowireType, aw := range GetAllWrapperAutowires() {
		logger.Blue("[Autowire Type] Found registered autowire type %s", autowireType)
		for sdID := range aw.GetAllStructDescriptors() {
			logger.Blue("[Autowire Struct Descriptor] Found type %s registered SD %s", autowireType, sdID)
		}
	}
}

func Load() error {
	printAutowireRegisteredStructDescriptor()

	// autowire all struct that can be entrance
	for _, aw := range GetAllWrapperAutowires() {
		for sdID := range aw.GetAllStructDescriptors() {
			sd := GetStructDescriptor(sdID)
			if sd == nil {
				continue
			}
			if parseCommonLoadAtOnceMetadataFromSDMetadata(sd.Metadata) || aw.CanBeEntrance() {
				_, err := aw.ImplWithoutParam(sdID, !sd.DisableProxy, false)
				if err != nil {
					return fmt.Errorf("[Autowire] Impl sd %s failed, reason is %s", sdID, err)
				}
			}
		}
	}
	return nil
}

func Impl(autowireType, key string, param interface{}) (interface{}, error) {
	return impl(autowireType, key, param, false, false)
}

func ImplWithProxy(autowireType, key string, param interface{}) (interface{}, error) {
	return impl(autowireType, key, param, true, false)
}

func ImplByForce(autowireType, key string, param interface{}) (interface{}, error) {
	return impl(autowireType, key, param, false, true)
}

func impl(autowireType, key string, param interface{}, expectWithProxy, force bool) (interface{}, error) {
	targetSDID := GetSDIDByAliasIfNecessary(key)

	// check expectWithProxy flag
	sd := GetStructDescriptor(targetSDID)
	if sd != nil && sd.DisableProxy {
		// if proxy is disabled by struct descriptor, set expectWithProxy to false
		expectWithProxy = false
	}

	for _, wrapperAutowire := range GetAllWrapperAutowires() {
		if wrapperAutowire.TagKey() == autowireType {
			return wrapperAutowire.ImplWithParam(targetSDID, param, expectWithProxy, force)
		}
	}
	logger.Red("[Autowire] SDID %s with autowire type %s not found in all autowires", key, autowireType)
	return nil, fmt.Errorf("[Autowire] SDID %s with autowire type %s not found in all autowires", key, autowireType)
}

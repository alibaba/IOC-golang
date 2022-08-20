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

package allimpls

import (
	"github.com/alibaba/ioc-golang/autowire"
)

const InterfaceMetadataKey = "interfaces"
const AutowireTypeMetadataKey = "autowireType"

func parseAllImpledIntefacesFromSDMetadata(metadata autowire.Metadata) []interface{} {
	autowireMetadata := autowire.ParseAutowireMetadataFromSDMetadata(metadata)
	if autowireMetadata == nil {
		return nil
	}
	allimplsMetadata, ok := autowireMetadata[autowire.CommonMetadataKey].(map[string]interface{})
	if !ok {
		return nil
	}
	result, ok := allimplsMetadata[autowire.CommonImplementsMetadataKey].([]interface{})
	if !ok {
		return nil
	}
	return result
}

func parseAllImplsItemAutowireTypeFromSDMetadata(metadata autowire.Metadata) string {
	autowireMetadata := autowire.ParseAutowireMetadataFromSDMetadata(metadata)
	if autowireMetadata == nil {
		return ""
	}
	allimplsMetadata, ok := autowireMetadata[Name].(map[string]interface{})
	if !ok {
		return ""
	}
	result, ok := allimplsMetadata[AutowireTypeMetadataKey].(string)
	if !ok {
		return ""
	}
	return result
}

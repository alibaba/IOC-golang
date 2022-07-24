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
	"fmt"
	"reflect"

	"github.com/alibaba/ioc-golang/autowire"
)

type sdIDParser struct {
	sdidFieldInterfaceReflectTypeMap map[string]reflect.Type
}

/*
Parse support parse field like:
ResourceServiceClient resources.ResourceServiceClient `grpc:"resource-service"`
to struct descriptor ID 'ResourceServiceClient-ResourceServiceClient'
*/
func (p *sdIDParser) Parse(fi *autowire.FieldInfo) (string, error) {
	if fi.FieldReflectType.Kind() != reflect.Slice {
		return "", fmt.Errorf("[Autowire allimpls] invalid field %s, field should be interface slice", fi)
	}
	if fi.FieldReflectType.Elem().Kind() != reflect.Interface {
		return "", fmt.Errorf("[Autowire allimpls] invalid field %s, field should be interface slice", fi)
	}
	return fi.FieldType, nil
}

func (p *sdIDParser) newFieldInterfaceSliceValue(sdid string) interface{} {
	fieldInterfaceReflectType, ok := p.sdidFieldInterfaceReflectTypeMap[sdid]
	if !ok {
		// never occurs
		return nil
	}

	return reflect.MakeSlice(reflect.SliceOf(fieldInterfaceReflectType), 0, 0)
}

var sdidParserSingleton *sdIDParser

func getSDIDParserSingleton() *sdIDParser {
	if sdidParserSingleton == nil {
		sdidParserSingleton = &sdIDParser{
			sdidFieldInterfaceReflectTypeMap: impledInterfaceSDIDTypeMap,
		}
	}
	return sdidParserSingleton
}

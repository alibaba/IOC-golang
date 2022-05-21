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

package triple

import (
	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/autowire/util"
)

type sdIDParser struct {
}

/*
Parse support parse field like:
ResourceServiceClient resources.ResourceServiceClient `triple:"resource-service"`
to struct descriptor ID 'ResourceServiceClient-ResourceServiceClient'
*/
func (p *sdIDParser) Parse(fi *autowire.FieldInfo) (string, error) {
	grpcInterfaceName := fi.FieldType
	return util.GetIdByNamePair(grpcInterfaceName, grpcInterfaceName), nil
}

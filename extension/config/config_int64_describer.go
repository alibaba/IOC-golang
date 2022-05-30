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

package config

//
// Why?
//
// In many scenarios, int64 may be us
// such as the nextId by the snowflake algorithm.
//

// +ioc:autowire=true
// +ioc:autowire:baseType=true
// +ioc:autowire:type=config
// +ioc:autowire:paramType=ConfigInt64
// +ioc:autowire:constructFunc=New
// +ioc:autowire:alias=ConfigInt64

type ConfigInt64 int64

func (ci *ConfigInt64) Value() int64 {
	return int64(*ci)
}

func (ci *ConfigInt64) New(impl *ConfigInt64) (*ConfigInt64, error) {
	*impl = *ci
	return impl, nil
}

func FromInt64(val int64) *ConfigInt64 {
	configInt64 := ConfigInt64(val)
	return &configInt64
}

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
// In many scenarios, float64 may be use,
// such as the price of order.
//

// +ioc:autowire=true
// +ioc:autowire:baseType=true
// +ioc:autowire:type=config
// +ioc:autowire:paramType=ConfigFloat64
// +ioc:autowire:constructFunc=New
// +ioc:autowire:alias=ConfigFloat64

type ConfigFloat64 float64

func (ci *ConfigFloat64) Value() float64 {
	return float64(*ci)
}

func (ci *ConfigFloat64) New(impl *ConfigFloat64) (*ConfigFloat64, error) {
	*impl = *ci
	return impl, nil
}

func FromFloat64(val float64) *ConfigFloat64 {
	configFloat64 := ConfigFloat64(val)
	return &configFloat64
}

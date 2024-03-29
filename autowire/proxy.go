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

// proxy function

var pf func(interface{}) interface{}

func RegisterProxyFunction(f func(interface{}) interface{}) {
	pf = f
}

func GetProxyFunction() func(interface{}) interface{} {
	if pf == nil {
		return func(i interface{}) interface{} {
			return i
		}
	}
	return pf
}

// proxy impl function

var pif func(interface{}, interface{}, string) error

func RegisterProxyImplFunction(f func(interface{}, interface{}, string) error) {
	pif = f
}

func GetProxyImplFunction() func(interface{}, interface{}, string) error {
	if pif == nil {
		return func(interface{}, interface{}, string) error {
			return nil
		}
	}
	return pif
}

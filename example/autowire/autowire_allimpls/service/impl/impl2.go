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

package impl

import "fmt"

// +ioc:autowire=true
// +ioc:autowire:type=allimpls
// +ioc:autowire:allimpls:interface=github.com/alibaba/ioc-golang/example/autowire/autowire_allimpls/service.Service

type serviceImpl2 struct {
}

func (s *serviceImpl2) GetHelloString(name string) string {
	return fmt.Sprintf("This is ServiceImpl1, hello %s", name)
}

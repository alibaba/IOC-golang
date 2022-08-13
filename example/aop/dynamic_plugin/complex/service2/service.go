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

package service2

import (
	"github.com/alibaba/ioc-golang/example/aop/dynamic_plugin/complex/service1"
	normalRedis "github.com/alibaba/ioc-golang/extension/state/redis"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:constructFunc=constructService2

type Service2 struct {
	Service1Singleton service1.Service1IOCInterface `singleton:""`
	Service1Normal    service1.Service1IOCInterface `normal:""`

	NormalDB1Redis normalRedis.RedisIOCInterface `normal:",db1-redis"`

	name string
}

func constructService2(s *Service2) (*Service2, error) {
	s.name = "default"
	return s, nil
}

func (s *Service2) LoadData() string {
	return ""
}

func (s *Service2) SetData(val string) {

}

func (s *Service2) GetName() string {
	return s.name
}

func (s *Service2) SetName(name string) {
	s.name = name
}

func (s *Service2) GetService1Normal() service1.Service1IOCInterface {
	return s.Service1Normal
}

func (s *Service2) GetService1Singleton() service1.Service1IOCInterface {
	return s.Service1Singleton
}

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

package main

import (
	"github.com/alibaba/ioc-golang/example/aop/dynamic_plugin/complex/service1"
	normalRedis "github.com/alibaba/ioc-golang/extension/state/redis"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type Service2 struct {
	Service1Singleton service1.Service1IOCInterface `singleton:""`
	Service1Normal    service1.Service1IOCInterface `normal:""`

	NormalDB1Redis normalRedis.RedisIOCInterface `normal:",db1-redis"`
}

func (s *Service2) LoadData() string {
	rsp, _ := s.NormalDB1Redis.Get("key1").Result()
	return rsp
}

func (s *Service2) SetData(val string) {
	s.NormalDB1Redis.Set("key1", val, -1)
}

func (s *Service2) GetName() string {
	return "plugin"
}

func (s *Service2) SetName(name string) {
}

func (s *Service2) GetService1Normal() service1.Service1IOCInterface {
	return s.Service1Normal
}

func (s *Service2) GetService1Singleton() service1.Service1IOCInterface {
	return s.Service1Singleton
}

// nolint
var Service2Plugin = Service2{}

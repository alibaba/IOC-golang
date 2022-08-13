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

package service1

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:type=normal
// +ioc:autowire:constructFunc=constructFunc

type Service1 struct {
	name string
}

func constructFunc(s *Service1) (*Service1, error) {
	s.name = "default"
	return s, nil
}

func (s *Service1) SetName(name string) {
	s.name = name
}

func (s *Service1) GetName() string {
	return s.name
}

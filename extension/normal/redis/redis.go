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

package redis

import (
	"time"

	"github.com/go-redis/redis"

	"github.com/alibaba/IOC-Golang/autowire/normal"
)

type Redis interface {
	GetRawClient() *redis.Client
	Set(key string, value interface{}, expiration time.Duration) (string, error)
	Get(key string) (string, error)
	HGetAll(key string) (map[string]string, error)
}

const SDID = "Redis-Impl"

// +ioc:autowire=true
// +ioc:autowire:type=normal
// +ioc:autowire:interface=Redis
// +ioc:autowire:paramType=Config
// +ioc:autowire:constructFunc=New

type Impl struct {
	client *redis.Client
}

func (i *Impl) GetRawClient() *redis.Client {
	return i.client
}

func (i *Impl) HGetAll(key string) (map[string]string, error) {
	return i.client.HGetAll(key).Result()
}

func (i *Impl) Get(key string) (string, error) {
	return i.client.Get(key).Result()
}

func (i *Impl) Set(key string, value interface{}, expiration time.Duration) (string, error) {
	return i.client.Set(key, value, expiration).Result()
}

var _ Redis = &Impl{}

func GetRedis(config *Config) (Redis, error) {
	mysqlImpl, err := normal.GetImpl(SDID, config)
	if err != nil {
		return nil, err
	}
	return mysqlImpl.(Redis), nil
}

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

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func clearEnv() {
	os.Setenv("REDIS_ADDRESS", "")
	os.Setenv(EnvKeyIOCGolangConfigPath, "")
	os.Setenv(EnvKeyIOCGolangEnv, "")
}

func Test_parseConfigSource(t *testing.T) {
	defer clearEnv()
	assert.Nil(t, os.Setenv(EnvKeyIOCGolangConfigPath, "./test/ioc_golang-config-source-env.yaml"))
	assert.Nil(t, os.Setenv("REDIS_ADDRESS", "localhost:16379"))
	assert.Nil(t, os.Setenv("REDIS_ADDRESS_EXPAND", "localhost:6388"))

	assert.Nil(t, Load())

	t.Run("test with env source ", func(t *testing.T) {

		redisConfig := &redisConfig{}

		assert.Nil(t, LoadConfigByPrefix("autowire.normal.<github.com/alibaba/ioc-golang/extension/normal/redis.Impl>.expand", redisConfig))
		assert.Equal(t, "0", redisConfig.DB)
		assert.Equal(t, "localhost:6388", redisConfig.Address)

		assert.Nil(t, LoadConfigByPrefix("autowire.normal.<github.com/alibaba/ioc-golang/extension/normal/redis.Impl>.param", redisConfig))
		assert.Equal(t, "0", redisConfig.DB)
		assert.Equal(t, "localhost:6379", redisConfig.Address)

		assert.Nil(t, LoadConfigByPrefix("autowire.normal.<github.com/alibaba/ioc-golang/extension/normal/redis.Impl>.env-redis.param", redisConfig))
		assert.Equal(t, "1", redisConfig.DB)
		assert.Equal(t, "localhost:16379", redisConfig.Address)

		assert.Nil(t, LoadConfigByPrefix("autowire.normal.<github.com/alibaba/ioc-golang/extension/normal/redis.Impl>.normal-redis.param", redisConfig))
		assert.Equal(t, "2", redisConfig.DB)
		assert.Equal(t, "localhost:26379", redisConfig.Address)
	})
}

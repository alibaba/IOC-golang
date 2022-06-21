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
	"strings"

	"github.com/fatih/color"
)

const (
	ConfigSourceKey     = "_ioc_golang_config_source"
	ConfigSourceEnvFlag = "env"
	ConfigEnvKeyPrefix  = "${"
	ConfigEnvKeySuffix  = "}"
)

func parseConfigSource(config Config) {
	envFlag := false
	if source, ok := config[ConfigSourceKey]; ok {
		if sourceStr, okStr := source.(string); okStr && sourceStr == ConfigSourceEnvFlag {
			color.Blue("[Config] %s under %v is set to %s, try to load from env", ConfigSourceKey, config, ConfigSourceEnvFlag)
			envFlag = true
		}
	}
	for k, v := range config {
		if val, ok := v.(string); ok {
			if expandValue, expand := ExpandConfigValueIfNecessary(val); expand {
				config[k] = expandValue
				continue
			}

			if envFlag {
				if envVal := os.Getenv(val); envVal != "" {
					config[k] = envVal
				} else {
					color.Blue("[Config] Try to load %s from env failed", val)
				}
			}
		} else if subConfig, ok := v.(Config); ok {
			parseConfigSource(subConfig)
		}
	}
}

func ExpandConfigValueIfNecessary(targetValue interface{}) (interface{}, bool) {
	if tv, ok := targetValue.(string); ok {
		if strings.HasPrefix(tv, ConfigEnvKeyPrefix) && strings.HasSuffix(tv, ConfigEnvKeySuffix) {
			expandValue := os.ExpandEnv(tv)
			if expandValue != "" {
				return expandValue, true
			}
		}
	}

	return targetValue, false
}

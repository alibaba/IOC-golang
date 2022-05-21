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
	"io/ioutil"

	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

var (
	supportedConfigTypes = []string{"yml", "yaml"}
)

func LoadOptions(opts ...Option) error {
	options := initOptions(opts...)
	if notSupportConfigType(options.ConfigType) {
		color.Red("[Config] Config file type:[%s] not supported now(yml, yaml)", options.ConfigType)
		return nil
	}

	targetMap := make(Config)

	configFiles := searchConfigFiles(options)
	for _, cf := range configFiles {
		contents, err := ioutil.ReadFile(cf)
		if err != nil {
			color.Red("[Config] Load ioc-golang config file failed. %v\n The load procedure is continue", err)
			return nil
		}

		var sub Config
		err = yaml.Unmarshal(contents, &sub)
		if err != nil {
			color.Red("[Config] yamlFile Unmarshal err: %v", err)
			return err
		}

		if len(sub) > 0 {
			targetMap = MergeMap(targetMap, sub)
		}

	}
	config = targetMap
	parseConfigSource(config)

	return nil
}

func initOptions(opts ...Option) *Options {
	options := newOptions()
	for _, opt := range opts {
		opt(options)
	}
	options.validate()

	return options
}

func notSupportConfigType(configType string) bool {
	return !stringSliceContains(supportedConfigTypes, configType)
}

func stringSliceContains(haystack []string, needle string) bool {
	for _, hs := range haystack {
		if hs == needle {
			return true
		}
	}

	return false
}

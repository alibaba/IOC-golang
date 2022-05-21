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
	"strings"

	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

type Config AnyMap

var config Config

func SetConfig(yamlBytes []byte) error {
	return yaml.Unmarshal(yamlBytes, &config)
}

func Load() error {
	configPath := GetConfigPath()
	color.Blue("[Config] Load config file from %s", configPath)

	absPath := determineAbsPath(configPath)

	yamlFile, err := ioutil.ReadFile(absPath)
	if err != nil {
		color.Red("Load ioc-golang config file failed. %v\nThe load procedure is continue\n", err)
		return nil
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		color.Red("yamlFile Unmarshal err: %v\n", err)
		return err
	}
	parseConfigSource(config)
	return nil
}

// LoadConfigByPrefix prefix is a.b.c, configStructPtr is interface ptr
func LoadConfigByPrefix(prefix string, configStructPtr interface{}) error {
	if configStructPtr == nil {
		return nil
	}
	configProperties := strings.Split(prefix, ".")
	return loadProperty(configProperties, 0, config, configStructPtr)
}

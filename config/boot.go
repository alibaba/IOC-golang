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
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

type Config map[string]interface{}

const (
	DefaultCallerSkipStep int = 3
)

var config Config

func SetConfig(yamlBytes []byte) error {
	return yaml.Unmarshal(yamlBytes, &config)
}

func Load(skips ...int) error {
	// skips - the number of stack frames to ascend, default 3
	skip := DefaultCallerSkipStep
	switch len(skips) {
	case 1:
		skip = skips[0]
	}

	configPath := GetConfigPath()

	if IsNotAbsPath(configPath) {
		currentDir := getCurrentAbPathByCaller(skip)
		color.Blue("[Config] Parse current dir %s", currentDir)
		configPath = filepath.Join(currentDir, configPath)
	}

	color.Blue("[Config] Load config file from %s", configPath)
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		color.Red("Load ioc-golang config file failed. %v\n The load procedure is continue\n", err)
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

// getCurrentAbPathByCaller - get current path by caller
func getCurrentAbPathByCaller(skip int) string {
	var currentAbsPath string
	_, filename, _, ok := runtime.Caller(skip)
	if ok {
		currentAbsPath = path.Dir(filename)
	}

	return currentAbsPath
}

// IsNotAbsPath - reports whether the path isn't absolute
func IsNotAbsPath(path string) bool {
	return !filepath.IsAbs(path)
}

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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"

	perrors "github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const (

	// EnvKeyIOCGolangConfigPath is absolute/relate path to ioc_golang.yaml
	EnvKeyIOCGolangConfigPath = "IOC_GOLANG_CONFIG_PATH" // default val is "../conf/ioc_golang.yaml"

	// EnvKeyIOCGolangEnv if is set to dev,then:
	// 1. choose config center with namespace dev
	// 2. choose config path like "config/ioc_golang_dev.yaml
	EnvKeyIOCGolangEnv = "IOC_GOLANG_ENV" //

	DefaultConfigPath = "../conf/ioc_golang.yaml"

	PathSeparator = string(os.PathSeparator)
	emptyString   = ""
	dotSeparator  = "."
)

func GetEnv() string {
	return os.Getenv(EnvKeyIOCGolangEnv)
}

func GetConfigPath() string {
	configPath := ""
	env := GetEnv()

	configFilePath := DefaultConfigPath
	if iocGolangConfigPath := os.Getenv(EnvKeyIOCGolangConfigPath); iocGolangConfigPath != "" {
		color.Blue("[Config] Environment %s is set to %s", EnvKeyIOCGolangConfigPath, iocGolangConfigPath)
		configFilePath = iocGolangConfigPath
	}
	prefix := strings.Split(configFilePath, ".yaml")
	// prefix == ["config/ioc_golang", ""]
	if len(prefix) != 2 {
		panic("Invalid config file path = " + configFilePath)
	}
	// get target env yaml file
	if env != "" {
		color.Blue("[Config] Environment %s is set to %s", EnvKeyIOCGolangEnv, env)
		configPath = prefix[0] + "_" + env + ".yaml"
	} else {
		configPath = configFilePath
	}
	return configPath
}

func loadProperty(splitedConfigName []string, index int, tempConfigMap Config, configStructPtr interface{}) error {
	subConfig, ok := tempConfigMap[splitedConfigName[index]]
	if !ok {
		return perrors.Errorf("property %s's key %s not found", splitedConfigName, splitedConfigName[index])
	}
	if index+1 == len(splitedConfigName) {
		targetConfigByte, err := yaml.Marshal(subConfig)
		if err != nil {
			return perrors.Errorf("property %s's key %s invalid, error = %s", splitedConfigName, splitedConfigName[index], err)
		}
		err = yaml.Unmarshal(targetConfigByte, configStructPtr)
		if err != nil {
			return perrors.Errorf("property %s's key %s doesn't match type %+v, error = %s", splitedConfigName, splitedConfigName[index], configStructPtr, err)
		}
		return nil
	}
	subMap, ok := subConfig.(Config)
	if !ok {
		return perrors.Errorf("property %s's key %s of config is not map[string]string, which is %+v", splitedConfigName,
			splitedConfigName[index], subConfig)
	}
	return loadProperty(splitedConfigName, index+1, subMap, configStructPtr)
}

func searchConfigFiles(opts *Options) []string {
	configFiles := make([]string, 0)

	if isNotEmptyStringSlice(opts.AbsPath) {
		for _, absPath := range opts.AbsPath {
			if !filepath.IsAbs(absPath) {
				panic(fmt.Sprintf("[Config] %s, abs path?", absPath))
			}
			configFiles = append(configFiles, absPath)
		}

		return configFiles
	}

	configNames := determineConfigFileName(opts)
	for _, configName := range configNames {
	PATH:
		for _, path := range opts.SearchPath {
			searchPath := determinePathSuffix(path) + configName
			absPath := determineAbsPath(searchPath)
			if stringSliceContains(configFiles, absPath) {
				continue PATH
			}
			if ok, _ := fileExists(absPath); ok { // xxx/config.yml
				configFiles = append(configFiles, absPath)
			}
		}
	}

	return configFiles
}

func isBlankString(src string) bool {
	return "" == src || "" == strings.TrimSpace(src)
}

func isNotBlankString(src string) bool {
	return !isBlankString(src)
}

func isEmptyStringSlice(src []string) bool {
	return 0 == len(src)
}

func isNotEmptyStringSlice(src []string) bool {
	return !isEmptyStringSlice(src)
}

func determineAbsPath(path string) string {
	if path == emptyString {
		path = dotSeparator
	}
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}

	p, err := filepath.Abs(path)
	if err == nil {
		return filepath.Clean(p)
	}

	return emptyString
}

func determinePathSuffix(searchPath string) string {
	if searchPath == emptyString {
		searchPath = dotSeparator
	}
	if strings.HasSuffix(searchPath, PathSeparator) {
		return searchPath
	}

	return searchPath + PathSeparator
}

func determineConfigFileName(opts *Options) []string {
	configNames := make([]string, len(opts.ProfilesActive)+1)              // dev,share
	configName := populateConfigName(opts.ConfigName, "", opts.ConfigType) // config
	configNames[0] = configName
	for i, profile := range opts.ProfilesActive {
		configNames[i+1] = populateConfigName(opts.ConfigName, profile, opts.ConfigType) // config_dev.yml, config_share.yml
	}

	return configNames
}

func populateConfigName(configName, profile, configType string) string {
	if profile == emptyString {
		return fmt.Sprintf("%s.%s", configName, configType) // config.yml
	}
	return fmt.Sprintf("%s_%s.%s", configName, profile, configType) // config_dev.yml
}

func fileExists(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err == nil {
		return !stat.IsDir(), nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

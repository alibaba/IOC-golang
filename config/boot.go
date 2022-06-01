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
type Option func(opts *Options)

const (
	YmlExtension            = "yml"
	YamlExtension           = "yaml"
	DefaultSearchConfigName = "config"
	DefaultSearchConfigType = YmlExtension // yaml

	emptySlice = 0

	YamlConfigSeparator = "."
)

var (
	config               Config
	supportedConfigTypes = []string{YmlExtension, YamlExtension}
	DefaultSearchPath    = []string{".", "./config", "./configs"}
)

type Options struct {
	AbsPath []string // abs path -> high level priority

	// default: config
	ConfigName string
	// default: yml,yaml
	ConfigType string
	// Search path of config files
	//
	// default: ., ./config, ./configs;
	//
	// priority: ./ > ./config > ./configs
	SearchPath []string
	// Profiles active
	//
	// default: []string
	//
	// e.g.:
	//
	// []string{"dev","share"}
	//
	// search target: config.yml/config_dev.yml/config_share.yml
	ProfilesActive []string
	// Depth of merging under multiple config files
	MergeDepth uint8
}

func (opts *Options) validate() {
	if isBlankString(opts.ConfigName) {
		opts.ConfigName = DefaultSearchConfigName
	}
	if isBlankString(opts.ConfigType) {
		opts.ConfigType = DefaultSearchConfigType
	}
	if isEmptyStringSlice(opts.SearchPath) {
		opts.SearchPath = DefaultSearchPath
	}
	if opts.MergeDepth == 0 {
		opts.MergeDepth = defaultMergeDepth
	}
}

// ----------------------------------------------------------------

func SetConfig(yamlBytes []byte) error {
	return yaml.Unmarshal(yamlBytes, &config)
}

func Load(opts ...Option) error {
	options := initOptions(opts...)
	if notSupportConfigType(options.ConfigType) {
		color.Red("[Config] Config file type:[%s] not supported now(yml, yaml)", options.ConfigType)
		return nil
	}

	targetMap := make(Config)

	configFiles := searchConfigFiles(options)
	if len(opts) == emptySlice {
		defaultConfigFile := loadDefaultConfigFileIfNecessary()
		if isNotBlankString(defaultConfigFile) {
			configFiles = append(configFiles, defaultConfigFile)
		}
	}

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

func loadDefaultConfigFileIfNecessary() string {
	configPath := GetConfigPath()
	color.Blue("[Config] Load default config file from %s", configPath)
	if isBlankString(configPath) {
		return configPath
	}
	absPath := determineAbsPath(configPath)

	return absPath
}

// ----------------------------------------------------------------

func WithAbsPath(absPath ...string) Option {
	return func(opts *Options) {
		opts.AbsPath = absPath
	}
}

func WithConfigName(configName string) Option {
	return func(opts *Options) {
		opts.ConfigName = configName
	}
}

func WithConfigType(configType string) Option {
	return func(opts *Options) {
		opts.ConfigType = configType
	}
}

func WithSearchPath(searchPath ...string) Option {
	return func(opts *Options) {
		opts.SearchPath = searchPath
	}
}

func WithProfilesActive(profilesActive ...string) Option {
	return func(opts *Options) {
		opts.ProfilesActive = profilesActive
	}
}

func WithMergeDepth(mergeDepth uint8) Option {
	return func(opts *Options) {
		opts.MergeDepth = mergeDepth
	}
}

// ----------------------------------------------------------------

// LoadConfigByPrefix prefix is like 'a.b.c' or 'a.b.<github.com/xxx/xx/xxx.Impl>.c', configStructPtr is interface ptr
func LoadConfigByPrefix(prefix string, configStructPtr interface{}) error {
	if configStructPtr == nil {
		return nil
	}
	splited := strings.Split(prefix, "<")
	var configProperties []string
	if len(splited) == 1 {
		configProperties = strings.Split(prefix, YamlConfigSeparator)
	} else {
		configProperties = strings.Split(splited[0], YamlConfigSeparator)
		backSplitedList := strings.Split(splited[1], ">")
		configProperties = append(configProperties, backSplitedList[0])
		configProperties = append(configProperties, strings.Split(backSplitedList[1], YamlConfigSeparator)...)
	}
	realConfigProperties := make([]string, 0)
	for _, v := range configProperties {
		if v != "" {
			realConfigProperties = append(realConfigProperties, v)
		}
	}
	return loadProperty(realConfigProperties, 0, config, configStructPtr)
}

// ----------------------------------------------------------------

func newOptions() *Options {
	return &Options{
		SearchPath:     make([]string, 0),
		ProfilesActive: make([]string, 0),
	}
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

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
	"os"
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
	DefaultSearchConfigType = YamlExtension // yaml

	SearchPathEnvKey    = "IOC_GOLANG_CONFIG_SEARCH_PATH"
	TypeEnvKey          = "IOC_GOLANG_CONFIG_TYPE"
	NameEnvKey          = "IOC_GOLANG_CONFIG_NAME"
	ActiveProfileEnvKey = "IOC_GOLANG_CONFIG_ACTIVE_PROFILE"

	YamlConfigSeparator = "."
	EnvValueSeparator   = ","
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

func (opts *Options) printLogs() {
	color.Blue("[Config] Config files load options is %+v", *opts)
}

func (opts *Options) loadFromEnv() {
	opts.SearchPath = loadSplitedStringsFromEnvWith(SearchPathEnvKey)
	opts.ConfigType = os.Getenv(TypeEnvKey)
	opts.ConfigName = os.Getenv(NameEnvKey)
	opts.ProfilesActive = loadSplitedStringsFromEnvWith(ActiveProfileEnvKey)
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

	options.printLogs()
	configFiles := searchConfigFiles(options)

	for _, cf := range configFiles {
		color.Blue("[Config] Loading config file %s", cf)
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
	return nil
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
	options.loadFromEnv()
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

func loadSplitedStringsFromEnvWith(envKey string) []string {
	val := os.Getenv(envKey)
	if isBlankString(val) {
		return []string{}
	} else {
		return strings.Split(val, EnvValueSeparator)
	}
}

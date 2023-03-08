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

	"github.com/alibaba/ioc-golang/logger"

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
	activeProfile        = make([]string, 0)
	supportedConfigTypes = []string{YmlExtension, YamlExtension}
	DefaultSearchPath    = []string{".", "./config", "./configs"}
)

type Options struct {
	AbsPath []string // abs path -> high level priority

	// default: config
	ConfigName string
	// default: yaml
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
	// Properties set by API
	Properties AnyMap
}

func (opts *Options) printLogs() {
	logger.Blue("[Config] Config files load options is %+v", *opts)
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
		logger.Red("[Config] Config file type:[%s] not supported now(yml, yaml)", options.ConfigType)
		return nil
	}

	targetMap := make(Config)

	options.printLogs()

	// set profile
	activeProfile = options.ProfilesActive

	configFiles := searchConfigFiles(options)

	for _, cf := range configFiles {
		logger.Blue("[Config] Loading config file %s", cf)
		contents, err := os.ReadFile(cf)
		if err != nil {
			logger.Red("[Config] Load ioc-golang config file failed. %v\n The load procedure is continue", err)
			return nil
		}

		var sub Config
		err = yaml.Unmarshal(contents, &sub)
		if err != nil {
			logger.Red("[Config] yamlFile Unmarshal err: %v", err)
			return err
		}

		if len(sub) > 0 {
			targetMap = MergeMap(targetMap, sub)
		}
	}
	addProperties(targetMap, options.Properties)

	// set config
	config = targetMap

	parseConfigIfNecessary(config)

	return nil
}

func addProperties(config Config, properties AnyMap) {
	for k, v := range properties {
		addProperty2Map(config, k, v)
	}
}

func addProperty2Map(cfg Config, key string, val interface{}) {
	first, others := separateFirstPrefixUnit(key)
	if first == "" {
		return
	}
	if others == "" {
		cfg[first] = val
		return
	}
	subMap, ok := cfg[first]
	if ok {
		addProperty2Map(subMap.(Config), others, val)
	} else {
		cfg[first] = make(Config)
		addProperty2Map(cfg[first].(Config), others, val)
	}
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

func AddProperty(key string, value interface{}) Option {
	return func(opts *Options) {
		opts.Properties[key] = value
	}
}

// ----------------------------------------------------------------

// LoadConfigByPrefix prefix is like 'a.b.c' or 'a.b.<github.com/xxx/xx/xxx.Impl>.c', configStructPtr is interface ptr
func LoadConfigByPrefix(prefix string, configStructPtr interface{}) error {
	if configStructPtr == nil {
		return nil
	}
	realConfigProperties := make([]string, 0)
	for _, v := range splitPrefix2Units(prefix) {
		if v != "" {
			v = expandIfNecessary(v)
			realConfigProperties = append(realConfigProperties, v)
		}
	}
	return loadProperty(realConfigProperties, 0, config, configStructPtr)
}

func GetActiveProfiles() []string {
	return activeProfile
}

func splitPrefix2Units(prefix string) []string {
	configProperties := make([]string, 0)
	if prefix == "" {
		return configProperties
	}
	splited := strings.Split(prefix, "<")
	if len(splited) == 1 {
		configProperties = strings.Split(prefix, YamlConfigSeparator)
	} else {
		if splited[0] != "" {
			configProperties = strings.Split(splited[0], YamlConfigSeparator)
		}
		backSplitedList := strings.Split(splited[1], ">")
		configProperties = append(configProperties, backSplitedList[0])
		if backSplitedList[1] != "" {
			configProperties = append(configProperties, strings.Split(strings.TrimPrefix(backSplitedList[1], YamlConfigSeparator), YamlConfigSeparator)...)
		}
	}
	return configProperties
}

func separateFirstPrefixUnit(prefix string) (string, string) {
	prefixUnits := splitPrefix2Units(prefix)
	if len(prefixUnits) == 0 {
		return "", ""
	}
	if len(prefixUnits) == 1 {
		return prefixUnits[0], ""
	}
	firstUnit := prefixUnits[0]
	if prefix[0] == '<' {
		return firstUnit, strings.TrimPrefix(prefix, "<"+firstUnit+">"+YamlConfigSeparator)
	}
	return firstUnit, strings.TrimPrefix(prefix, firstUnit+YamlConfigSeparator)
}

// ----------------------------------------------------------------

func newOptions() *Options {
	return &Options{
		SearchPath:     make([]string, 0),
		ProfilesActive: make([]string, 0),
		Properties:     make(AnyMap, 0),
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

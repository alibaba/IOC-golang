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
		color.Red("Config file type:[%s] not supported now(yml, yaml)\n", options.ConfigType)
		return nil
	}

	targetMap := make(Config)

	configFiles := searchConfigFiles(options)
	for _, cf := range configFiles {
		contents, err := ioutil.ReadFile(cf)
		if err != nil {
			color.Red("Load ioc-golang config file failed. %v\n The load procedure is continue\n", err)
			return nil
		}

		var sub Config
		err = yaml.Unmarshal(contents, &sub)
		if err != nil {
			color.Red("yamlFile Unmarshal err: %v\n", err)
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

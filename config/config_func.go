package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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
		configNames[i+1] = populateConfigName(opts.ConfigName, profile, opts.ConfigType) // config_test.yml
	}

	return configNames
}

func populateConfigName(configName, profile, configType string) string {
	if profile == emptyString {
		return fmt.Sprintf("%s.%s", configName, configType) // config.yml
	}
	return fmt.Sprintf("%s_%s.%s", configName, profile, configType) // config_test.yml
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

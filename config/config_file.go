package config

import (
	"os"
)

const (
	PathSeparator = string(os.PathSeparator)
	emptyString   = ""
	dotSeparator  = "."
)

func searchConfigFiles(opts *Options) []string {
	configFiles := make([]string, 0)

	if isNotEmptyStringSlice(opts.AbsPath) {
		for _, absPath := range opts.AbsPath {
			configFiles = append(configFiles, absPath)

			return configFiles
		}
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

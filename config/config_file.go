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

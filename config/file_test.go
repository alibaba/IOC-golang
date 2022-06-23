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
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func clearEnv() {
	os.Unsetenv("REDIS_ADDRESS")
	os.Unsetenv(TypeEnvKey)
	os.Unsetenv(SearchPathEnvKey)
	os.Unsetenv(NameEnvKey)
	os.Unsetenv(ActiveProfileEnvKey)
}

func Test_searchConfigFiles(t *testing.T) {
	type args struct {
		opts *Options
	}

	wd, _ := os.Getwd()

	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Test searchConfigFiles()",
			args: args{
				opts: &Options{
					ConfigName:     "ioc_golang",
					ProfilesActive: []string{},
					ConfigType:     "yaml",
					SearchPath:     []string{".", "test"},
				},
			},
			want: []string{filepath.Join(wd, "test/ioc_golang.yaml")},
		},
		{
			name: "Test searchConfigFiles()-.",
			args: args{
				opts: &Options{
					ConfigName:     "boot",
					ProfilesActive: []string{"test"},
					ConfigType:     "go",
					SearchPath:     []string{"."},
				},
			},
			want: []string{filepath.Join(wd, "boot.go"), filepath.Join(wd, "boot_test.go")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, searchConfigFiles(tt.args.opts), "searchConfigFiles(%v)", tt.args.opts)
		})
	}
}

func Test_determineAbsPath(t *testing.T) {
	type args struct {
		path string
	}

	wd, _ := os.Getwd()

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test determineAbsPath()-\"\"",
			args: args{
				path: "",
			},
			want: wd,
		},
		{
			name: "Test determineAbsPath()-.",
			args: args{
				path: ".",
			},
			want: wd,
		},
		{
			name: "Test determineAbsPath()-test",
			args: args{
				path: "test",
			},
			want: filepath.Join(wd, "test"),
		},
		{
			name: "Test determineAbsPath()-./test/ioc_golang.yaml",
			args: args{
				path: "./test/ioc_golang.yaml",
			},
			want: filepath.Join(wd, "./test/ioc_golang.yaml"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, determineAbsPath(tt.args.path), "determineAbsPath(%v)", tt.args.path)
		})
	}
}

func Test_determinePathSuffix(t *testing.T) {
	type args struct {
		searchPath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test determinePathSuffix()-.",
			args: args{
				searchPath: ".",
			},
			want: "." + PathSeparator,
		},
		{
			name: "Test determinePathSuffix()-./",
			args: args{
				searchPath: "." + PathSeparator,
			},
			want: "." + PathSeparator,
		},
		{
			name: "Test determinePathSuffix()-./config",
			args: args{
				searchPath: "." + PathSeparator + "config",
			},
			want: "." + PathSeparator + "config" + PathSeparator,
		},
		{
			name: "Test determinePathSuffix()-\"\"",
			args: args{
				searchPath: "",
			},
			want: "." + PathSeparator,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, determinePathSuffix(tt.args.searchPath), "determinePathSuffix(%v)", tt.args.searchPath)
		})
	}
}

func Test_determineConfigFileName(t *testing.T) {
	type args struct {
		opts *Options
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Test determineConfigFileName()-dev-test-prod",
			args: args{
				opts: &Options{
					ConfigName:     "config",
					ProfilesActive: []string{"dev", "test", "prod"},
					ConfigType:     "yml",
				},
			},
			want: []string{"config.yml", "config_dev.yml", "config_test.yml", "config_prod.yml"},
		},
		{
			name: "Test determineConfigFileName()",
			args: args{
				opts: &Options{
					ConfigName:     "config",
					ProfilesActive: []string{},
					ConfigType:     "yml",
					MergeDepth:     16,
				},
			},
			want: []string{"config.yml"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, determineConfigFileName(tt.args.opts), "determineConfigFileName(%v)", tt.args.opts)
		})
	}
}

func Test_populateConfigName(t *testing.T) {
	type args struct {
		configName string
		profile    string
		configType string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test populateConfigName()-config_dev.yml",
			args: args{
				configName: "config",
				profile:    "dev",
				configType: "yml",
			},
			want: "config_dev.yml",
		},
		{
			name: "Test populateConfigName()-config_dev.yaml",
			args: args{
				configName: "config",
				profile:    "dev",
				configType: "yaml",
			},
			want: "config_dev.yaml",
		},
		{
			name: "Test populateConfigName()-config_test.yml",
			args: args{
				configName: "config",
				profile:    "test",
				configType: "yml",
			},
			want: "config_test.yml",
		},
		{
			name: "Test populateConfigName()-config_test.yaml",
			args: args{
				configName: "config",
				profile:    "test",
				configType: "yaml",
			},
			want: "config_test.yaml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, populateConfigName(tt.args.configName, tt.args.profile, tt.args.configType), "populateConfigName(%v, %v, %v)", tt.args.configName, tt.args.profile, tt.args.configType)
		})
	}
}

func Test_fileExists(t *testing.T) {
	type args struct {
		path string
	}

	wd, _ := os.Getwd()

	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Test fileExists()-exist",
			args: args{
				path: filepath.Join(wd, "boot.go"),
			},
		},
		{
			name: "Test fileExists()-exist",
			args: args{
				path: filepath.Join(wd, "config_file.go"),
			},
		},
		{
			name: "Test fileExists()-not-exist",
			args: args{
				path: filepath.Join(wd, "hello.conf"),
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fileExists(tt.args.path)
			if (err == nil) != tt.wantErr {
				return
			}
			assert.Equalf(t, tt.want, got, "fileExists(%v)", tt.args.path)
		})
	}
}

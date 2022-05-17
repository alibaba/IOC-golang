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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfigPath(t *testing.T) {
	defer clearEnv()
	tests := []struct {
		iocGolangConfigPath string
		isDefault           bool
		env                 string
		name                string
		want                string
	}{
		{
			isDefault: true,
			name:      "default config path",
			want:      DefaultConfigPath,
		},
		{
			isDefault: true,
			name:      "default config path with en",
			env:       "dev",
			want:      "../conf/ioc_golang_dev.yaml",
		},
		{
			isDefault:           false,
			iocGolangConfigPath: "./test/ioc_golang.yaml",
			name:                "given config path with env ",
			env:                 "dev",
			want:                "./test/ioc_golang_dev.yaml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.isDefault {
				assert.Nil(t, os.Setenv(EnvKeyIOCGolangConfigPath, tt.iocGolangConfigPath))
			}
			if tt.env != "" {
				assert.Nil(t, os.Setenv(EnvKeyIOCGolangEnv, tt.env))
			}
			assert.Equalf(t, tt.want, GetConfigPath(), "GetConfigPath()")
		})
	}
}

func TestGetIOCGolangEnv(t *testing.T) {
	defer clearEnv()
	tests := []struct {
		name string
		want string
	}{
		{
			"test set env",
			"dev",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Nil(t, os.Setenv(EnvKeyIOCGolangEnv, tt.want))
			assert.Equalf(t, tt.want, GetEnv(), "GetEnv()")
		})
	}
}

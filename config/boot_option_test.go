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
	"testing"

	perrors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestLoad_empty(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test Load()",
			args: args{
				opts: []Option{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Load(tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}

			strValue := ""
			assert.Error(t, perrors.New("property [autowire config strValue]'s key autowire not found"),
				LoadConfigByPrefix("autowire#config#strValue", &strValue))
		})
	}
}

func TestLoad_options(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test Load()",
			args: args{
				opts: []Option{
					WithConfigName("ioc_golang"),
					WithConfigType("yaml"),
					WithSearchPath("./test"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Load(tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}

			strValue := "strValue"
			intValue := 123
			assert.Nil(t, LoadConfigByPrefix("autowire#config#strValue", &strValue))
			assert.Nil(t, LoadConfigByPrefix("autowire#config#intValue", &intValue))
		})
	}
}

func TestLoad_profile_active(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test Load()",
			args: args{
				opts: []Option{
					WithConfigName("ioc_golang"),
					WithConfigType("yaml"),
					WithSearchPath("./test"),
					WithProfilesActive("dev"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Load(tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}

			strValue := "strValue"
			intValue := 123
			boolValue := true
			mapValue := "mapValue1"
			assert.Nil(t, LoadConfigByPrefix("autowire#config#strValue", &strValue))
			assert.Nil(t, LoadConfigByPrefix("autowire#config#intValue", &intValue))
			assert.Nil(t, LoadConfigByPrefix("profilesActive#shared#boolValue", &boolValue))
			assert.Nil(t, LoadConfigByPrefix("profilesActive#shared#mapValue#mapKey1", &mapValue))

			sliceValue := []string{}
			assert.Nil(t, LoadConfigByPrefix("profilesActive#shared#sliceValue", &sliceValue))
			assert.Equal(t, "sliceStr1", sliceValue[0])
			assert.Equal(t, "sliceStr2", sliceValue[1])
			assert.Equal(t, "sliceStr3", sliceValue[2])
		})
	}
}

func TestLoad_abs_path(t *testing.T) {
	type args struct {
		opts []Option
	}

	wd, _ := os.Getwd()

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test Load()",
			args: args{
				opts: []Option{
					WithAbsPath(
						filepath.Join(wd, "./test/ioc_golang.yaml"),
						filepath.Join(wd, "./test/ioc_golang_dev.yaml"),
					),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Load(tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}

			strValue := "strValue"
			intValue := 123
			boolValue := true
			mapValue := "mapValue1"
			assert.Nil(t, LoadConfigByPrefix("autowire#config#strValue", &strValue))
			assert.Nil(t, LoadConfigByPrefix("autowire#config#intValue", &intValue))
			assert.Nil(t, LoadConfigByPrefix("profilesActive#shared#boolValue", &boolValue))
			assert.Nil(t, LoadConfigByPrefix("profilesActive#shared#mapValue#mapKey1", &mapValue))

			sliceValue := []string{}
			assert.Nil(t, LoadConfigByPrefix("profilesActive#shared#sliceValue", &sliceValue))
			assert.Equal(t, "sliceStr1", sliceValue[0])
			assert.Equal(t, "sliceStr2", sliceValue[1])
			assert.Equal(t, "sliceStr3", sliceValue[2])
		})
	}
}

func TestLoad_abs_path_panic(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test Load()-panic",
			args: args{
				opts: []Option{WithAbsPath("./test/ioc_golang.yaml")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					assert.Equal(t, fmt.Sprintf("%v", err), "[Config] ./test/ioc_golang.yaml, abs path?")
				}
			}()
			if err := Load(tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

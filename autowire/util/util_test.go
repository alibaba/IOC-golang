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

package util

import (
	"reflect"
	"testing"
)

type StructFoo struct {
}

type InterfaceFoo interface {
}

func TestGetIdByInterfaceAndImplPtr(t *testing.T) {
	type args struct {
		interfaceStruct interface{}
		implStructPtr   interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test get id by interface and impl",
			args: args{
				interfaceStruct: new(InterfaceFoo),
				implStructPtr:   &StructFoo{},
			},
			want: "InterfaceFoo-StructFoo",
		},
		{
			name: "test get id by impl",
			args: args{
				implStructPtr: &StructFoo{},
			},
			want: "-StructFoo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetIdByInterfaceAndImplPtr(tt.args.interfaceStruct, tt.args.implStructPtr); got != tt.want {
				t.Errorf("GetIdByInterfaceAndImplPtr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetIdByNamePair(t *testing.T) {
	type args struct {
		interfaceName string
		structPtrName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "get id by name pair",
			args: args{
				interfaceName: "InterfaceFoo",
				structPtrName: "StructFoo",
			},
			want: "InterfaceFoo-StructFoo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetIdByNamePair(tt.args.interfaceName, tt.args.structPtrName); got != tt.want {
				t.Errorf("GetIdByNamePair() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetStructName(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "get struct name",
			args: args{
				v: &StructFoo{},
			},
			want: "StructFoo",
		},
		{
			name: "get interface name",
			args: args{
				v: new(InterfaceFoo),
			},
			want: "InterfaceFoo",
		},
		{
			name: "get nil name",
			args: args{
				v: nil,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetStructName(tt.args.v); got != tt.want {
				t.Errorf("GetStructName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTypeFromInterface(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want reflect.Type
	}{
		{
			name: "get type from interface",
			args: args{
				v: new(InterfaceFoo),
			},
			want: func() reflect.Type {
				return reflect.TypeOf(new(InterfaceFoo)).Elem()
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTypeFromInterface(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTypeFromInterface() = %v, want %v", got, tt.want)
			}
		})
	}
}

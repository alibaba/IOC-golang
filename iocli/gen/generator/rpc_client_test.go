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

package generator

import (
	"testing"

	"github.com/alibaba/ioc-golang/iocli/gen/generator/plugin/common"
)

func Test_matchFunctionByStructName(t *testing.T) {
	type args struct {
		functionSignature string
		structName        string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantStr string
	}{
		{
			args: args{
				functionSignature: "func(s * MyStruct)(s*MyStruct)(string,error){",
				structName:        "MyStruct",
			},
			want:    true,
			wantStr: "(s*MyStruct)(string,error){",
		},
		{
			args: args{
				functionSignature: "func   (  s    *     MyStruct)(s*MyStruct)error{",
				structName:        "MyStruct",
			},
			want:    true,
			wantStr: "(s*MyStruct)error{",
		},
		{
			args: args{
				functionSignature: "func   (  s    *     MyStruct)(s*MyStruct){",
				structName:        "MyStruct2",
			},
			want: false,
		},
		{
			args: args{
				functionSignature: "func   (  s    *     MyStruct2)(s*MyStruct){",
				structName:        "MyStruct",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if body, ok := common.MatchFunctionByStructName(tt.args.functionSignature, tt.args.structName); ok != tt.want || body != tt.wantStr {
				t.Errorf("matchFunctionByStructName() = %s, %t, want %s, %t", body, ok, tt.wantStr, tt.want)
			}
		})
	}
}

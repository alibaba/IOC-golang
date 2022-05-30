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

package sdid_parser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang/autowire"
)

func TestGetDefaultSDIDParser(t *testing.T) {
	t.Run("Get Default SDID Parser equal", func(t *testing.T) {
		got1 := GetDefaultSDIDParser()
		got2 := GetDefaultSDIDParser()
		assert.Equal(t, got1, got2)
	})

	t.Run("Get Default SDID Parser not nil", func(t *testing.T) {
		assert.NotNil(t, GetDefaultSDIDParser())
	})
}

func Test_defaultSDIDParser_Parse(t *testing.T) {
	type args struct {
		fi *autowire.FieldInfo
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test default sdid parse normal interface field info",
			args: args{
				fi: &autowire.FieldInfo{
					FieldName: "MyRedis",
					FieldType: "Redis",
					TagKey:    "normal",
					TagValue:  "Impl",
				},
			},
			want:    "Impl#Redis",
			wantErr: false,
		},
		{
			name: "test default sdid parse normal struct ptr field info",
			args: args{
				fi: &autowire.FieldInfo{
					FieldName: "MyRedis",
					FieldType: "",
					TagKey:    "normal",
					TagValue:  "StructImpl",
				},
			},
			want:    "StructImpl",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &defaultSDIDParser{}
			got, err := p.Parse(tt.args.fi)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

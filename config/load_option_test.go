package config

import (
	"testing"

	perrors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestLoadOptions_empty(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test LoadOptions()",
			args: args{
				opts: []Option{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := LoadOptions(tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}

			strValue := ""
			assert.Error(t, perrors.New("property [autowire config strValue]'s key autowire not found"),
				LoadConfigByPrefix("autowire.config.strValue", &strValue))
		})
	}
}

func TestLoadOptions(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test LoadOptions()",
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
			if err := LoadOptions(tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}

			strValue := "strValue"
			intValue := 123
			assert.Nil(t, LoadConfigByPrefix("autowire.config.strValue", &strValue))
			assert.Nil(t, LoadConfigByPrefix("autowire.config.intValue", &intValue))
		})
	}
}

func TestLoadOptions_profile_active(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test LoadOptions()",
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
			if err := LoadOptions(tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}

			strValue := "strValue"
			intValue := 123
			boolValue := true
			mapValue := "mapValue1"
			assert.Nil(t, LoadConfigByPrefix("autowire.config.strValue", &strValue))
			assert.Nil(t, LoadConfigByPrefix("autowire.config.intValue", &intValue))
			assert.Nil(t, LoadConfigByPrefix("profilesActive.shared.boolValue", &boolValue))
			assert.Nil(t, LoadConfigByPrefix("profilesActive.shared.mapValue.mapKey1", &mapValue))

			sliceValue := []string{}
			assert.Nil(t, LoadConfigByPrefix("profilesActive.shared.sliceValue", &sliceValue))
			assert.Equal(t, "sliceStr1", sliceValue[0])
			assert.Equal(t, "sliceStr2", sliceValue[1])
			assert.Equal(t, "sliceStr3", sliceValue[2])
		})
	}
}

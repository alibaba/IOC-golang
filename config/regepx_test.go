package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isEnv(t *testing.T) {
	type args struct {
		envValue string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test isEnv-true",
			args: args{
				envValue: "${REDIS_ADDRESS_EXPAND}",
			},
			want: true,
		},
		{
			name: "test isEnv-false-1",
			args: args{
				envValue: "REDIS_ADDRESS_EXPAND",
			},
			want: false,
		},
		{
			name: "test isEnv-false-2",
			args: args{
				envValue: "${REDIS_ADDRESS_EXPAND",
			},
			want: false,
		},
		{
			name: "test isEnv-false-3",
			args: args{
				envValue: "REDIS_ADDRESS_EXPAND}",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, isEnv(tt.args.envValue), "isEnv(%v)", tt.args.envValue)
		})
	}
}

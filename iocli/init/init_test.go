package init

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func Test_parseGOMod(t *testing.T) {
	type args struct {
		regex string
		src   string
		name  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test find()",
			args: args{
				regex: `module\s+(?P<name>[\S]+)`,
				src:   "module github.com/alibaba/ioc-golang\n\ngo 1.17\n\nrequire (...)\n",
				name:  "$name",
			},
			want: "github.com/alibaba/ioc-golang",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := find(tt.args.regex, tt.args.src, tt.args.name); got != tt.want {
				t.Errorf("find() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseEVNecessary(t *testing.T) {
	type args struct {
		path string
	}

	ps := string(os.PathSeparator)
	goopath := fmt.Sprintf("%sUsers%sphotowey%sgo", ps, ps, ps)
	_ = os.Setenv("GOOPATH", goopath)
	_ = os.Setenv("HELLO_WORLD", "hello")

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test parseEVNecessary()-1",
			args: args{
				path: fmt.Sprintf("$GOOPATH%ssrc%sgithub.com%salibaba%s$HELLO_WORLD", ps, ps, ps, ps),
			},
			want: fmt.Sprintf("%s%ssrc%sgithub.com%salibaba%shello", goopath, ps, ps, ps, ps),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseEVNecessary(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseEVNecessary() = %v, want %v", got, tt.want)
			}
		})
	}
}

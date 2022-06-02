package inject

import (
	"testing"
)

func Test_parseInterfacePackage(t *testing.T) {
	type args struct {
		serviceFullName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test parseInterfacePackage",
			args: args{
				serviceFullName: "github.com/author/project/package/subPackage/interfacePackage.InterfaceName",
			},
			want: "github.com/author/project/package/subPackage/interfacePackage",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseInterfacePackage(tt.args.serviceFullName); got != tt.want {
				t.Errorf("parseInterfacePackage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseInterfaceName(t *testing.T) {
	type args struct {
		serviceFullName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test parseInterfaceName",
			args: args{
				serviceFullName: "github.com/author/project/package/subPackage/interfacePackage.InterfaceName",
			},
			want: "InterfaceName",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseInterfaceName(tt.args.serviceFullName); got != tt.want {
				t.Errorf("parseInterfaceName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseInterfacePackageAlias(t *testing.T) {
	type args struct {
		c            *copyMethodMaker
		otherPackage string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test parseInterfacePackageAlias",
			args: args{
				c: &copyMethodMaker{
					importsList: &importsList{
						byPath:  make(map[string]string, 0),
						byAlias: make(map[string]string, 0),
					},
				},
				otherPackage: "github.com/author/project/package/subPackage/interfacePackage",
			},
			want: "interfacePackage",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseInterfacePackageAlias(tt.args.c, tt.args.otherPackage); got != tt.want {
				t.Errorf("parseInterfacePackageAlias() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isEligibleInterfaceReferencePath(t *testing.T) {
	type args struct {
		interfaceReferencePath string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test isEligibleInterfaceReferencePath-true",
			args: args{
				interfaceReferencePath: "github.com/author/project/package/subPackage/interfacePackage.InterfaceName",
			},
			want: true,
		},
		{
			name: "Test isEligibleInterfaceReferencePath-false-1",
			args: args{
				interfaceReferencePath: "github.com/author/project/package/subPackage/interfacePackage",
			},
			want: false,
		},
		{
			name: "Test isEligibleInterfaceReferencePath-false-2",
			args: args{
				interfaceReferencePath: "github.com/author/project/package/subPackage/interfacePackage/",
			},
			want: false,
		},
		{
			name: "Test isEligibleInterfaceReferencePath-false-3",
			args: args{
				interfaceReferencePath: "github.com/author/project/package/subPackage/interfacePackage/.",
			},
			want: false,
		},
		{
			name: "Test isEligibleInterfaceReferencePath-false-4",
			args: args{
				interfaceReferencePath: "InterfaceName",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isEligibleInterfaceReferencePath(tt.args.interfaceReferencePath); got != tt.want {
				t.Errorf("isEligibleInterfaceReferencePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

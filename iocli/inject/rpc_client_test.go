package inject

import "testing"

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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if body, ok := matchFunctionByStructName(tt.args.functionSignature, tt.args.structName); ok != tt.want || body != tt.wantStr {
				t.Errorf("matchFunctionByStructName() = %s, %t, want %s, %t", body, ok, tt.wantStr, tt.want)
			}
		})
	}
}

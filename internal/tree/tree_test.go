package tree

import "testing"

func Test_splitModuleVersion(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		module  string
		version string
	}{
		{
			name:    "no module or version",
			args:    args{s: ""},
			module:  "",
			version: "",
		},
		{
			name:    "module and version",
			args:    args{s: "github.com/ankittk/go-dependency-tree-parser@v0.1.0"},
			module:  "github.com/ankittk/go-dependency-tree-parser",
			version: "v0.1.0",
		},
		{
			name:    "module only",
			args:    args{s: "github.com/ankittk/go-dependency-tree-parser"},
			module:  "github.com/ankittk/go-dependency-tree-parser",
			version: "",
		},
		{
			name:    "version only",
			args:    args{s: "@v0.1.0"},
			module:  "",
			version: "v0.1.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := splitModuleVersion(tt.args.s)
			if got != tt.module {
				t.Errorf("splitModuleVersion() got = %v, want %v", got, tt.module)
			}
			if got1 != tt.version {
				t.Errorf("splitModuleVersion() got1 = %v, want %v", got1, tt.version)
			}
		})
	}
}

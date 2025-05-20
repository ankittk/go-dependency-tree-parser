package parser

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestDefaultModParser_RunGoModTidyAndGraph(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `
module github.com/ankittk/testmod

go 1.20

require (
	github.com/sirupsen/logrus v1.9.0
	github.com/stretchr/testify v1.8.0
)
`
	err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644)
	if err != nil {
		t.Fatalf("failed to create go.mod: %v", err)
	}

	mainGo := `
package main

import (
	_ "github.com/sirupsen/logrus"
	_ "github.com/stretchr/testify"
)

func main() {}
`
	err = os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(mainGo), 0644)
	if err != nil {
		t.Fatalf("failed to create main.go: %v", err)
	}

	modParser := NewDefaultModParser(tmpDir)

	t.Run("RunGoModTidy", func(t *testing.T) {
		if err := modParser.RunGoModTidy(); err != nil {
			t.Fatalf("RunGoModTidy failed: %v", err)
		}
	})

	t.Run("RunGoModGraph", func(t *testing.T) {
		output, err := modParser.RunGoModGraph()
		if err != nil {
			t.Fatalf("RunGoModGraph failed: %v", err)
		}

		expectedDeps := []string{
			"github.com/ankittk/testmod github.com/sirupsen/logrus@v1.9.0",
			"github.com/ankittk/testmod github.com/stretchr/testify@v1.8.0",
		}

		for _, dep := range expectedDeps {
			if !strings.Contains(output, dep) {
				t.Errorf("expected dependency %q not found in mod graph output", dep)
			}
		}
	})
}

func TestDefaultModParser_ParseModGraph(t *testing.T) {
	type fields struct {
		repoPath string
	}
	type args struct {
		modGraph string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string][]string
	}{
		{
			name:   "empty graph",
			fields: fields{repoPath: "/tmp/dtree/repo"},
			args:   args{modGraph: ""},
			want:   map[string][]string{},
		},
		{
			name:   "single dependency",
			fields: fields{repoPath: "/tmp/dtree/repo"},
			args:   args{modGraph: "A B\n"},
			want:   map[string][]string{"A": {"B"}},
		},
		{
			name:   "multiple dependencies",
			fields: fields{repoPath: "/tmp/dtree/repo"},
			args:   args{modGraph: "A B\nA C\nB D\n"},
			want: map[string][]string{
				"A": {"B", "C"},
				"B": {"D"},
			},
		},
		{
			name:   "self dependency",
			fields: fields{repoPath: "/tmp/dtree/repo"},
			args:   args{modGraph: "A A\n"},
			want:   map[string][]string{"A": {"A"}},
		},
		{
			name:   "trailing spaces",
			fields: fields{repoPath: "/tmp/dtree/repo"},
			args:   args{modGraph: "a@v1 b@v1 \n b@v1  c@v1 "},
			want: map[string][]string{
				"a@v1": {"b@v1"},
				"b@v1": {"c@v1"},
			},
		},
		{
			name:   "malformed lines ignored",
			fields: fields{repoPath: "/tmp/dtree/repo"},
			args:   args{modGraph: "a@v1 b@v1\nmalformedline\nb@v1 c@v1"},
			want: map[string][]string{
				"a@v1": {"b@v1"},
				"b@v1": {"c@v1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &DefaultModParser{
				repoPath: tt.fields.repoPath,
			}
			if got := p.ParseModGraph(tt.args.modGraph); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseModGraph() = %v, want %v", got, tt.want)
			}
		})
	}
}

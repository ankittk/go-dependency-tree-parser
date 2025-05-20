package tree

import (
	"reflect"
	"testing"
)

type mockModParser struct {
	tidyFunc  func() error
	graphFunc func() (string, error)
	parseFunc func(modGraph string) map[string][]string
}

func (m *mockModParser) RunGoModTidy() error {
	return m.tidyFunc()
}

func (m *mockModParser) RunGoModGraph() (string, error) {
	return m.graphFunc()
}

func (m *mockModParser) ParseModGraph(modGraph string) map[string][]string {
	return m.parseFunc(modGraph)
}

func TestBuildDependencyTree(t *testing.T) {
	mockParser := &mockModParser{
		tidyFunc: func() error {
			return nil
		},
		graphFunc: func() (string, error) {
			return "a@v1 b@v1\nb@v1 c@v1", nil
		},
		parseFunc: func(modGraph string) map[string][]string {
			return map[string][]string{
				"a@v1": {"b@v1"},
				"b@v1": {"c@v1"},
				"c@v1": {},
			}
		},
	}

	tb := &treeBuilder{modParser: mockParser}
	got, err := tb.BuildDependencyTree()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []*Artifact{
		{
			Name:    "a",
			Version: "v1",
			Dependencies: []*Artifact{
				{
					Name:    "b",
					Version: "v1",
					Dependencies: []*Artifact{
						{
							Name:         "c",
							Version:      "v1",
							Dependencies: []*Artifact{},
						},
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("BuildDependencyTree() = %+v, want %+v", got, want)
	}
}

func TestBuildTree(t *testing.T) {
	type args struct {
		module  string
		graph   map[string][]string
		visited map[string]*Artifact
		path    map[string]bool
	}
	tests := []struct {
		name string
		args args
		want *Artifact
	}{
		{
			name: "Simple tree with one dependency",
			args: args{
				module: "a@v1",
				graph: map[string][]string{
					"a@v1": {"b@v1"},
					"b@v1": {},
				},
				visited: map[string]*Artifact{},
				path:    map[string]bool{},
			},
			want: &Artifact{
				Name:    "a",
				Version: "v1",
				Dependencies: []*Artifact{
					{
						Name:         "b",
						Version:      "v1",
						Dependencies: []*Artifact{},
					},
				},
			},
		},
		{
			name: "Cycle is handled gracefully",
			args: args{
				module: "a@v1",
				graph: map[string][]string{
					"a@v1": {"b@v1"},
					"b@v1": {"a@v1"}, // cycle
				},
				visited: map[string]*Artifact{},
				path:    map[string]bool{},
			},
			want: &Artifact{
				Name:    "a",
				Version: "v1",
				Dependencies: []*Artifact{
					{
						Name:    "b",
						Version: "v1",
						Dependencies: []*Artifact{
							{
								Name:    "a",
								Version: "v1",
							},
						},
					},
				},
			},
		},
		{
			name: "Self-loop is skipped",
			args: args{
				module: "a@v1",
				graph: map[string][]string{
					"a@v1": {"a@v1"},
				},
				visited: map[string]*Artifact{},
				path:    map[string]bool{},
			},
			want: &Artifact{
				Name:         "a",
				Version:      "v1",
				Dependencies: []*Artifact{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := &treeBuilder{}
			got := tb.buildTree(tt.args.module, tt.args.graph, tt.args.visited, tt.args.path)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildTree() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

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

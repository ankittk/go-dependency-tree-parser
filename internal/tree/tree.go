package tree

import (
	"errors"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/ankittk/go-dependency-tree-parser/internal/parser"
)

// Artifact represents the structure to hold the dependency tree
type Artifact struct {
	Name    string `json:"name"`
	Version string `json:"version"`

	// Dependencies hold the list of child artifacts.
	// It is a nested structure, where each child artifact can have its own dependencies.
	// It is often nil or empty for leaf nodes.
	Dependencies []*Artifact `json:"dependencies,omitempty"`
}

// Builder is an interface that defines the method to build a dependency tree
type Builder interface {
	// BuildDependencyTree builds the dependency tree for the given module
	BuildDependencyTree() ([]*Artifact, error)
}

// treeBuilder is a struct that implements the Builder interface
type treeBuilder struct {
	modParser parser.ModParser
}

// NewBuilder creates a new instance of treeBuilder
func NewBuilder(modParser parser.ModParser) Builder {
	return &treeBuilder{
		modParser: modParser,
	}
}

// BuildDependencyTree builds the dependency tree for the given module
func (tb *treeBuilder) BuildDependencyTree() ([]*Artifact, error) {
	log.Debugf("Building dependency tree ...")

	if err := tb.modParser.RunGoModTidy(); err != nil {
		return nil, fmt.Errorf("failed to run go mod tidy: %w", err)
	}

	modGraph, err := tb.modParser.RunGoModGraph()
	if err != nil {
		return nil, fmt.Errorf("failed to run go mod graph: %w", err)
	}

	// parse the module graph into a map of parent-to-child dependencies
	graph := tb.modParser.ParseModGraph(modGraph)
	if len(graph) == 0 {
		return nil, errors.New("empty module graph")
	}

	// calculate incoming edges
	incoming := map[string]int{}
	for parent, children := range graph {
		if _, ok := incoming[parent]; !ok {
			incoming[parent] = 0
		}
		for _, child := range children {
			incoming[child]++
		}
	}

	// identifying root nodes
	var roots []string
	for mod, count := range incoming {
		if count == 0 {
			roots = append(roots, mod)
		}
	}

	visited := make(map[string]*Artifact)
	clusters := make([]*Artifact, 0, len(roots))

	// build the tree for each root node
	for _, root := range roots {
		clusters = append(clusters,
			tb.buildTree(root, graph, visited, make(map[string]bool)),
		)
	}

	return clusters, nil
}

// buildTree recursively builds the dependency tree for a module and
// uses the 'path' map to detect cycles and avoid infinite recursion
func (tb *treeBuilder) buildTree(module string, graph map[string][]string, visited map[string]*Artifact, path map[string]bool) *Artifact {
	if art, ok := visited[module]; ok {
		// If the artifact is already visited, check if it's in the current path
		if !path[module] {
			return art
		}
		return &Artifact{
			Name:    art.Name,
			Version: art.Version,
		}
	}

	name, version := splitModuleVersion(module)
	artifact := &Artifact{
		Name:         name,
		Version:      version,
		Dependencies: make([]*Artifact, 0),
	}

	// Mark the artifact as visited and add it to the current path
	visited[module] = artifact
	path[module] = true

	for _, dep := range graph[module] {
		if dep == module {
			// Skip self-loop
			continue
		}

		// this is a recursive call to build the tree for the dependency
		child := tb.buildTree(dep, graph, visited, path)
		artifact.Dependencies = append(artifact.Dependencies, child)
	}

	// backtrack: remove the module from the current path
	delete(path, module)
	return artifact
}

// splitModuleVersion splits a module string into its name and version components
func splitModuleVersion(s string) (string, string) {
	i := strings.LastIndex(s, "@")
	if i == -1 {
		return s, ""
	}

	// return the module name and version
	return s[:i], s[i+1:]
}

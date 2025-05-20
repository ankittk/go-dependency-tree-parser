package parser

import (
	"fmt"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

// ModParser is an interface that defines methods for parsing Go modules
type ModParser interface {
	RunGoModTidy() error
	RunGoModGraph() (string, error)
	ParseModGraph(modGraph string) map[string][]string
}

// DefaultModParser is a struct that implements the ModParser interface
type DefaultModParser struct {
	repoPath string
}

// NewDefaultModParser creates a new instance of DefaultModParser
func NewDefaultModParser(path string) ModParser {
	return &DefaultModParser{
		repoPath: path,
	}
}

// RunGoModTidy runs 'go mod tidy' command to clean up the go.mod file
func (p *DefaultModParser) RunGoModTidy() error {
	log.Debugf("Running go mod tidy ...")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = p.repoPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go mod tidy failed: %w", err)
	}
	return nil
}

// RunGoModGraph runs 'go mod graph' command and returns the output
func (p *DefaultModParser) RunGoModGraph() (string, error) {
	log.Debugf("Running go mod graph ...")
	cmd := exec.Command("go", "mod", "graph")
	cmd.Dir = p.repoPath
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("go mod graph failed: %w", err)
	}
	return string(output), nil
}

// ParseModGraph parses the output of 'go mod graph' and returns a map of parent to child dependencies
func (p *DefaultModParser) ParseModGraph(modGraph string) map[string][]string {
	log.Debugf("Parsing go mod graph ...")
	graph := make(map[string][]string)
	for _, line := range strings.Split(modGraph, "\n") {
		// Trim whitespace and skip empty lines
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Split the line into parent and child
		parts := strings.Fields(line)
		if len(parts) != 2 {
			continue
		}

		parent, child := parts[0], parts[1]
		graph[parent] = append(graph[parent], child)
	}
	return graph
}

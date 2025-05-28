package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/ankittk/go-dependency-tree-parser/internal/github"
	"github.com/ankittk/go-dependency-tree-parser/internal/parser"
	"github.com/ankittk/go-dependency-tree-parser/internal/tree"
)

var parseCmd = &cobra.Command{
	Use:   "parse <repo> <tag>",
	Short: "Analyze and cache the Go dependencies of a GitHub repository at a given tag or branch.",
	Long: `
The parse command analyzes and caches the Go dependencies of a GitHub repository at a given tag or branch.
It uses 'go mod' to fetch dependencies and stores them in a local cache (output.json).
`,
	Example: `
  dtree parse github.com/etcd-io/etcd v3.6.0-rc.0
  dtree parse github.com/etcd-io/etcd master
`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := args[0]
		tag := args[1]

		if verbose {
			log.SetLevel(log.DebugLevel)
		}

		return parseDependencies(repo, tag)
	},
}

func init() {
	rootCmd.AddCommand(parseCmd)
	parseCmd.PersistentFlags().StringP("repo", "r", "", "Repository to parse")
	parseCmd.PersistentFlags().StringP("tag", "t", "", "Tag or branch to parse")
}

// parseDependencies clones the repo, runs go mod graph, builds the tree and writes output.json
func parseDependencies(repo, tag string) error {
	log.Infof("Starting dependency parsing for repo: %s at %s", repo, tag)

	gitClient := github.NewDefaultGitClient()

	log.Debug("Cloning repository...")
	clonedPath, err := gitClient.Clone(repo, tag)
	if err != nil {
		log.Errorf("Failed to clone repository: %v", err)
		return err
	}
	log.Debugf("Repository cloned to: %s", clonedPath)

	goModPath := filepath.Join(clonedPath, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		log.Errorf("go.mod not found at path: %s", goModPath)
		return err
	}

	log.Debug("Parsing go.mod and building dependency graph...")
	modParser := parser.NewDefaultModParser(clonedPath)
	tb := tree.NewBuilder(modParser)
	artifactTree, err := tb.BuildDependencyTree()
	if err != nil {
		log.Errorf("Failed to build dependency tree: %v", err)
		return err
	}

	jsonBytes, err := json.MarshalIndent(artifactTree, "", "  ")
	if err != nil {
		log.Errorf("Failed to marshal dependency tree to JSON: %v", err)
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Errorf("Failed to get current working directory: %v", err)
		return err
	}
	outputFile := filepath.Join(cwd, "output.json")

	if err := os.WriteFile(outputFile, jsonBytes, 0644); err != nil {
		log.Errorf("Failed to write output to file: %v", err)
		return err
	}
	log.Infof("Dependency parsing completed successfully. Output written to %s", outputFile)

	return nil
}

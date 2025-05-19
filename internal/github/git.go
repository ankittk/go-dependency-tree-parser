package github

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

// GitClient is an interface that defines the methods for interacting with Git repositories
type GitClient interface {
	Clone(repoUrl, tag string) (string, error)
}

// DefaultGitClient is the default implementation of the GitClient interface
type DefaultGitClient struct {
	cacheDir string
}

func NewDefaultGitClient() GitClient {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Errorf("Error getting user home dir: %v", err)
		return nil
	}

	// Use ~/.dtree as the cache directory
	cacheDir := filepath.Join(homeDir, ".dtree")

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(cacheDir, os.ModePerm); err != nil {
		log.Errorf("Error creating cache dir: %v", err)
		return nil
	}

	return &DefaultGitClient{
		cacheDir: cacheDir,
	}
}

// Clone clones the specified Git repository at the given tag or branch
func (g *DefaultGitClient) Clone(repoURL, tag string) (string, error) {
	normalizedURL := normalizeGitURL(repoURL)

	// Extract the repo name from the URL path (e.g., etcd from github.com/etcd-io/etcd)
	repoName := filepath.Base(repoURL)

	// Destination directory: ~/.dtree/etcd-v3.5.5
	clonePath := filepath.Join(g.cacheDir, fmt.Sprintf("%s-%s", repoName, tag))

	if _, err := os.Stat(clonePath); err == nil {
		log.Debugf("%s already exists, fetching repo", clonePath)
		cmd := exec.Command("git", "-C", clonePath, "fetch", "origin", tag)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Errorf("git fetch failed: %v", err)
			return "", fmt.Errorf("git fetch failed: %w", err)
		}
		return clonePath, nil
	}

	// shallow clone as we only need the latest commit
	// --depth 1: only fetch the latest commit
	// --branch <tag>: checkout the specified tag or branch
	// --single-branch: only fetch the specified branch
	cmd := exec.Command("git", "clone", "--depth", "1", "--branch", tag, normalizedURL, clonePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Infof("Cloning %s at %s into %s", normalizedURL, tag, clonePath)

	if err := cmd.Run(); err != nil {
		log.Errorf("git clone failed: %v", err)
		return "", fmt.Errorf("git clone failed: %w", err)
	}

	return clonePath, nil
}

// normalizeGitURL normalizes the Git URL to ensure it has the correct format
func normalizeGitURL(input string) string {
	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") || strings.HasSuffix(input, ".git") {
		return input
	}
	return "https://" + input + ".git"
}

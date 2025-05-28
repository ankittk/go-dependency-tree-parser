# dtree - Go Dependency Tree Parser

**dtree** is a CLI tool to analyze and visualize the Go module dependencies of a GitHub repository at a specific tag.

It clones the repository, parses the `go.mod` file using `go mod graph`,
and outputs a raw dependency graph in structured JSON.
---

## Overview
- Clone a GitHub repository at a specific tag.
- Analyze the go.mod and recursively resolve all dependencies.
- Output a structured JSON dependency tree.
- Cache results for reuse and inspection.

---

## Example Output

The output is a JSON file structured as follows:

```json
{
  "module": {
    "path": "github.com/user/repo",
    "version": "v1.0.0",
    "dependencies": [
      {
        "path": "github.com/another/dependency",
        "version": "v1.2.3",
        "dependencies": []
      }
    ]
  }
}
```

---

## Installation

### Prerequisites

- Go 1.18 or later
- Git

### Build

To build the project:

```bash
make build
```

This will create a binary named `dtree` in the `bin/` directory.

### Run

To run the tool:

```bash
./bin/dtree parse <repository> <tag-or-branch>
```

Example:

```bash
./bin/dtree parse github.com/etcd-io/etcd v3.6.0
```
OR
You can also use the `make` command to run the tool:

```bash
make run repo=github.com/etcd-io/etcd tag=v3.6.0
```

---

## Testing

Run unit tests with:

```bash
make test
```


### Linting
To run linters:

```bash
make lint
```

---

## Cleanup

To remove generated files:

```bash
make clean
```

---

## CLI Flags

| Flag        | Description                                         |
|-------------|-----------------------------------------------------|
| `--repo`    | GitHub repository (e.g., `github.com/etcd-io/etcd`) |
| `--tag`     | Tag or branch to check out (e.g., `v3.6.0`)         |
| `--verbose` | Enable verbose debug output                         |

---

## Project Structure

```text
├── cmd
│   ├── dtree.go               # Command-line entry point logic
│   └── parse.go               # Argument parsing or subcommand logic
├── go.mod                     # Go module definition
├── go.sum                     # Go dependencies checksums
├── internal
│   ├── github
│   │   └── git.go             # GitHub repository handling (e.g., clone, checkout)
│   ├── parser
│   │   ├── modgraph.go        # go mod graph parsing logic
│   │   └── modgraph_test.go   # Unit tests for modgraph.go
│   └── tree
│       ├── tree.go            # Core tree-building logic
│       └── tree_test.go       # Unit tests for tree logic
├── LICENSE                    # License file
├── main.go                    # Main entry point for the application
├── Makefile                   # Build, test, lint automation
└── README.md                  # Project documentation
```

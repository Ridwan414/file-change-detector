# File Change Detector Client Library

A Go client library for detecting file changes using Merkle trees. This library provides a simple interface to create snapshots of directories and compare them to detect modifications, additions, and deletions.

## Features

- **Merkle Tree Based**: Uses SHA-256 hashing to create a cryptographic fingerprint of directory contents
- **Change Detection**: Identifies modified, added, and deleted files between snapshots
- **CSV Storage**: Stores snapshots in CSV format for easy inspection and portability
- **Client Interface**: Clean API for integration into other applications
- **Efficient Comparison**: Only changed branches of the tree need to be examined

## Installation

```bash
go get github.com/Ridwan414/file-change-detector
```

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    "log"
    "github.com/Ridwan414/file-change-detector/merkle"
)

func main() {
    // Create a client with storage directory
    client := merkle.NewClient("merkle_states")
    
    // Create a snapshot of a directory
    snapshot, err := client.CreateSnapshot("./my-folder")
    if err != nil {
        log.Fatal(err)
    }
    
    // Save the snapshot
    err = client.SaveSnapshot(snapshot, "./my-folder")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Snapshot created with root hash: %x\n", snapshot.RootHash)
}
```

### Comparing Snapshots

```go
// Find the latest snapshot
latestFile, err := client.FindLatestSnapshot("./my-folder")
if err != nil {
    log.Fatal("No previous snapshot found")
}

// Load previous snapshot
previousSnapshot, err := client.LoadSnapshot(latestFile)
if err != nil {
    log.Fatal(err)
}

// Create current snapshot
currentSnapshot, err := client.CreateSnapshot("./my-folder")
if err != nil {
    log.Fatal(err)
}

// Compare snapshots
report := client.CompareSnapshots(previousSnapshot, currentSnapshot)

// Display changes
merkle.PrintChangeReport(report)

// Or access changes programmatically
for _, change := range report.Changes {
    switch change.ChangeType {
    case merkle.Modified:
        fmt.Printf("Modified: %s\n", change.FileName)
    case merkle.Added:
        fmt.Printf("Added: %s\n", change.FileName)
    case merkle.Deleted:
        fmt.Printf("Deleted: %s\n", change.FileName)
    }
}
```

## API Reference

### Client Interface

```go
type Client interface {
    // Create a snapshot of a directory
    CreateSnapshot(folderPath string) (*TreeState, error)
    
    // Save a snapshot to storage
    SaveSnapshot(state *TreeState, folderPath string) error
    
    // Load a snapshot from storage
    LoadSnapshot(filename string) (*TreeState, error)
    
    // Find the most recent snapshot for a folder
    FindLatestSnapshot(folderPath string) (string, error)
    
    // Compare two snapshots
    CompareSnapshots(oldState, newState *TreeState) *ChangeReport
    
    // Get the Merkle tree for a folder
    GetTree(folderPath string) (*MerkleTree, error)
}
```

### Types

```go
// TreeState represents a snapshot
type TreeState struct {
    Timestamp  time.Time
    RootHash   []byte
    FileHashes map[string][]byte
}

// ChangeReport contains comparison results
type ChangeReport struct {
    OldTimestamp time.Time
    NewTimestamp time.Time
    OldRootHash  []byte
    NewRootHash  []byte
    Changes      []FileChange
}

// FileChange represents a single file change
type FileChange struct {
    FileName   string
    ChangeType ChangeType
    OldHash    []byte
    NewHash    []byte
}

// ChangeType enumeration
const (
    Modified ChangeType = iota
    Added
    Deleted
)
```

## Command Line Usage

The package includes a command-line tool:

```bash
# Create initial snapshot
go run main.go /path/to/folder

# Compare with previous snapshot
go run main.go /path/to/folder --compare
```

## Storage Format

Snapshots are stored as CSV files with the following format:
- Filename: `state_<foldername>_<timestamp>.csv`
- Columns: `timestamp,root_hash,file_path,file_hash`

## Use Cases

- **Backup Verification**: Ensure backup integrity by comparing snapshots
- **Change Monitoring**: Track modifications in configuration directories
- **Deployment Validation**: Verify file deployments match expectations
- **Security Auditing**: Detect unauthorized file modifications
- **Version Control**: Lightweight alternative for tracking file changes

## How It Works

1. **File Hashing**: Each file is hashed using SHA-256
2. **Tree Construction**: Files become leaf nodes, sorted alphabetically
3. **Parent Nodes**: Created by hashing concatenated child hashes
4. **Root Hash**: Final hash represents entire directory state
5. **Comparison**: Only branches with different hashes are examined

## Performance

- **O(n)** for creating snapshots (n = number of files)
- **O(log n)** average case for finding changes
- **Space efficient**: Only stores hashes, not file contents

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License
package merkle

import (
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// ChangeType represents the type of change detected
type ChangeType int

const (
	Modified ChangeType = iota
	Added
	Deleted
)

// FileChange represents a change detected in a file
type FileChange struct {
	FileName   string
	ChangeType ChangeType
	OldHash    []byte
	NewHash    []byte
}

// ChangeReport contains all changes detected between two states
type ChangeReport struct {
	OldTimestamp time.Time
	NewTimestamp time.Time
	OldRootHash  []byte
	NewRootHash  []byte
	Changes      []FileChange
}

// Client interface for the Merkle tree file change detector
type Client interface {
	// CreateSnapshot creates a Merkle tree snapshot of the specified folder
	CreateSnapshot(folderPath string) (*TreeState, error)

	// SaveSnapshot saves a tree state to storage
	SaveSnapshot(state *TreeState, folderPath string) error

	// LoadSnapshot loads a specific snapshot from storage
	LoadSnapshot(filename string) (*TreeState, error)

	// FindLatestSnapshot finds the most recent snapshot for a folder
	FindLatestSnapshot(folderPath string) (string, error)

	// CompareSnapshots compares two tree states and returns a change report
	CompareSnapshots(oldState, newState *TreeState) *ChangeReport

	// GetTree returns the Merkle tree for a folder
	GetTree(folderPath string) (*MerkleTree, error)
}

// MerkleClient implements the Client interface
type MerkleClient struct {
	storageDir string
}

// NewClient creates a new Merkle tree client
func NewClient(storageDir string) Client {
	return &MerkleClient{
		storageDir: storageDir,
	}
}

// MerkleNode represents a node in the Merkle tree
type MerkleNode struct {
	Hash     []byte
	Left     *MerkleNode
	Right    *MerkleNode
	IsLeaf   bool
	FileName string
}

// MerkleTree represents the complete Merkle tree
type MerkleTree struct {
	Root *MerkleNode
}

// TreeState represents a snapshot of the Merkle tree at a point in time
type TreeState struct {
	Timestamp  time.Time
	RootHash   []byte
	FileHashes map[string][]byte // filename -> hash
}

// CreateSnapshot creates a Merkle tree snapshot of the specified folder
func (c *MerkleClient) CreateSnapshot(folderPath string) (*TreeState, error) {
	tree, err := c.GetTree(folderPath)
	if err != nil {
		return nil, err
	}

	state := &TreeState{
		Timestamp:  time.Now(),
		RootHash:   tree.Root.Hash,
		FileHashes: make(map[string][]byte),
	}

	collectFileHashes(tree.Root, state.FileHashes)
	return state, nil
}

// GetTree returns the Merkle tree for a folder
func (c *MerkleClient) GetTree(folderPath string) (*MerkleTree, error) {
	return createMerkleTreeFromFolder(folderPath)
}

// SaveSnapshot saves a tree state to storage
func (c *MerkleClient) SaveSnapshot(state *TreeState, folderPath string) error {
	// Create storage directory if it doesn't exist
	if err := os.MkdirAll(c.storageDir, 0755); err != nil {
		return err
	}

	// Generate filename with timestamp
	filename := fmt.Sprintf("%s/state_%s_%s.csv", c.storageDir,
		filepath.Base(folderPath),
		state.Timestamp.Format("20060102_150405"))

	// Create CSV file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"timestamp", "root_hash", "file_path", "file_hash"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data rows
	timestampStr := state.Timestamp.Format(time.RFC3339)
	rootHashStr := hex.EncodeToString(state.RootHash)

	for fileName, hash := range state.FileHashes {
		row := []string{
			timestampStr,
			rootHashStr,
			fileName,
			hex.EncodeToString(hash),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// LoadSnapshot loads a specific snapshot from storage
func (c *MerkleClient) LoadSnapshot(filename string) (*TreeState, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, err
	}

	// Validate header
	expectedHeader := []string{"timestamp", "root_hash", "file_path", "file_hash"}
	for i, h := range expectedHeader {
		if header[i] != h {
			return nil, fmt.Errorf("invalid CSV header")
		}
	}

	state := &TreeState{
		FileHashes: make(map[string][]byte),
	}

	// Read data rows
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// Parse timestamp
		if state.Timestamp.IsZero() {
			state.Timestamp, _ = time.Parse(time.RFC3339, row[0])
		}

		// Parse root hash
		if state.RootHash == nil {
			state.RootHash, _ = hex.DecodeString(row[1])
		}

		// Parse file hash
		fileHash, _ := hex.DecodeString(row[3])
		state.FileHashes[row[2]] = fileHash
	}

	return state, nil
}

// FindLatestSnapshot finds the most recent snapshot for a folder
func (c *MerkleClient) FindLatestSnapshot(folderPath string) (string, error) {
	folderName := filepath.Base(folderPath)
	pattern := fmt.Sprintf("%s/state_%s_*.csv", c.storageDir, folderName)
	files, err := filepath.Glob(pattern)
	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return "", fmt.Errorf("no previous state found for folder: %s", folderName)
	}

	// Sort files by name (which includes timestamp)
	sort.Strings(files)

	// Return the most recent file
	return files[len(files)-1], nil
}

// CompareSnapshots compares two tree states and returns a change report
func (c *MerkleClient) CompareSnapshots(oldState, newState *TreeState) *ChangeReport {
	report := &ChangeReport{
		OldTimestamp: oldState.Timestamp,
		NewTimestamp: newState.Timestamp,
		OldRootHash:  oldState.RootHash,
		NewRootHash:  newState.RootHash,
		Changes:      []FileChange{},
	}

	// Find modified files
	for fileName, newHash := range newState.FileHashes {
		if oldHash, exists := oldState.FileHashes[fileName]; exists {
			if !equalHashes(oldHash, newHash) {
				report.Changes = append(report.Changes, FileChange{
					FileName:   fileName,
					ChangeType: Modified,
					OldHash:    oldHash,
					NewHash:    newHash,
				})
			}
		}
	}

	// Find added files
	for fileName, hash := range newState.FileHashes {
		if _, exists := oldState.FileHashes[fileName]; !exists {
			report.Changes = append(report.Changes, FileChange{
				FileName:   fileName,
				ChangeType: Added,
				NewHash:    hash,
			})
		}
	}

	// Find deleted files
	for fileName, hash := range oldState.FileHashes {
		if _, exists := newState.FileHashes[fileName]; !exists {
			report.Changes = append(report.Changes, FileChange{
				FileName:   fileName,
				ChangeType: Deleted,
				OldHash:    hash,
			})
		}
	}

	return report
}

// Helper functions (not exported)

func hashData(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func hashFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return nil, err
	}

	return hasher.Sum(nil), nil
}

func buildMerkleTree(nodes []*MerkleNode) *MerkleNode {
	if len(nodes) == 0 {
		return nil
	}

	if len(nodes) == 1 {
		return nodes[0]
	}

	var nextLevel []*MerkleNode

	for i := 0; i < len(nodes); i += 2 {
		var left, right *MerkleNode
		left = nodes[i]

		if i+1 < len(nodes) {
			right = nodes[i+1]
		} else {
			right = nodes[i]
		}

		combinedHash := append(left.Hash, right.Hash...)
		parentHash := hashData(combinedHash)

		parent := &MerkleNode{
			Hash:   parentHash,
			Left:   left,
			Right:  right,
			IsLeaf: false,
		}

		nextLevel = append(nextLevel, parent)
	}

	return buildMerkleTree(nextLevel)
}

func createMerkleTreeFromFolder(folderPath string) (*MerkleTree, error) {
	var leafNodes []*MerkleNode

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			fileHash, err := hashFile(path)
			if err != nil {
				return err
			}

			relPath, _ := filepath.Rel(folderPath, path)

			node := &MerkleNode{
				Hash:     fileHash,
				Left:     nil,
				Right:    nil,
				IsLeaf:   true,
				FileName: relPath,
			}

			leafNodes = append(leafNodes, node)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(leafNodes) == 0 {
		return nil, fmt.Errorf("no files found in folder")
	}

	sort.Slice(leafNodes, func(i, j int) bool {
		return leafNodes[i].FileName < leafNodes[j].FileName
	})

	root := buildMerkleTree(leafNodes)

	return &MerkleTree{Root: root}, nil
}

func collectFileHashes(node *MerkleNode, fileHashes map[string][]byte) {
	if node == nil {
		return
	}

	if node.IsLeaf {
		fileHashes[node.FileName] = node.Hash
	} else {
		collectFileHashes(node.Left, fileHashes)
		collectFileHashes(node.Right, fileHashes)
	}
}

func equalHashes(h1, h2 []byte) bool {
	if len(h1) != len(h2) {
		return false
	}
	for i := range h1 {
		if h1[i] != h2[i] {
			return false
		}
	}
	return true
}

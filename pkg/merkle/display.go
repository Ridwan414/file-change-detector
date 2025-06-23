package merkle

import (
	"fmt"
)

// PrintTree prints the Merkle tree structure
func PrintTree(node *MerkleNode, depth int) {
	if node == nil {
		return
	}

	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}

	if node.IsLeaf {
		fmt.Printf("%s[FILE] %s: %x\n", indent, node.FileName, node.Hash[:8])
	} else {
		fmt.Printf("%s[NODE] Hash: %x\n", indent, node.Hash[:8])
		PrintTree(node.Left, depth+1)
		PrintTree(node.Right, depth+1)
	}
}

// PrintChangeReport prints a formatted change report
func PrintChangeReport(report *ChangeReport) {
	fmt.Println("\n=== Change Detection Report ===")
	fmt.Printf("Comparing states from %s to %s\n",
		report.OldTimestamp.Format("2006-01-02 15:04:05"),
		report.NewTimestamp.Format("2006-01-02 15:04:05"))

	// Check root hash
	if !equalHashes(report.OldRootHash, report.NewRootHash) {
		fmt.Println("\nRoot hash changed - files have been modified")
		fmt.Printf("Old root: %x\n", report.OldRootHash[:16])
		fmt.Printf("New root: %x\n", report.NewRootHash[:16])
	} else {
		fmt.Println("\nNo changes detected - root hash is identical")
		return
	}

	// Count changes by type
	modifiedCount := 0
	addedCount := 0
	deletedCount := 0

	for _, change := range report.Changes {
		switch change.ChangeType {
		case Modified:
			modifiedCount++
		case Added:
			addedCount++
		case Deleted:
			deletedCount++
		}
	}

	// Print modified files
	fmt.Println("\nModified files:")
	if modifiedCount == 0 {
		fmt.Println("  None")
	} else {
		for _, change := range report.Changes {
			if change.ChangeType == Modified {
				fmt.Printf("  [MODIFIED] %s\n", change.FileName)
				fmt.Printf("    Old hash: %x\n", change.OldHash[:16])
				fmt.Printf("    New hash: %x\n", change.NewHash[:16])
			}
		}
	}

	// Print added files
	fmt.Println("\nAdded files:")
	if addedCount == 0 {
		fmt.Println("  None")
	} else {
		for _, change := range report.Changes {
			if change.ChangeType == Added {
				fmt.Printf("  [ADDED] %s (hash: %x)\n", change.FileName, change.NewHash[:16])
			}
		}
	}

	// Print deleted files
	fmt.Println("\nDeleted files:")
	if deletedCount == 0 {
		fmt.Println("  None")
	} else {
		for _, change := range report.Changes {
			if change.ChangeType == Deleted {
				fmt.Printf("  [DELETED] %s (hash: %x)\n", change.FileName, change.OldHash[:16])
			}
		}
	}

	fmt.Printf("\nSummary: %d modified, %d added, %d deleted\n",
		modifiedCount, addedCount, deletedCount)
}

// GetChangeTypeString returns a string representation of the change type
func GetChangeTypeString(changeType ChangeType) string {
	switch changeType {
	case Modified:
		return "MODIFIED"
	case Added:
		return "ADDED"
	case Deleted:
		return "DELETED"
	default:
		return "UNKNOWN"
	}
}

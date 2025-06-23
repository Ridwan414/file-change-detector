package main

import (
	"fmt"
	"os"

	"github.com/Ridwan414/file-change-detector/pkg/merkle"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <folder_path> [--compare]")
		fmt.Println("  --compare: Compare with the most recent saved state")
		os.Exit(1)
	}

	folderPath := os.Args[1]
	compareMode := false

	// Check for --compare flag
	if len(os.Args) > 2 && os.Args[2] == "--compare" {
		compareMode = true
	}

	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		fmt.Printf("Error: Folder '%s' does not exist\n", folderPath)
		os.Exit(1)
	}

	// Create client with storage directory
	client := merkle.NewClient("merkle_states")

	fmt.Printf("Creating Merkle tree for folder: %s\n", folderPath)

	// Get the Merkle tree
	tree, err := client.GetTree(folderPath)
	if err != nil {
		fmt.Printf("Error creating Merkle tree: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nMerkle Tree Root Hash: %x\n", tree.Root.Hash)
	fmt.Println("\nTree Structure:")
	merkle.PrintTree(tree.Root, 0)

	// Create current snapshot
	currentState, err := client.CreateSnapshot(folderPath)
	if err != nil {
		fmt.Printf("Error creating snapshot: %v\n", err)
		os.Exit(1)
	}

	// Compare with previous state if requested
	if compareMode {
		latestFile, err := client.FindLatestSnapshot(folderPath)
		if err != nil {
			fmt.Printf("\nNo previous state to compare with: %v\n", err)
		} else {
			fmt.Printf("\nLoading previous state from: %s\n", latestFile)
			previousState, err := client.LoadSnapshot(latestFile)
			if err != nil {
				fmt.Printf("Error loading previous state: %v\n", err)
			} else {
				report := client.CompareSnapshots(previousState, currentState)
				merkle.PrintChangeReport(report)
			}
		}
	}

	// Save current state
	if err := client.SaveSnapshot(currentState, folderPath); err != nil {
		fmt.Printf("Error saving tree state: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nTree state saved successfully\n")
}

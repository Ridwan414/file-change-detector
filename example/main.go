package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Ridwan414/file-change-detector/v1/pkg/merkle"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <folder_path>")
		os.Exit(1)
	}

	folderPath := os.Args[1]

	// Create a client with storage directory
	client := merkle.NewClient("merkle_states")

	// Check if this is the first run
	latestFile, err := client.FindLatestSnapshot(folderPath)
	if err != nil {
		// First run - create initial snapshot
		fmt.Println("Creating initial snapshot...")
		snapshot, err := client.CreateSnapshot(folderPath)
		if err != nil {
			log.Fatal(err)
		}

		// Save the snapshot
		err = client.SaveSnapshot(snapshot, folderPath)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Initial snapshot created with root hash: %x\n", snapshot.RootHash)
		return
	}

	// Load previous snapshot
	previousSnapshot, err := client.LoadSnapshot(latestFile)
	if err != nil {
		log.Fatal(err)
	}

	// Create current snapshot
	currentSnapshot, err := client.CreateSnapshot(folderPath)
	if err != nil {
		log.Fatal(err)
	}

	// Compare snapshots
	report := client.CompareSnapshots(previousSnapshot, currentSnapshot)

	// Display changes
	merkle.PrintChangeReport(report)

	// Save current snapshot
	err = client.SaveSnapshot(currentSnapshot, folderPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nSnapshot saved successfully!")
}

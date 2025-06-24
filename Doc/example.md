# File Change Detector - Example Usage

This document demonstrates how to use the File Change Detector with the test-folder structure, showing CSV data generation, file modification, and change detection.

## Test Folder Structure

The `test-folder` directory contains the following structure:

```
test-folder/
├── new-file.txt                                    (22B, 2 lines)
├── test-text-2.txt                                 (21B, 1 line)
├── test-folder-1/
│   ├── test-text-1.txt                             (21B, 1 line)
│   └── new-test-folder/
│       ├── test-text-1.txt                         (21B, 1 line)
│       └── test-text-2.txt                         (21B, 1 line)
└── test-folder-2/
    ├── test-text-1.txt                             (25B, 2 lines)
    └── test-text-2.txt                             (22B, 1 line)
```

## Step 1: Initial Snapshot and CSV Data Generation

Run the command to create the initial Merkle tree snapshot:

```bash
go run cmd/main.go test-folder
```

### Output:
```
Creating Merkle tree for folder: test-folder

Merkle Tree Root Hash: c3f0e775c1da05224cf3853734107741262c57b7036726cff69a2184a482c5be

Tree Structure:
[NODE] Hash: c3f0e775c1da0522
  [NODE] Hash: 8415a6deb572f9f1
    [NODE] Hash: b3d90d59e439a079
      [FILE] new-file.txt: 6f09abcb5b65e6e2
      [FILE] test-folder-1/new-test-folder/test-text-1.txt: bc6333c6f9c263d4
    [NODE] Hash: 721319829b8c52a3
      [FILE] test-folder-1/new-test-folder/test-text-2.txt: bf2bd557ba244ec6
      [FILE] test-folder-1/test-text-1.txt: ec0b0d64ff043208
  [NODE] Hash: 694beb4d84550ff3
    [NODE] Hash: b6d875518ec9ce9b
      [FILE] test-folder-2/test-text-1.txt: 79800b5bcf3e35f3
      [FILE] test-folder-2/test-text-2.txt: 28c8bb6b48b4681d
    [NODE] Hash: c653d18ba511e234
      [FILE] test-text-2.txt: 8da57ffc24f7e983
      [FILE] test-text-2.txt: 8da57ffc24f7e983

Tree state saved successfully
```

## Step 2: Generated CSV Data

The command generates a CSV file in the `merkle_states/` directory with the following format:

**File:** `merkle_states/state_test-folder_20250623_141911.csv`

```csv
timestamp,root_hash,file_path,file_hash
2025-06-23T14:19:11+06:00,c3f0e775c1da05224cf3853734107741262c57b7036726cff69a2184a482c5be,test-folder-2/test-text-1.txt,79800b5bcf3e35f34283199660e70fb26a750103378608efd17bf5fde03b2453
2025-06-23T14:19:11+06:00,c3f0e775c1da05224cf3853734107741262c57b7036726cff69a2184a482c5be,test-folder-2/test-text-2.txt,28c8bb6b48b4681d5421cd922a491008f55a9b1567b8eeef6f4645a6007b5264
2025-06-23T14:19:11+06:00,c3f0e775c1da05224cf3853734107741262c57b7036726cff69a2184a482c5be,test-text-2.txt,8da57ffc24f7e98320481f7aff6470799670a618566c6c395e4dbe1e8ea3db96
2025-06-23T14:19:11+06:00,c3f0e775c1da05224cf3853734107741262c57b7036726cff69a2184a482c5be,new-file.txt,6f09abcb5b65e6e29787b925b4a1086e7ba65788fafa265ba7761a26ee807e54
2025-06-23T14:19:11+06:00,c3f0e775c1da05224cf3853734107741262c57b7036726cff69a2184a482c5be,test-folder-1/new-test-folder/test-text-1.txt,bc6333c6f9c263d45f3e0d01a7a629bdb7e49dc86c5c4d8c8971f64ef20711a7
2025-06-23T14:19:11+06:00,c3f0e775c1da05224cf3853734107741262c57b7036726cff69a2184a482c5be,test-folder-1/new-test-folder/test-text-2.txt,bf2bd557ba244ec6783b4bbd3c3cbf6d2b19a187657820b02445212e240a83a1
2025-06-23T14:19:11+06:00,c3f0e775c1da05224cf3853734107741262c57b7036726cff69a2184a482c5be,test-folder-1/test-text-1.txt,ec0b0d64ff0432082e75ec8da7867595d1fe26d97a09a3baeafa1dc179e8a3fc
```

### CSV Data Explanation:
- **timestamp**: When the snapshot was taken
- **root_hash**: The root hash of the entire Merkle tree
- **file_path**: Relative path of each file in the directory
- **file_hash**: SHA-256 hash of each individual file

## Step 3: Modifying a File

To demonstrate change detection, we modify the content of `test-folder/new-file.txt`:

**Original content:**
```
This is a new file 23
```

**Modified content:**
```
This is a new file 23
This file has been modified for testing change detection!
```

## Step 4: Running the Compare Command

After modifying the file, run the compare command:

```bash
go run cmd/main.go test-folder --compare
```

### Output:
```
Creating Merkle tree for folder: test-folder

Merkle Tree Root Hash: 4007d952963ec6e065df27ce4fe776010a86cc17b7f990ec4d955e912a1ab62d

Tree Structure:
[NODE] Hash: 4007d952963ec6e0
  [NODE] Hash: 0fc77a4226752290
    [NODE] Hash: 17498693933de8b9
      [FILE] new-file.txt: 314211469cb73c75
      [FILE] test-folder-1/new-test-folder/test-text-1.txt: bc6333c6f9c263d4
    [NODE] Hash: 721319829b8c52a3
      [FILE] test-folder-1/new-test-folder/test-text-2.txt: bf2bd557ba244ec6
      [FILE] test-folder-1/test-text-1.txt: ec0b0d64ff043208
  [NODE] Hash: 694beb4d84550ff3
    [NODE] Hash: b6d875518ec9ce9b
      [FILE] test-folder-2/test-text-1.txt: 79800b5bcf3e35f3
      [FILE] test-folder-2/test-text-2.txt: 28c8bb6b48b4681d
    [NODE] Hash: c653d18ba511e234
      [FILE] test-text-2.txt: 8da57ffc24f7e983
      [FILE] test-text-2.txt: 8da57ffc24f7e983

Loading previous state from: merkle_states/state_test-folder_20250623_141911.csv

=== Change Detection Report ===
Comparing states from 2025-06-23 14:19:11 to 2025-06-23 14:19:38

Root hash changed - files have been modified
Old root: c3f0e775c1da05224cf3853734107741
New root: 4007d952963ec6e065df27ce4fe77601

Modified files:
  [MODIFIED] new-file.txt
    Old hash: 6f09abcb5b65e6e29787b925b4a1086e
    New hash: 314211469cb73c7598c0cbe884feb17e

Added files:
  None

Deleted files:
  None

Summary: 1 modified, 0 added, 0 deleted

Tree state saved successfully
```

## Key Observations

1. **Root Hash Change**: The root hash changed from `c3f0e775c1da0522...` to `4007d952963ec6e0...`, indicating changes in the directory structure.

2. **File Detection**: The system correctly identified that `new-file.txt` was modified.

3. **Hash Comparison**: 
   - Old hash: `6f09abcb5b65e6e29787b925b4a1086e`
   - New hash: `314211469cb73c7598c0cbe884feb17e`

4. **Precise Tracking**: Only the modified file is reported, showing the efficiency of the Merkle tree approach.

5. **Automatic Storage**: The new state is automatically saved for future comparisons.

## Usage Summary

1. **Initial Snapshot**: `go run cmd/main.go <folder_path>`
2. **Change Detection**: `go run cmd/main.go <folder_path> --compare`
3. **CSV Storage**: Snapshots are stored in `merkle_states/` directory
4. **Change Types**: Supports detection of modified, added, and deleted files

This demonstrates the complete workflow of the File Change Detector using Merkle trees for efficient change detection.

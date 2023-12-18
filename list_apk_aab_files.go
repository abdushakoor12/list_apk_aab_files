package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// Prompt the user to enter a directory path
	fmt.Print("Enter directory path: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	dirPath := scanner.Text()

	// Validate and list APK and AAB files
	totalSizeMB, filesToDelete, err := listApkAabFiles(dirPath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(filesToDelete) > 0 {
		// Ask the user for confirmation before cleanup
		fmt.Printf("\nTotal size of APK and AAB files: %.2f MB\n", totalSizeMB)
		fmt.Print("Do you want to clean up these files? (y/n): ")
		scanner.Scan()
		answer := scanner.Text()

		if answer == "y" || answer == "Y" {
			// Clean up the files
			if err := cleanupFiles(filesToDelete); err != nil {
				fmt.Println("Error cleaning up files:", err)
				return
			}
			fmt.Println("Files cleaned up successfully.")
		} else {
			fmt.Println("Files were not cleaned up.")
		}
	} else {
		fmt.Println("\nNo APK or AAB files found.")
	}

	fmt.Print("\nPress any key to exit...")
	scanner.Scan()
}

func listApkAabFiles(dirPath string) (float64, []string, error) {
	var totalSizeMB float64
	var filesToDelete []string

	// Validate if the provided path is a directory
	fileInfo, err := os.Stat(dirPath)
	if err != nil {
		return 0, nil, err
	}
	if !fileInfo.IsDir() {
		return 0, nil, fmt.Errorf("%s is not a directory", dirPath)
	}

	// Walk through the directory recursively
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Handle "Access is denied" error by skipping the folder
			if os.IsPermission(err) {
				fmt.Printf("Skipping folder: %s (Access is denied)\n", path)
				return filepath.SkipDir
			}
			return err
		}

		// Check if the file is either APK or AAB
		if isApkOrAabFile(path) {
			sizeMB := float64(info.Size()) / (1024 * 1024)
			fmt.Printf("%s - %.2f MB\n", path, sizeMB)
			totalSizeMB += sizeMB
			filesToDelete = append(filesToDelete, path)
		}

		return nil
	})

	return totalSizeMB, filesToDelete, err
}

func isApkOrAabFile(filePath string) bool {
	ext := filepath.Ext(filePath)
	return ext == ".apk" || ext == ".aab"
}

func cleanupFiles(filesToDelete []string) error {
	for _, file := range filesToDelete {
		if err := os.Remove(file); err != nil {
			return err
		}
	}
	return nil
}

package main

import (
	"fmt"
	"os"
)

func s() {
	// Specify the path to the folder you want to delete
	targetFolder := "0TK2YIli7K1leLovkQiNik"

	// Call the function to delete the folder and its contents recursively
	err := deleteFolder(targetFolder)
	if err != nil {
		fmt.Printf("Error deleting folder: %v\n", err)
	} else {
		fmt.Printf("Folder %s and its contents deleted successfully\n", targetFolder)
	}
}

func deleteFolder(targetFolder string) error {
	return removePath(targetFolder)
}

func removePath(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		fmt.Printf("Error removing path %s: %v\n", path, err)
	}
	return err
}

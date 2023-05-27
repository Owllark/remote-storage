package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"remote-storage/server/file_system"
)

func f(w http.ResponseWriter, r *http.Request) {

	_, err := file_system.CreateFile("", "test.txt")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error creating file"))
		return
	}

	files, err := getFilesInDirectory("storage/")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(200)
	w.Write([]byte(files[0]))
}

func d(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(200)
	w.Write([]byte("Hello, world!"))
}

func getFilesInDirectory(directoryPath string) ([]string, error) {
	var fileList []string

	// Read the directory contents
	entries, err := os.ReadDir(directoryPath)
	if err != nil {
		return nil, err
	}

	// Iterate over the entries
	for _, entry := range entries {
		filePath := filepath.Join(directoryPath, entry.Name())

		// Check if the current entry is a file
		if !entry.IsDir() {
			fileList = append(fileList, filePath)
		} else {
			// Recursively call the function for subdirectories
			subDirFiles, err := getFilesInDirectory(filePath)
			if err != nil {
				log.Printf("Error reading subdirectory '%s': %s\n", filePath, err)
				continue
			}
			fileList = append(fileList, subDirFiles...)

		}

	}

	return fileList, nil
}

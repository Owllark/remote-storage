package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"remote-storage/schemas"
)

func f(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(200)
	w.Write([]byte("Hello, world!"))
}

func Cd(w http.ResponseWriter, r *http.Request) {
	var request schemas.CdRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	curPath, err := fs.Cd(request.Path)

	w.WriteHeader(200)
	var response schemas.CdResponse
	response.Path = curPath
	if err != nil {
		response.Message = err.Error()
	}
	body, err := json.Marshal(response)
	w.Write([]byte(body))
}

func Ls(w http.ResponseWriter, r *http.Request) {
	var request schemas.LsRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	output, err := fs.Ls(request.DirPath)

	w.WriteHeader(200)
	var response schemas.LsResponse
	response.CommandOutput = output
	if err != nil {
		response.Message = err.Error()
	}
	body, err := json.Marshal(response)
	w.Write([]byte(body))
}

func MkDir(w http.ResponseWriter, r *http.Request) {

	var request schemas.MkDirRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	curPath, err := fs.MkDir(request.Path, request.DirName)

	w.WriteHeader(200)
	var response schemas.MkDirResponse
	response.Path = curPath
	if err != nil {
		response.Message = err.Error()
	}
	body, err := json.Marshal(response)
	w.Write([]byte(body))
}

func Rename(w http.ResponseWriter, r *http.Request) {
	var request schemas.RenameRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	path, err := fs.RenameCmd(request.DirPath, request.OldName, request.NewName)

	w.WriteHeader(200)
	var response schemas.RenameResponse
	response.Path = path
	if err != nil {
		response.Message = err.Error()
	}
	body, err := json.Marshal(response)
	w.Write([]byte(body))
}

func Move(w http.ResponseWriter, r *http.Request) {
	var request schemas.MoveRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	path, err := fs.CopyCmd(request.SrcDirPath, request.FileName, request.DestDirPath)

	w.WriteHeader(200)
	var response schemas.CopyResponse
	response.Path = path
	if err != nil {
		response.Message = err.Error()
	}
	body, err := json.Marshal(response)
	w.Write([]byte(body))
}

func Copy(w http.ResponseWriter, r *http.Request) {
	var request schemas.CopyRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	path, err := fs.CopyCmd(request.SrcDirPath, request.FileName, request.DestDirPath)

	w.WriteHeader(200)
	var response schemas.CopyResponse
	response.Path = path
	if err != nil {
		response.Message = err.Error()
	}
	body, err := json.Marshal(response)
	w.Write([]byte(body))
}

func Delete(w http.ResponseWriter, r *http.Request) {
	var request schemas.DeleteRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	path, err := fs.DeleteCmd(request.DirPath, request.FileName)

	w.WriteHeader(200)
	var response schemas.CopyResponse
	response.Path = path
	if err != nil {
		response.Message = err.Error()
	}
	body, err := json.Marshal(response)
	w.Write([]byte(body))
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

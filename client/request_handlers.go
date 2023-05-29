package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"remote-storage/schemas"
)

func CreateMkDirRequest(path, dirname string) (*http.Request, error) {
	var body = schemas.MkDirRequest{Path: path, DirName: dirname}
	request, err := createRequest(body, "mkdir")
	return request, err
}

func CreateCdRequest(location string) (*http.Request, error) {
	var body = schemas.CdRequest{Path: location}
	request, err := createRequest(body, "cd")
	return request, err
}

func CreateLsRequest(location string) (*http.Request, error) {
	var body = schemas.LsRequest{DirPath: location}
	request, err := createRequest(body, "ls")
	return request, err
}

func CreateRenameRequest(location, oldName, newName string) (*http.Request, error) {
	var body = schemas.RenameRequest{
		DirPath: location,
		OldName: oldName,
		NewName: newName,
	}
	request, err := createRequest(body, "rename")
	return request, err
}

func CreateMoveRequest(srcLocation, filename, destLocation string) (*http.Request, error) {
	var body = schemas.MoveRequest{
		SrcDirPath:  srcLocation,
		FileName:    filename,
		DestDirPath: destLocation,
	}
	request, err := createRequest(body, "move")
	return request, err
}

func CreateCopyRequest(srcLocation, filename, destLocation string) (*http.Request, error) {
	var body = schemas.CopyRequest{
		SrcDirPath:  srcLocation,
		FileName:    filename,
		DestDirPath: destLocation,
	}
	request, err := createRequest(body, "delete")
	return request, err
}

func CreateDeleteRequest(location, filename string) (*http.Request, error) {
	var body = schemas.DeleteRequest{
		DirPath:  location,
		FileName: filename,
	}
	request, err := createRequest(body, "delete")
	return request, err
}

func createRequest(body any, api string) (*http.Request, error) {
	var request *http.Request
	var err error
	var bodyBytes []byte
	bodyBytes, err = json.Marshal(body)
	if err != nil {
		return request, err
	}
	request, err = http.NewRequest("POST", serverUrl+api, bytes.NewReader(bodyBytes))
	request.Header.Set("Content-Type", "application/json")
	return request, err
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"remote-storage/schemas"
	"strconv"
)

func CreateMkDirRequest(path, dirname string) (*http.Request, error) {
	var body = schemas.MkDirRequest{Path: path, DirName: dirname}
	request, err := createRequestJson(body, "mkdir")
	return request, err
}

func CreateCdRequest(location string) (*http.Request, error) {
	var body = schemas.CdRequest{Path: location}
	request, err := createRequestJson(body, "cd")
	return request, err
}

func CreateLsRequest(location string) (*http.Request, error) {
	var body = schemas.LsRequest{DirPath: location}
	request, err := createRequestJson(body, "ls")
	return request, err
}

func CreateRenameRequest(location, oldName, newName string) (*http.Request, error) {
	var body = schemas.RenameRequest{
		DirPath: location,
		OldName: oldName,
		NewName: newName,
	}
	request, err := createRequestJson(body, "rename")
	return request, err
}

func CreateMoveRequest(srcLocation, filename, destLocation string) (*http.Request, error) {
	var body = schemas.MoveRequest{
		SrcDirPath:  srcLocation,
		FileName:    filename,
		DestDirPath: destLocation,
	}
	request, err := createRequestJson(body, "move")
	return request, err
}

func CreateCopyRequest(srcLocation, filename, destLocation string) (*http.Request, error) {
	var body = schemas.CopyRequest{
		SrcDirPath:  srcLocation,
		FileName:    filename,
		DestDirPath: destLocation,
	}
	request, err := createRequestJson(body, "copy")
	return request, err
}

func CreateDeleteRequest(location, filename string) (*http.Request, error) {
	var body = schemas.DeleteRequest{
		DirPath:  location,
		FileName: filename,
	}
	request, err := createRequestJson(body, "delete")
	return request, err
}

func CreateStartUploadRequest(location, filename string, chunksNum int) (*http.Request, error) {
	var body = schemas.StartUploadRequest{
		Location:  location,
		FileName:  filename,
		ChunksNum: chunksNum,
	}
	request, err := createRequestJson(body, "upload")
	return request, err
}

func CreateUploadChunkRequest(chunk []byte, id int) (*http.Request, error) {
	body := chunk
	params := url.Values{}
	params.Set("id", strconv.Itoa(id))
	requestURL := fmt.Sprintf("%s?%s", "upload/chunk", params.Encode())
	request, err := createRequestPlain(body, requestURL)
	return request, err
}

func CreateCompleteUploadRequest() (*http.Request, error) {
	var body = schemas.CompleteUploadResponse{}
	request, err := createRequestJson(body, "upload/completed")
	return request, err
}

func CreateStartDownloadRequest(location, filename string) (*http.Request, error) {
	var body = schemas.StartDownloadRequest{
		Location: location,
		FileName: filename,
	}
	request, err := createRequestJson(body, "download")
	return request, err
}

func CreateDownloadChunkRequest(id int) (*http.Request, error) {
	body := ""
	params := url.Values{}
	params.Set("id", strconv.Itoa(id))
	requestURL := fmt.Sprintf("%s?%s", "download/chunk", params.Encode())
	request, err := createRequestPlain([]byte(body), requestURL)
	return request, err
}

func CreateAuthenticationRequest(name, password string) (*http.Request, error) {
	var body = schemas.AuthenticateRequest{
		Name:     name,
		Password: password,
	}
	request, err := createRequestJson(body, "authenticate")
	return request, err
}

func createRequestJson(body any, api string) (*http.Request, error) {
	var request *http.Request
	var err error
	var bodyBytes []byte
	bodyBytes, err = json.Marshal(body)
	if err != nil {
		return request, err
	}
	request, err = http.NewRequest("POST", serverUrl+api, bytes.NewReader(bodyBytes))
	request.Header.Set("Content-Type", "application/json")
	cookie := &http.Cookie{
		Name:  "token",
		Value: jwtToken,
	}
	request.AddCookie(cookie)
	return request, err
}

func createRequestPlain(bodyBytes []byte, api string) (*http.Request, error) {
	var request *http.Request
	var err error
	request, err = http.NewRequest("POST", serverUrl+api, bytes.NewReader(bodyBytes))
	request.Header.Set("Content-Type", "text/plain")
	cookie := &http.Cookie{
		Name:  "token",
		Value: jwtToken,
	}
	request.AddCookie(cookie)
	return request, err
}

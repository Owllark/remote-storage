package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"remote-storage/common"
	"strconv"
)

func CreateMkDirRequest(path, dirname string) (*http.Request, error) {
	var body = common.MkDirRequest{Path: path, DirName: dirname}
	request, err := createPostRequestJson(body, "mkdir")
	return request, err
}

func CreateCdRequest(location string) (*http.Request, error) {
	var body = common.CdRequest{Path: location}
	request, err := createPostRequestJson(body, "cd")
	return request, err
}

func CreateLsRequest(location string) (*http.Request, error) {
	var body = common.LsRequest{DirPath: location}
	request, err := createPostRequestJson(body, "ls")
	return request, err
}

func CreateRenameRequest(location, oldName, newName string) (*http.Request, error) {
	var body = common.RenameRequest{
		DirPath: location,
		OldName: oldName,
		NewName: newName,
	}
	request, err := createPostRequestJson(body, "rename")
	return request, err
}

func CreateMoveRequest(srcLocation, filename, destLocation string) (*http.Request, error) {
	var body = common.MoveRequest{
		SrcDirPath:  srcLocation,
		FileName:    filename,
		DestDirPath: destLocation,
	}
	request, err := createPostRequestJson(body, "move")
	return request, err
}

func CreateCopyRequest(srcLocation, filename, destLocation string) (*http.Request, error) {
	var body = common.CopyRequest{
		SrcDirPath:  srcLocation,
		FileName:    filename,
		DestDirPath: destLocation,
	}
	request, err := createPostRequestJson(body, "copy")
	return request, err
}

func CreateDeleteRequest(location, filename string) (*http.Request, error) {
	var body = common.DeleteRequest{
		DirPath:  location,
		FileName: filename,
	}
	request, err := createPostRequestJson(body, "delete")
	return request, err
}

func CreateStartUploadRequest(location, filename string, chunksNum int) (*http.Request, error) {
	var body = common.StartUploadRequest{
		Location:  location,
		FileName:  filename,
		ChunksNum: chunksNum,
	}
	request, err := createPostRequestJson(body, "upload")
	return request, err
}

func CreateUploadChunkRequest(chunk []byte, id int) (*http.Request, error) {
	var body = common.UploadChunkRequest{
		Id:   id,
		Data: chunk,
	}
	request, err := createPostRequestJson(body, "upload/chunk")
	q := request.URL.Query()
	q.Add("id", strconv.Itoa(id))
	request.URL.RawQuery = q.Encode()
	return request, err
}

func CreateCompleteUploadRequest() (*http.Request, error) {
	var body = common.CompleteUploadResponse{}
	request, err := createPostRequestJson(body, "upload/completed")
	return request, err
}

func CreateStartDownloadRequest(location, filename string) (*http.Request, error) {
	var body = common.StartDownloadRequest{
		Location: location,
		FileName: filename,
	}
	request, err := createPostRequestJson(body, "download")
	return request, err
}

func CreateDownloadChunkRequest(id int) (*http.Request, error) {
	var body = ""
	request, err := createPostRequestJson(body, "download/chunk")
	q := request.URL.Query()
	q.Add("id", strconv.Itoa(id))
	request.URL.RawQuery = q.Encode()
	return request, err
}

func CreateAuthenticationRequest(name, password string) (*http.Request, error) {
	var body = common.AuthenticateRequest{
		Name:     name,
		Password: password,
	}
	request, err := createPostRequestJson(body, "authenticate")
	return request, err
}

func CreateGetFileSystemStateRequest() (*http.Request, error) {
	request, err := createGetRequest("state")
	return request, err
}

func CreateLogoutRequest() (*http.Request, error) {
	request, err := createGetRequest("logout")
	return request, err
}

func CreateRefreshRequest() (*http.Request, error) {
	request, err := createGetRequest("refresh")
	return request, err
}

func createGetRequest(api string) (*http.Request, error) {
	request, err := http.NewRequest("GET", serverUrl+api, nil)
	cookie := &http.Cookie{
		Name:  "token",
		Value: jwtToken,
	}
	request.AddCookie(cookie)
	return request, err
}
func createPostRequestJson(body any, api string) (*http.Request, error) {
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

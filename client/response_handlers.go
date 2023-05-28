package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"remote-storage/schemas"
)

func HandleMkDirResponse(path, dirname string) (*http.Request, error) {
	var request *http.Request
	var err error
	var body = schemas.MkDirRequest{Path: path, DirName: dirname}
	var bodyBytes []byte
	bodyBytes, err = json.Marshal(body)
	if err != nil {
		return request, err
	}
	request, err = http.NewRequest("POST", serverUrl, bytes.NewReader(bodyBytes))
	request.Header.Set("Content-Type", "application/json")
	return request, err
}

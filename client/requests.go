package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"remote-storage/client/http_client"
	"remote-storage/common"
)

func MkDir(client http_client.Client, location, dirname string) (common.MkDirResponse, error) {
	var responseInf common.MkDirResponse
	request, err := CreateMkDirRequest(location, dirname)
	if err != nil {
		return responseInf, err
	}
	err = DoRequest(client, request, &responseInf)
	return responseInf, err
}

func Rename(client http_client.Client, location, oldName, newName string) (common.RenameResponse, error) {
	var responseInf common.RenameResponse
	request, err := CreateRenameRequest(location, oldName, newName)
	if err != nil {
		return responseInf, err
	}
	err = DoRequest(client, request, &responseInf)
	return responseInf, err
}

func Move(client http_client.Client, srcLocation, filename, destLocation string) (common.MoveResponse, error) {
	var responseInf common.MoveResponse
	request, err := CreateMoveRequest(srcLocation, filename, destLocation)
	if err != nil {
		return responseInf, err
	}
	err = DoRequest(client, request, &responseInf)
	return responseInf, err
}

func Copy(client http_client.Client, srcLocation, filename, destLocation string) (common.CopyResponse, error) {
	var responseInf common.CopyResponse
	request, err := CreateCopyRequest(srcLocation, filename, destLocation)
	if err != nil {
		return responseInf, err
	}
	err = DoRequest(client, request, &responseInf)
	return responseInf, err
}

func Delete(client http_client.Client, location, filename string) (common.DeleteResponse, error) {
	var responseInf common.DeleteResponse
	request, err := CreateDeleteRequest(location, filename)
	if err != nil {
		return responseInf, err
	}
	err = DoRequest(client, request, &responseInf)
	return responseInf, err
}

func StartUpload(client http_client.Client, location, filename string, chunksNum int) (common.StartUploadResponse, error) {
	var responseInf common.StartUploadResponse
	request, err := CreateStartUploadRequest(location, filename, chunksNum)
	if err != nil {
		return responseInf, err
	}
	err = DoRequest(client, request, &responseInf)
	return responseInf, err
}

func UploadChunk(client http_client.Client, chunk []byte, id int) (common.UploadChunkResponse, error) {
	var responseInf common.UploadChunkResponse
	request, err := CreateUploadChunkRequest(chunk, id)
	if err != nil {
		return responseInf, err
	}
	err = DoRequest(client, request, &responseInf)
	return responseInf, err
}

func CompleteUpload(client http_client.Client, chunk []byte, id int) (common.CompleteUploadResponse, error) {
	var responseInf common.CompleteUploadResponse
	request, err := CreateCompleteUploadRequest()
	if err != nil {
		return responseInf, err
	}
	err = DoRequest(client, request, &responseInf)
	return responseInf, err
}

func StartDownload(client http_client.Client, location, filename string, chunksNum int) (common.StartDownloadResponse, error) {
	var responseInf common.StartDownloadResponse
	request, err := CreateStartDownloadRequest(location, filename)
	if err != nil {
		return responseInf, err
	}
	err = DoRequest(client, request, &responseInf)
	return responseInf, err
}

func DownloadChunk(client http_client.Client, id int) (common.DownloadChunkResponse, error) {
	var responseInf common.DownloadChunkResponse
	request, err := CreateDownloadChunkRequest(id)
	if err != nil {
		return responseInf, err
	}
	err = DoRequest(client, request, &responseInf)
	return responseInf, err
}

func GetFileSystemState(client http_client.Client) (common.GetStateResponse, error) {
	var responseInf common.GetStateResponse
	request, err := CreateGetFileSystemStateRequest()
	if err != nil {
		return responseInf, err
	}
	err = DoRequest(client, request, &responseInf)
	return responseInf, err
}

func DoRequest[scheme any](client http_client.Client, request *http.Request, responseInf *scheme) error {
	response, err := client.DoRequest(request)
	if err != nil {
		return errors.New("Problem with server connection")
	}
	if response.StatusCode != 200 {
		switch response.StatusCode {
		case 401:
			{
				return errors.New("Error: Unauthorized")
			}
		case 404:
			{
				return errors.New("Error: Not found")
			}
		case 500:
			{
				return errors.New("Server error occurred")
			}
		}

	}
	body, _ := io.ReadAll(response.Body)
	err = json.Unmarshal(body, &responseInf)

	return nil
}

func Logout(client http_client.Client) error {
	request, err := CreateLogoutRequest()
	if err != nil {
		return err
	}
	response, err := client.DoRequest(request)
	if err != nil {
		return errors.New("Problem with server connection")
	}
	if response.StatusCode != 200 {
		return errors.New("Server error occurred")
	}
	return nil
}

func UploadFile(client http_client.Client, path, filename, remoteLocation string) error {

	chunks, err := DivideFileIntoChunks(path, CHUNK_SIZE)
	if err != nil {
		return err
	}

	chunksUploaded := 0
	chunksAmount := len(chunks)
	_, err = StartUpload(client, remoteLocation, filename, chunksAmount)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	for i, chunk := range chunks {
		request, err := CreateUploadChunkRequest(chunk, i)
		if err != nil {
			return err
		}
		response, err := client.DoRequest(request)
		if err == nil && response.StatusCode == 200 {
			chunksUploaded++
		}
		PrintProgressBar(chunksUploaded, chunksAmount)
	}
	var allChunksAccepted bool
	for !allChunksAccepted {
		request, err := CreateCompleteUploadRequest()
		if err != nil {
			return err
		}
		response, err := client.DoRequest(request)
		if err != nil {
			return errors.New("Problem with server connection")
		}
		var responseInf common.CompleteUploadResponse
		body, _ := io.ReadAll(response.Body)
		json.Unmarshal(body, &responseInf)
		if len(responseInf.MissedChunks) == 0 {
			allChunksAccepted = true
			PrintProgressBar(chunksUploaded, chunksAmount)
			fmt.Println()
			fmt.Println(responseInf.Message)
		} else {
			for _, i := range responseInf.MissedChunks {
				request, err := CreateUploadChunkRequest(chunks[i], i)
				if err != nil {
					return err
				}
				response, err := client.DoRequest(request)
				if err == nil && response.StatusCode == 200 {
					chunksUploaded++
				}
				PrintProgressBar(chunksUploaded, chunksAmount)
			}
		}
	}
	return nil
}

func DownloadFile(client http_client.Client, downloadLocation, filename, remoteLocation string) error {
	request, err := CreateStartDownloadRequest(remoteLocation, filename)
	if err != nil {
		return err
	}
	response, err := client.DoRequest(request)
	if err != nil {
		return errors.New("Problem with server connection")
	}
	if response.StatusCode != 200 {
		return errors.New("error downloading file")
	}
	var responseInf common.StartDownloadResponse
	body, _ := io.ReadAll(response.Body)
	json.Unmarshal(body, &responseInf)
	file, err := os.OpenFile(downloadLocation+filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	chunkNum := responseInf.ChunksNum
	for i := 0; i < chunkNum; i++ {
		isChunkDownloaded := false
		request, err := CreateDownloadChunkRequest(i)
		if err != nil {
			return err
		}
		for !isChunkDownloaded {
			response, err := client.DoRequest(request)
			if err != nil {
				continue
			}
			if response.StatusCode == 200 {
				isChunkDownloaded = true
				var responseInf common.DownloadChunkResponse
				body, _ := io.ReadAll(response.Body)
				json.Unmarshal(body, &responseInf)
				file.Write(responseInf.Data)

			}
		}
		PrintProgressBar(i+1, chunkNum)

	}
	file.Close()
	return nil
}

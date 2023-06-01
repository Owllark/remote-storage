package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"remote-storage/client/http_client"
	"remote-storage/schemas"
	"strings"
)

const (
	barSize       = 100
	serverUrl     = "http://localhost:8080/"
	pathSeparator = string(os.PathSeparator)
	CHUNK_SIZE    = 64 * 1024
)

type State struct {
	curPath string
}

var state State

var username string

var jwtToken string

func main() {
	client := http_client.NewHttpClient("")
	reader := bufio.NewReader(os.Stdin)
	var isOver = false
	state.curPath = ""
	isAuthenticated := false
	for !isAuthenticated {
		fmt.Println("Enter your name:")
		name, _ := reader.ReadString('\n')
		fmt.Println("Enter your password:")
		password, _ := reader.ReadString('\n')
		name = strings.TrimSuffix(name, "\r\n")
		password = strings.TrimSuffix(password, "\r\n")
		request, _ := CreateAuthenticationRequest(name, password)
		response, err := client.DoRequest(request)
		if err != nil {
			fmt.Println("Problem with server connection")
			continue
		}
		if response.StatusCode == 401 {
			fmt.Println("Authentication failed")
			continue
		} else if response.StatusCode != 200 {
			fmt.Println("Server error occurred")
			continue
		}
		isTokenReceived := false
		for _, cookie := range response.Cookies() {
			if cookie.Name == "token" {
				jwtToken = cookie.Value
				isTokenReceived = true
			}
		}
		if !isTokenReceived {
			fmt.Println("Error accepting necessary cookies")
			continue
		}
		isAuthenticated = true
		var responseInf schemas.AuthenticateResponse
		body, _ := io.ReadAll(response.Body)
		json.Unmarshal(body, &responseInf)
		username = name
		fmt.Println("Authenticated successfully")
	}
	for !isOver {
		fmt.Print(fmt.Sprintf("%s%s>", username, state.curPath))
		var req string
		req, _ = reader.ReadString('\n')
		req = req[0 : len(req)-2]
		arguments := strings.Split(req, " ")
		switch arguments[0] {
		case "":
			{
				continue
			}
		case "mkdir":
			{
				if len(arguments) < 2 {
					fmt.Println("Not enough arguments")
					continue
				} else if len(arguments) > 2 {
					fmt.Println("Too much arguments")
					continue
				}
				location, dirname, err := getLocationAndFileName(arguments[1])
				if err != nil {
					fmt.Println(err)
					continue
				}
				request, err := CreateMkDirRequest(location, dirname)
				if err != nil {
					fmt.Println(err)
					continue
				}
				response, err := client.DoRequest(request)
				if err != nil {
					fmt.Println("Problem with server connection")
					continue
				}
				if response.StatusCode != 200 {
					fmt.Println("Server error occurred")
					continue
				}
				body, _ := io.ReadAll(response.Body)
				var responseInf schemas.MkDirResponse
				json.Unmarshal(body, &responseInf)
				if responseInf.Message != "" {
					fmt.Println(responseInf.Message)
				}
			}
		case "cd":
			{
				if len(arguments) < 2 {
					fmt.Println("Not enough arguments")
					continue
				} else if len(arguments) > 2 {
					fmt.Println("Too much arguments")
					continue
				}
				location := arguments[1]
				request, err := CreateCdRequest(location)
				if err != nil {
					fmt.Println(err)
					continue
				}
				response, err := client.DoRequest(request)
				if err != nil {
					fmt.Println("Problem with server connection")
					continue
				}
				if response.StatusCode != 200 {
					fmt.Println("Server error occurred")
					continue
				}
				body, _ := io.ReadAll(response.Body)
				var responseInf schemas.CdResponse
				json.Unmarshal(body, &responseInf)
				state.curPath = responseInf.Path
				if responseInf.Message != "" {
					fmt.Println(responseInf.Message)
				}
			}
		case "ls":
			{
				if len(arguments) < 2 {
					fmt.Println("Not enough arguments")
					continue
				} else if len(arguments) > 2 {
					fmt.Println("Too much arguments")
					continue
				}
				location := arguments[1]
				request, err := CreateLsRequest(location)
				if err != nil {
					fmt.Println(err)
					continue
				}
				response, err := client.DoRequest(request)
				if err != nil {
					fmt.Println("Problem with server connection")
					continue
				}
				if response.StatusCode != 200 {
					fmt.Println("Server error occurred")
					continue
				}
				body, _ := io.ReadAll(response.Body)
				var responseInf schemas.LsResponse
				json.Unmarshal(body, &responseInf)
				if responseInf.Message != "" {
					fmt.Println(responseInf.Message)
				}
				output := responseInf.CommandOutput
				for _, s := range output {
					fmt.Println(s)
				}

			}
		case "rename":
			{
				if len(arguments) < 3 {
					fmt.Println("Not enough arguments")
					continue
				} else if len(arguments) > 3 {
					fmt.Println("Too much arguments")
					continue
				}
				location, oldName, err := getLocationAndFileName(arguments[1])
				newName := arguments[2]
				if err != nil {
					fmt.Println(err)
					continue
				}
				request, err := CreateRenameRequest(location, oldName, newName)
				if err != nil {
					fmt.Println(err)
					continue
				}
				response, err := client.DoRequest(request)
				if err != nil {
					fmt.Println("Problem with server connection")
					continue
				}
				if response.StatusCode != 200 {
					fmt.Println("Server error occurred")
					continue
				}
				body, _ := io.ReadAll(response.Body)
				var responseInf schemas.RenameResponse
				json.Unmarshal(body, &responseInf)
				if responseInf.Message != "" {
					fmt.Println(responseInf.Message)
				}
			}
		case "move":
			{
				if len(arguments) < 3 {
					fmt.Println("Not enough arguments")
					continue
				} else if len(arguments) > 3 {
					fmt.Println("Too much arguments")
					continue
				}
				oldLocation, filename, err := getLocationAndFileName(arguments[1])
				newLocation := arguments[2]
				if err != nil {
					fmt.Println(err)
					continue
				}
				request, err := CreateMoveRequest(oldLocation, filename, newLocation)
				if err != nil {
					fmt.Println(err)
					continue
				}
				response, err := client.DoRequest(request)
				if err != nil {
					fmt.Println("Problem with server connection")
					continue
				}
				if response.StatusCode != 200 {
					fmt.Println("Server error occurred")
					continue
				}
				body, _ := io.ReadAll(response.Body)
				var responseInf schemas.MoveResponse
				json.Unmarshal(body, &responseInf)
				if responseInf.Message != "" {
					fmt.Println(responseInf.Message)
				}
			}
		case "copy":
			{
				if len(arguments) < 3 {
					fmt.Println("Not enough arguments")
					continue
				} else if len(arguments) > 3 {
					fmt.Println("Too much arguments")
					continue
				}
				srcLocation, filename, err := getLocationAndFileName(arguments[1])
				destLocation := arguments[2]
				if err != nil {
					fmt.Println(err)
					continue
				}
				request, err := CreateCopyRequest(srcLocation, filename, destLocation)
				if err != nil {
					fmt.Println(err)
					continue
				}
				response, err := client.DoRequest(request)
				if err != nil {
					fmt.Println("Problem with server connection")
					continue
				}
				if response.StatusCode != 200 {
					fmt.Println("Server error occurred")
					continue
				}
				body, _ := io.ReadAll(response.Body)
				var responseInf schemas.CopyResponse
				json.Unmarshal(body, &responseInf)
				if responseInf.Message != "" {
					fmt.Println(responseInf.Message)
				}
			}
		case "delete":
			{
				if len(arguments) < 2 {
					fmt.Println("Not enough arguments")
					continue
				} else if len(arguments) > 2 {
					fmt.Println("Too much arguments")
					continue
				}
				location, filename, err := getLocationAndFileName(arguments[1])
				if err != nil {
					fmt.Println(err)
					continue
				}
				request, err := CreateDeleteRequest(location, filename)
				if err != nil {
					fmt.Println(err)
					continue
				}
				response, err := client.DoRequest(request)
				if err != nil {
					fmt.Println("Problem with server connection")
					continue
				}
				body, _ := io.ReadAll(response.Body)
				var responseInf schemas.DeleteResponse
				json.Unmarshal(body, &responseInf)
				if responseInf.Message != "" {
					fmt.Println(responseInf.Message)
				}
			}
		case "upload":
			{
				if len(arguments) < 2 {
					fmt.Println("Not enough arguments")
					continue
				} else if len(arguments) > 3 {
					fmt.Println("Too much arguments")
					continue
				}
				path := arguments[1]
				chunks, err := DivideFileIntoChunks(path, CHUNK_SIZE)
				if err != nil {
					fmt.Println(err)
					continue
				}
				_, filename, err := getLocationAndFileName(path)
				if err != nil {
					fmt.Println(err)
					continue
				}
				var location string
				if len(arguments) == 2 {
					location = "/"
				} else {
					location = arguments[2]
				}

				chunksUploaded := 0
				chunksAmount := len(chunks)
				request, err := CreateStartUploadRequest(location, filename, len(chunks))
				if err != nil {
					fmt.Println(err)
					continue
				}
				_, err = client.DoRequest(request)
				if err != nil {
					fmt.Println("Problem with server connection")
					continue
				}
				for i, chunk := range chunks {
					request, err := CreateUploadChunkRequest(chunk, i)
					if err != nil {
						fmt.Println(err)
						continue
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
						fmt.Println(err)
						continue
					}
					response, err := client.DoRequest(request)
					if err != nil {
						fmt.Println("Problem with server connection")
						continue
					}
					var responseInf schemas.CompleteUploadResponse
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
								fmt.Println(err)
								continue
							}
							response, err := client.DoRequest(request)
							if err == nil && response.StatusCode == 200 {
								chunksUploaded++
							}
							PrintProgressBar(chunksUploaded, chunksAmount)
						}
					}
				}

			}
		case "download":
			{
				if len(arguments) < 2 {
					fmt.Println("Not enough arguments")
					continue
				} else if len(arguments) > 3 {
					fmt.Println("Too much arguments")
					continue
				}
				var downloadLocation string
				if len(arguments) == 2 {
					downloadLocation = ""
				} else {
					downloadLocation = arguments[2]
				}
				location, filename, err := getLocationAndFileName(arguments[1])
				if err != nil {
					fmt.Println(err)
					continue
				}
				request, err := CreateStartDownloadRequest(location, filename)
				if err != nil {
					fmt.Println(err)
					continue
				}
				response, err := client.DoRequest(request)
				if err != nil {
					fmt.Println("Problem with server connection")
					continue
				}
				if response.StatusCode != 200 {
					fmt.Println("error downloading file")
					continue
				}
				var responseInf schemas.StartDownloadResponse
				body, _ := io.ReadAll(response.Body)
				json.Unmarshal(body, &responseInf)
				file, err := os.OpenFile(downloadLocation+filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
				chunkNum := responseInf.ChunksNum
				for i := 0; i < chunkNum; i++ {
					isChunkDownloaded := false
					request, err := CreateDownloadChunkRequest(i)
					if err != nil {
						fmt.Println(err)
						break
					}
					for !isChunkDownloaded {
						response, err := client.DoRequest(request)
						if err != nil {
							fmt.Println("Problem with server connection")
							continue
						}
						if response.StatusCode == 200 {
							isChunkDownloaded = true
							data, _ := io.ReadAll(response.Body)
							file.Write(data)

						}
					}

				}
				file.Close()
			}
		default:
			{
				fmt.Println(fmt.Sprintf("WRONG REQUEST: %s UNDEFINED", arguments[0]))
			}

		}
	}

}

func PrintProgressBar(n, total int) {
	var progress int
	var percent int
	if total == 0 {
		progress = barSize
		percent = 100
	} else {
		progress = n * barSize / total
		percent = n * 100 / total
	}
	fmt.Printf("\r[%s%s] %d%%", getProgressBar(progress), getEmptyBar(barSize-progress), percent)
}

func DivideFileIntoChunks(path string, chunkSize int) ([][]byte, error) {
	var chunks [][]byte
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return chunks, err
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()

	fullChunksAmount := int(fileSize) / chunkSize
	lastChunkSize := int(fileSize) % chunkSize
	chunkAmount := fullChunksAmount
	if lastChunkSize > 0 {
		chunkAmount++
	}
	chunks = make([][]byte, chunkAmount)

	for i := 0; i < fullChunksAmount; i++ {
		chunk := make([]byte, chunkSize)
		_, err := file.Read(chunk)
		if err != nil {
			return chunks, err
		}
		chunks[i] = chunk
	}

	if lastChunkSize > 0 {
		lastChunk := make([]byte, lastChunkSize)
		_, err := file.Read(lastChunk)
		if err != nil {
			return chunks, err
		}
		chunks[chunkAmount-1] = lastChunk
	}

	return chunks, err
}

func getLocationAndFileName(path string) (string, string, error) {
	var location, filename string
	if string(path[len(path)-1]) == pathSeparator {
		return location, filename, errors.New("unexpected separator symbol at the end of the path")
	}
	dirList := splitPath(path)
	if string(path[0]) == pathSeparator {
		location = pathSeparator + joinPath(dirList[:len(dirList)-1])
		filename = dirList[len(dirList)-1]
	} else {
		location = joinPath(dirList[:len(dirList)-1])
		filename = dirList[len(dirList)-1]
	}
	return location, filename, nil
}

func joinPath(path []string) string {
	var res string
	for _, s := range path {
		if s == "" {
			continue
		}
		res += s + pathSeparator
	}
	return res
}

func splitPath(path string) []string {
	res := strings.Split(path, pathSeparator)
	return res
}

func getProgressBar(progress int) string {
	bar := ""
	for i := 0; i < progress; i++ {
		bar += "="
	}
	return bar
}

func getEmptyBar(emptyCount int) string {
	empty := ""
	for i := 0; i < emptyCount; i++ {
		empty += " "
	}
	return empty
}

package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"remote-storage/client/file_system"
	"remote-storage/client/http_client"
	"remote-storage/common"
	"strings"
	"time"
)

const (
	barSize       = 100
	serverUrl     = "http://localhost:8080/"
	pathSeparator = string(os.PathSeparator)
	CHUNK_SIZE    = 64 * 1024
)

type State struct {
	fs *file_system.FileSystemBrowser
}

var state State

var username string

var jwtToken string

var tokenExpires time.Time

var tokenDuration time.Duration

var client = http_client.NewHttpClient("")

func main() {
	reader := bufio.NewReader(os.Stdin)
	var isOver = false
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
				tokenExpires = cookie.Expires
			}
		}
		if !isTokenReceived {
			fmt.Println("Error accepting necessary cookies")
			continue
		}

		isAuthenticated = true
		var responseInf common.AuthenticateResponse
		body, _ := io.ReadAll(response.Body)
		json.Unmarshal(body, &responseInf)
		username = name
		fmt.Println("Authenticated successfully")
		tokenDuration = tokenExpires.Sub(time.Now()) - 30*time.Second
		go RefreshTokenByTimer(tokenDuration)

	}

	state.fs = file_system.NewFileSystemBrowser(common.FileInfo{})
	UpdateFileSystemState(client)

	for !isOver {
		fmt.Print(fmt.Sprintf("%s%s>", username+pathSeparator, state.fs.GetCurrentPath()))
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
				response, err := MkDir(client, location, dirname)
				if err != nil {
					fmt.Println(err)
					continue
				}
				if response.Message != "" {
					fmt.Println(response.Message)
				}
				UpdateFileSystemState(client)

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
				err := state.fs.Cd(location)
				if err != nil {
					fmt.Println(err)
					continue
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
				dirPath := arguments[1]
				output, err := state.fs.Ls(dirPath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				for _, line := range output {
					fmt.Println(line)
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
				response, err := Rename(client, location, oldName, newName)
				if err != nil {
					fmt.Println(err)
					continue
				}
				if response.Message != "" {
					fmt.Println(response.Message)
				}
				UpdateFileSystemState(client)
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
				newLocation := getLocation(arguments[2])
				if err != nil {
					fmt.Println(err)
					continue
				}
				response, err := Move(client, oldLocation, filename, newLocation)
				if err != nil {
					fmt.Println(err)
					continue
				}
				if response.Message != "" {
					fmt.Println(response.Message)
				}
				UpdateFileSystemState(client)
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
				destLocation := getLocation(arguments[2])
				if err != nil {
					fmt.Println(err)
					continue
				}
				response, err := Copy(client, srcLocation, filename, destLocation)
				if err != nil {
					fmt.Println(err)
					continue
				}
				if response.Message != "" {
					fmt.Println(response.Message)
				}
				UpdateFileSystemState(client)
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
				response, err := Delete(client, location, filename)
				if err != nil {
					fmt.Println(err)
					continue
				}
				if response.Message != "" {
					fmt.Println(response.Message)
				}
				UpdateFileSystemState(client)
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
				_, filename, err := getLocationAndFileName(path)
				if err != nil {
					fmt.Println(err)
					continue
				}
				var location string
				if len(arguments) == 2 {
					location = state.fs.GetCurrentPath()
				} else {
					location = getLocation(arguments[2])
				}

				err = UploadFile(client, path, filename, location)
				if err != nil {
					fmt.Println(err)
					continue
				}
				UpdateFileSystemState(client)
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
				err = DownloadFile(client, downloadLocation, filename, location)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println()
				fmt.Println("downloaded successfully")
				UpdateFileSystemState(client)
			}
		case "logout":
			{
				err := Logout(client)
				if err != nil {
					fmt.Println(err)
					continue
				}
				return
			}
		default:
			{
				fmt.Println(fmt.Sprintf("WRONG REQUEST: %s UNDEFINED", arguments[0]))
			}

		}
	}

}

func RefreshTokenByTimer(duration time.Duration) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			Refresh()
		}
	}
}

func Refresh() {
	request, _ := CreateRefreshRequest()
	response, err := client.DoRequest(request)
	if err != nil || response.StatusCode != 200 {
		return
	}
	for _, cookie := range response.Cookies() {
		if cookie.Name == "token" {
			jwtToken = cookie.Value
		}
	}
}

func UpdateFileSystemState(client http_client.Client) {
	response, err := GetFileSystemState(client)
	if err != nil {
		fmt.Println("Unable to update file system state: ", err)
		return
	}
	state.fs.SetFiles(response.Info)
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
		location = joinPath(dirList[:len(dirList)-1])
		filename = dirList[len(dirList)-1]
	} else {
		location = state.fs.GetCurrentPath() + joinPath(dirList[:len(dirList)-1])
		filename = dirList[len(dirList)-1]
	}
	return location, filename, nil
}

func getLocation(path string) string {
	var location = path
	if string(path[len(path)-1]) != pathSeparator {
		location += pathSeparator
	}
	if string(path[0]) != pathSeparator {
		location = state.fs.GetCurrentPath() + location
	} else {
		location = location[len(pathSeparator):]
	}
	return location
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
	var res []string
	parts := strings.Split(path, string(os.PathSeparator))
	for _, part := range parts {
		if part != "" {
			res = append(res, part)
		}
	}

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

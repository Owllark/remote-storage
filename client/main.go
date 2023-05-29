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
	"time"
)

const (
	total         = 100
	barSize       = 20
	serverUrl     = "http://localhost:8080/"
	pathSeparator = string(os.PathSeparator)
)

type State struct {
	curPath string
}

var state State

func main() {
	client := http_client.NewHttpClient("")
	reader := bufio.NewReader(os.Stdin)
	var isOver = false
	state.curPath = ""
	for !isOver {
		fmt.Print(fmt.Sprintf("%s>", state.curPath))
		var req string
		req, _ = reader.ReadString('\n')
		req = req[0 : len(req)-2]
		arguments := strings.Split(req, " ")
		switch arguments[0] {
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
				response := client.DoRequest(request)
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
				response := client.DoRequest(request)
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
				response := client.DoRequest(request)
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
				response := client.DoRequest(request)
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
				response := client.DoRequest(request)
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
				response := client.DoRequest(request)
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
				location, fileame, err := getLocationAndFileName(arguments[1])
				if err != nil {
					fmt.Println(err)
					continue
				}
				request, err := CreateDeleteRequest(location, fileame)
				if err != nil {
					fmt.Println(err)
					continue
				}
				response := client.DoRequest(request)
				if response.StatusCode != 200 {
					fmt.Println("Server error occurred")
					continue
				}
				body, _ := io.ReadAll(response.Body)
				var responseInf schemas.DeleteResponse
				json.Unmarshal(body, &responseInf)
				if responseInf.Message != "" {
					fmt.Println(responseInf.Message)
				}
			}
		default:
			{
				fmt.Println(fmt.Sprintf("WRONG REQUEST: %s UNDEFINED", arguments[0]))
			}

		}
	}

	for i := 0; i <= total; i++ {
		progress := i * barSize / total
		fmt.Printf("\r[%s%s] %d%%", getProgressBar(progress), getEmptyBar(barSize-progress), i)
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("\nTask completed!")
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

package main

import (
	"bufio"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"remote-storage/server/src/db"
	"remote-storage/server/src/file_system"
	"remote-storage/server/src/helper"
	"strings"
)

var rootFs file_system.FileSystem

const pathSeparator = string(os.PathSeparator)

type uploadState struct {
	chunksNum      int
	tempDir        string
	fileName       string
	fileLocation   string
	chunksGotten   int
	receivedChunks []bool
}

type downloadState struct {
	chunks [][]byte
}

type Client struct {
	upload   uploadState
	download downloadState
	inf      helper.ClientInf
	fs       file_system.FileSystem
	//rootDir  string
}

var clients []Client

var database db.StorageDatabasePG

const (
	dbUser     = "postgres"
	dbHost     = "localhost"
	dbName     = "remote_storage"
	dbPassword = "password"
)

var jwtKey = []byte("secret")

const storagePath = "storage" + pathSeparator

func main() {
	/*var route = router.NewHttpRouter()
	route.AddHandler("/test", f)
	route.AddHandler("/download", f)
	route.AddHandler("/rename", Rename)
	route.AddHandler("/move", Move)
	route.AddHandler("/copy", Copy)
	route.AddHandler("/delete", RemoveAll)

	route.AddHandler("/cd", Cd)
	route.AddHandler("/mkdir", MkDir)
	route.AddHandler("/ls", Ls)
	route.AddHandler("/tree", f)
	route.AddHandler("/find", f)

	route.AddHandler("/upload", StartUploading)
	route.AddHandler("/upload_chunk/*id", UploadChunk)

	route.Listen()*/

	err := database.Connect(dbUser, dbPassword, dbName, dbHost)
	if err != nil {
		log.Fatal(err)
	}

	ClientsInit()

	rootFs.SetRootDir("storage" + pathSeparator)
	r := mux.NewRouter()
	r.HandleFunc("/test", f)
	r.HandleFunc("/rename", Rename)
	r.HandleFunc("/move", Move)
	r.HandleFunc("/copy", Copy)
	r.HandleFunc("/delete", Delete)

	r.HandleFunc("/cd", Cd)
	r.HandleFunc("/mkdir", MkDir)
	r.HandleFunc("/ls", Ls)
	r.HandleFunc("/tree", f)
	r.HandleFunc("/find", f)

	r.HandleFunc("/upload", StartUploading)
	r.HandleFunc("/upload/chunk", UploadChunk)
	r.HandleFunc("/upload/completed", UploadComplete)

	r.HandleFunc("/download", StartDownloading)
	r.HandleFunc("/download/chunk", DownloadChunk)

	r.HandleFunc("/authenticate", Authenticate)

	go AcceptInput()
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}

}

func ClientsInit() {
	clientList, err := database.GetClients()
	if err != nil {
		log.Fatal("Failed to initialize clients:", err)
	} else {
		for _, client := range clientList {
			var fs file_system.FileSystem
			fs.SetRootDir(client.RootDir)
			newClient := Client{
				inf: helper.ClientInf{Name: client.Name},
				fs:  fs,
			}
			clients = append(clients, newClient)
		}
	}
}

func AcceptInput() {
	reader := bufio.NewReader(os.Stdin)
	for {
		var req string
		req, _ = reader.ReadString('\n')
		req = req[0 : len(req)-2]
		arguments := strings.Split(req, " ")
		switch arguments[0] {
		case "newclient":
			{
				if len(arguments) < 3 {
					fmt.Println("Not enough arguments")
					continue
				} else if len(arguments) > 3 {
					fmt.Println("Too much arguments")
					continue
				}
				name := arguments[1]
				password := arguments[2]
				clientInf := helper.ClientInf{
					Name:    name,
					RootDir: storagePath + name + pathSeparator,
				}
				err := database.CreateClient(name, password, clientInf.RootDir)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("Client created successfully")
				}
				rootFs.MkDir(pathSeparator, clientInf.RootDir)
				var fs file_system.FileSystem
				fs.SetRootDir(clientInf.RootDir)
				newClient := Client{
					inf: helper.ClientInf{Name: clientInf.Name},
					fs:  fs,
				}
				clients = append(clients, newClient)

			}
		case "getclients":
			{
				if len(arguments) < 1 {
					fmt.Println("Not enough arguments")
					continue
				} else if len(arguments) > 1 {
					fmt.Println("Too much arguments")
					continue
				}
				clientNames, err := database.GetClients()
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("Clients:")
					for _, c := range clientNames {
						fmt.Println(c.Name)
					}
				}

			}
		default:
			{

			}

		}
	}
}

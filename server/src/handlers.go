package main

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"io"
	"net/http"
	"os"
	"remote-storage/schemas"
	"remote-storage/server/src/helper"
	"strconv"
	"time"
)

const chunkSize = 64 * 1024

func f(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(300)
	w.Write([]byte("Hello, world!"))
}

func Cd(w http.ResponseWriter, r *http.Request) {
	client, err := checkBearer(r)
	if err != nil {
		w.WriteHeader(401)
		return
	}
	var request schemas.CdRequest
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	curPath, err := client.fs.Cd(request.Path)

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
	client, err := checkBearer(r)
	if err != nil {
		w.WriteHeader(401)
		return
	}
	var request schemas.LsRequest
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	output, err := client.fs.Ls(request.DirPath)

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
	client, err := checkBearer(r)
	if err != nil {
		w.WriteHeader(401)
		return
	}
	var request schemas.MkDirRequest
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	curPath, err := client.fs.MkDir(request.Path, request.DirName)

	w.WriteHeader(200)
	var response schemas.MkDirResponse
	response.Path = curPath
	if err != nil {
		response.Message = err.Error()
	} else {
		response.Message = "directory created successfully"
	}
	body, err := json.Marshal(response)
	w.Write([]byte(body))
}

func Rename(w http.ResponseWriter, r *http.Request) {
	client, err := checkBearer(r)
	if err != nil {
		w.WriteHeader(401)
		return
	}
	var request schemas.RenameRequest
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	path, err := client.fs.RenameCmd(request.DirPath, request.OldName, request.NewName)

	w.WriteHeader(200)
	var response schemas.RenameResponse
	response.Path = path
	if err != nil {
		response.Message = err.Error()
	} else {
		response.Message = "renamed successfully"
	}
	body, err := json.Marshal(response)
	w.Write([]byte(body))
}

func Move(w http.ResponseWriter, r *http.Request) {
	client, err := checkBearer(r)
	if err != nil {
		w.WriteHeader(401)
		return
	}
	var request schemas.MoveRequest
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	path, err := client.fs.MoveCmd(request.SrcDirPath, request.FileName, request.DestDirPath)

	w.WriteHeader(200)
	var response schemas.CopyResponse
	response.Path = path
	if err != nil {
		response.Message = err.Error()
	} else {
		response.Message = "moved successfully"
	}
	body, err := json.Marshal(response)
	w.Write([]byte(body))
}

func Copy(w http.ResponseWriter, r *http.Request) {
	client, err := checkBearer(r)
	if err != nil {
		w.WriteHeader(401)
		return
	}
	var request schemas.CopyRequest
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	path, err := client.fs.CopyCmd(request.SrcDirPath, request.FileName, request.DestDirPath)

	w.WriteHeader(200)
	var response schemas.CopyResponse
	response.Path = path
	if err != nil {
		response.Message = err.Error()
	} else {
		response.Message = "copied successfully"
	}
	body, err := json.Marshal(response)
	w.Write([]byte(body))
}

func Delete(w http.ResponseWriter, r *http.Request) {
	client, err := checkBearer(r)
	if err != nil {
		w.WriteHeader(401)
		return
	}
	var request schemas.DeleteRequest
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	path, err := client.fs.DeleteCmd(request.DirPath, request.FileName)

	w.WriteHeader(200)
	var response schemas.CopyResponse
	response.Path = path
	if err != nil {
		response.Message = err.Error()
	} else {
		response.Message = "deleted successfully"
	}
	body, err := json.Marshal(response)
	w.Write([]byte(body))
}

func StartUploading(w http.ResponseWriter, r *http.Request) {
	client, err := checkBearer(r)
	if err != nil {
		w.WriteHeader(401)
		return
	}
	var request schemas.StartUploadRequest
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	tempDir := request.Location + request.FileName + "___temp___"
	client.fs.Mkdir(tempDir, 0644)

	client.upload.chunksNum = request.ChunksNum
	client.upload.tempDir = tempDir
	client.upload.fileName = request.FileName
	client.upload.fileLocation = request.Location
	client.upload.chunksGotten = 0
	client.upload.receivedChunks = make([]bool, client.upload.chunksNum)

	w.WriteHeader(200)
	var response schemas.StartUploadResponse
	if err != nil {
		response.Message = err.Error()
	} else {
		response.Message = "upload started successfully"
	}
	body, err := json.Marshal(response)
	w.Write([]byte(body))

}

func UploadChunk(w http.ResponseWriter, r *http.Request) {
	client, err := checkBearer(r)
	if err != nil {
		w.WriteHeader(401)
		return
	}
	var id string
	var data []byte
	id = r.URL.Query().Get("id")
	num, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}
	data, _ = io.ReadAll(r.Body)
	filePath := client.upload.tempDir + string(os.PathSeparator) + id + ".bin"
	file, err := client.fs.Create(filePath)
	file.Close()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(200)
		client.upload.receivedChunks[num] = true
	}
	client.fs.Write(filePath, data)

}

func UploadComplete(w http.ResponseWriter, r *http.Request) {
	client, err := checkBearer(r)
	if err != nil {
		w.WriteHeader(401)
		return
	}
	var response schemas.CompleteUploadResponse

	var missedChunks = make([]int, 0)
	for i, _ := range client.upload.receivedChunks {
		if !client.upload.receivedChunks[i] {
			missedChunks = append(missedChunks, i)
		}
	}
	if len(missedChunks) == 0 {
		err = client.fs.AssembleFiles(client.upload.fileLocation, client.upload.tempDir, client.upload.fileName)
		if err != nil {
			w.WriteHeader(500)
		}
		client.fs.RemoveAll(client.upload.tempDir)
		response.Message = "upload completed successfully"
		body, _ := json.Marshal(response)
		w.Write([]byte(body))
	} else {
		response.MissedChunks = missedChunks
		body, _ := json.Marshal(response)
		w.Write([]byte(body))
	}
}

func StartDownloading(w http.ResponseWriter, r *http.Request) {
	client, err := checkBearer(r)
	if err != nil {
		w.WriteHeader(401)
		return
	}
	var request schemas.StartDownloadRequest
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	client.download.chunks, err = client.fs.DivideFileIntoChunks(request.Location+request.FileName, chunkSize)

	w.WriteHeader(200)
	var response schemas.StartDownloadResponse
	response.ChunksNum = len(client.download.chunks)
	if err != nil {
		response.Message = err.Error()
	} else {
		response.Message = "download started successfully"
	}
	body, err := json.Marshal(response)
	w.Write([]byte(body))

}

func DownloadChunk(w http.ResponseWriter, r *http.Request) {
	client, err := checkBearer(r)
	if err != nil {
		w.WriteHeader(401)
		return
	}
	var id string
	id = r.URL.Query().Get("id")
	num, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}
	w.WriteHeader(200)
	w.Write(client.download.chunks[num])
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	var err error
	var request schemas.AuthenticateRequest
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	hashedPassword, _ := database.GetHashedPassword(request.Name)
	if helper.Hash(request.Name+request.Password) != hashedPassword {
		w.WriteHeader(401)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: request.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Finally, we set the Client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	var response schemas.AuthenticateResponse
	response.RootDir = request.Name + pathSeparator
	body, err := json.Marshal(response)
	w.Write([]byte(body))
	w.WriteHeader(200)
}

func checkBearer(r *http.Request) (*Client, error) {
	var res *Client
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("token")
	if err != nil {
		return nil, err
	}

	// Get the JWT string from the cookie
	tknStr := c.Value

	// Initialize a new instance of `Claims`
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !tkn.Valid {
		return nil, err
	}
	username := claims.Username
	for i, _ := range clients {
		if clients[i].inf.Name == username {
			res = &clients[i]
		}
	}
	return res, nil
}

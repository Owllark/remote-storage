package storagesvc

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"io"
	"net/http"
	"net/url"
	"remote-storage/server/authsvc"
	fs "remote-storage/server/storagesvc/repository/filesystem"
	"strings"
)

type Endpoints struct {
	GetStateEndpoint endpoint.Endpoint
	MkDirEndpoint    endpoint.Endpoint
	RenameEndpoint   endpoint.Endpoint
	MoveEndpoint     endpoint.Endpoint
	CopyEndpoint     endpoint.Endpoint
	DeleteEndpoint   endpoint.Endpoint
	DownloadEndpoint endpoint.Endpoint
	UploadEndpoint   endpoint.Endpoint
	cookies          map[string]*http.Cookie
}

type Cookies []*http.Cookie
type cookiesContextKey string

const contextKeyRequestCookie = cookiesContextKey("request_cookies")

// MakeServerEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the Service.
func MakeServerEndpoints(s Service) Endpoints {
	endpoints := Endpoints{
		GetStateEndpoint: MakeGetStateEndpoint(s),
		MkDirEndpoint:    MakeMkDirEndpoint(s),
		RenameEndpoint:   MakeRenameEndpoint(s),
		MoveEndpoint:     MakeMoveEndpoint(s),
		CopyEndpoint:     MakeCopyEndpoint(s),
		DeleteEndpoint:   MakeDeleteEndpoint(s),
		DownloadEndpoint: MakeDownloadEndpoint(s),
		UploadEndpoint:   MakeUploadEndpoint(s),
	}

	return endpoints
}

// MakeClientEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the remote instance, via a transport/http.Client.
// Useful in the storagesvc client.
func MakeClientEndpoints(instance string) (Endpoints, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	tgt, err := url.Parse(instance)
	if err != nil {
		return Endpoints{}, err
	}
	tgt.Path = ""

	options := []httptransport.ClientOption{}

	return Endpoints{
		GetStateEndpoint: httptransport.NewClient("POST", tgt, encodeGetStateRequest, decodeGetStateResponse, options...).Endpoint(),
		MkDirEndpoint:    httptransport.NewClient("POST", tgt, encodeMkDirRequest, decodeMkDirResponse, options...).Endpoint(),
		RenameEndpoint:   httptransport.NewClient("POST", tgt, encodeRenameRequest, decodeRenameResponse, options...).Endpoint(),
		MoveEndpoint:     httptransport.NewClient("POST", tgt, encodeMoveRequest, decodeMoveResponse, options...).Endpoint(),
		CopyEndpoint:     httptransport.NewClient("POST", tgt, encodeCopyRequest, decodeCopyResponse, options...).Endpoint(),
		DeleteEndpoint:   httptransport.NewClient("POST", tgt, encodeDeleteRequest, decodeDeleteResponse, options...).Endpoint(),
		DownloadEndpoint: httptransport.NewClient("POST", tgt, encodeDownloadRequest, decodeDownloadResponse, options...).Endpoint(),
		UploadEndpoint:   httptransport.NewClient("POST", tgt, encodeUploadRequest, decodeUdploadResponse, options...).Endpoint(),
	}, nil
}

// GetState implements Service. Primarily useful in a client.
func (e Endpoints) GetState(ctx context.Context) (fs.FileInfo, error) {
	e.setCookies(ctx)
	request := getStateRequest{
		UserRootDir: userRootDir,
	}
	response, err := e.GetStateEndpoint(ctx, request)
	if err != nil {
		return fs.FileInfo{}, err
	}
	resp := response.(getStateResponse)
	return resp.Info, errors.New(resp.Error)
}

// MkDir implements Service. Primarily useful in a client.
func (e Endpoints) MkDir(ctx context.Context, path string, dirName string) (string, error) {
	e.setCookies(ctx)
	request := mkDirRequest{
		Path:    path,
		DirName: dirName,
	}
	response, err := e.MkDirEndpoint(ctx, request)
	if err != nil {
		return "", err
	}
	resp := response.(mkDirResponse)
	return resp.Path, errors.New(resp.Error)
}

// Rename implements Service. Primarily useful in a client.
func (e Endpoints) Rename(ctx context.Context, dirPath string, oldName string, newName string) (string, error) {
	e.setCookies(ctx)
	request := renameRequest{
		DirPath: dirPath,
		OldName: oldName,
		NewName: newName,
	}
	response, err := e.RenameEndpoint(ctx, request)
	if err != nil {
		return "", err
	}
	resp := response.(renameResponse)
	return resp.Path, errors.New(resp.Error)
}

// Move implements Service. Primarily useful in a client.
func (e Endpoints) Move(ctx context.Context, srcDirPath string, fileName string, destDirPath string) (string, error) {
	e.setCookies(ctx)
	request := moveRequest{
		SrcDirPath:  srcDirPath,
		FileName:    fileName,
		DestDirPath: destDirPath,
	}
	response, err := e.MoveEndpoint(ctx, request)
	if err != nil {
		return "", err
	}
	resp := response.(moveResponse)
	return resp.Path, errors.New(resp.Error)
}

// Copy implements Service. Primarily useful in a client.
func (e Endpoints) Copy(ctx context.Context, srcDirPath string, fileName string, destDirPath string) (string, error) {
	e.setCookies(ctx)
	request := copyRequest{
		SrcDirPath:  srcDirPath,
		FileName:    fileName,
		DestDirPath: destDirPath,
	}
	response, err := e.CopyEndpoint(ctx, request)
	if err != nil {
		return "", err
	}
	resp := response.(copyResponse)
	return resp.Path, errors.New(resp.Error)
}

// Delete implements Service. Primarily useful in a client.
func (e Endpoints) Delete(ctx context.Context, dirPath string, fileName string) (string, error) {
	e.setCookies(ctx)
	request := deleteRequest{
		DirPath:  dirPath,
		FileName: fileName,
	}
	response, err := e.DeleteEndpoint(ctx, request)
	if err != nil {
		return "", err
	}
	resp := response.(deleteResponse)
	return resp.Path, errors.New(resp.Error)
}

// Download implements Service. Primarily useful in a client.
func (e Endpoints) Download(ctx context.Context, dirPath string, fileName string) (io.ReadCloser, error) {
	e.setCookies(ctx)
	request := downloadRequest{
		DirPath:  dirPath,
		FileName: fileName,
	}
	response, err := e.DownloadEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	resp := response.(downloadResponse)
	return resp.Buffer, errors.New(resp.Error)
}

// Upload implements Service. Primarily useful in a client.
func (e Endpoints) Upload(ctx context.Context, dirPath string, fileName string, contents io.ReadCloser) error {
	e.setCookies(ctx)
	request := uploadRequest{
		DirPath:  dirPath,
		FileName: fileName,
		Contents: contents,
	}
	response, err := e.UploadEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(uploadResponse)
	return errors.New(resp.Error)
}

func (e Endpoints) getAuthSvc() authsvc.Service {
	return nil
}

func (e Endpoints) AddCookie(cookie *http.Cookie) {
	if e.cookies == nil {
		e.cookies = make(map[string]*http.Cookie)
	}
	e.cookies[cookie.Name] = cookie
}

func (e Endpoints) DeleteCookie(cookie *http.Cookie) {
	if e.cookies == nil {
		e.cookies = make(map[string]*http.Cookie)
	}
	delete(e.cookies, cookie.Value)
}

func (e Endpoints) setCookies(ctx context.Context) context.Context {

	var cookieArr Cookies
	for _, v := range e.cookies {
		cookieArr = append(cookieArr, v)
	}
	return context.WithValue(ctx, contextKeyRequestCookie, cookieArr)
}

func MakeGetStateEndpoint(svc Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(getStateRequest)
		resp, err := svc.GetState(nil)
		if err != nil {
			return getStateResponse{resp, err.Error()}, nil
		}
		return getStateResponse{resp, ""}, nil
	}
}

func MakeMkDirEndpoint(svc Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(mkDirRequest)
		resp, err := svc.MkDir(nil, req.Path, req.DirName)
		if err != nil {
			return mkDirResponse{resp, err.Error()}, nil
		}
		return mkDirResponse{resp, ""}, nil
	}
}

func MakeRenameEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(renameRequest)
		resp, err := svc.Rename(ctx, req.DirPath, req.OldName, req.NewName)
		if err != nil {
			return renameResponse{resp, err.Error()}, nil
		}
		return renameResponse{resp, ""}, nil
	}
}

func MakeMoveEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(moveRequest)
		resp, err := svc.Move(ctx, req.SrcDirPath, req.FileName, req.DestDirPath)
		if err != nil {
			return moveResponse{resp, err.Error()}, nil
		}
		return moveResponse{resp, ""}, nil
	}
}

func MakeCopyEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(copyRequest)
		resp, err := svc.Copy(ctx, req.SrcDirPath, req.FileName, req.DestDirPath)
		if err != nil {
			return copyResponse{resp, err.Error()}, nil
		}
		return copyResponse{resp, ""}, nil
	}
}

func MakeDeleteEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteRequest)
		resp, err := svc.Delete(ctx, req.DirPath, req.FileName)
		if err != nil {
			return deleteResponse{resp, err.Error()}, nil
		}
		return deleteResponse{resp, ""}, nil
	}
}

func MakeDownloadEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(downloadRequest)
		resp, err := svc.Download(ctx, req.DirPath, req.FileName)
		if err != nil {
			return downloadResponse{resp, err.Error()}, nil
		}
		return downloadResponse{resp, ""}, nil
	}
}

func MakeUploadEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(uploadRequest)
		err := svc.Upload(ctx, req.DirPath, req.FileName, req.Contents)
		if err != nil {
			return uploadResponse{err.Error()}, nil
		}
		return uploadResponse{""}, nil
	}
}

type getStateRequest struct {
	UserRootDir string `json:"user_root_dir,omitempty"`
}

type getStateResponse struct {
	Info  fs.FileInfo `json:"info"`
	Error string      `json:"error,omitempty"`
}

func (r getStateResponse) error() error {
	return errors.New(r.Error)
}

type mkDirRequest struct {
	// path to directory where new directory must be created
	Path string `json:"path"`
	// name of the directory to be created
	DirName string `json:"dir_name,omitempty"`
}

type mkDirResponse struct {
	// path to created directory
	Path string `json:"path,omitempty"`
	// message if something went wrong
	Error string `json:"error,omitempty"`
}

func (r mkDirResponse) error() error {
	return errors.New(r.Error)
}

type renameRequest struct {
	// path to directory
	DirPath string `json:"dir_path"`
	// old name of file
	OldName string `json:"old_name,omitempty"`
	// new name of file
	NewName string `json:"new_name,omitempty"`
}

type renameResponse struct {
	// path to the renamed file
	Path string `json:"path,omitempty"`
	// message if something went wrong
	Error string `json:"error,omitempty"`
}

func (r renameResponse) error() error {
	return errors.New(r.Error)
}

type moveRequest struct {
	// path to source directory
	SrcDirPath string `json:"src_dir_path"`
	// name of file to be moved
	FileName string `json:"file_name,omitempty"`
	// path to destination directory
	DestDirPath string `json:"dest_dir_path"`
}

type moveResponse struct {
	// path to the moved file
	Path string `json:"path,omitempty"`
	// message if something went wrong
	Error string `json:"error,omitempty"`
}

func (r moveResponse) error() error {
	return errors.New(r.Error)
}

type copyRequest struct {
	// path to source directory
	SrcDirPath string `json:"src_dir_path"`
	// name of file to be copied
	FileName string `json:"file_name,omitempty"`
	// path to destination directory
	DestDirPath string `json:"dest_dir_path"`
}

type copyResponse struct {
	// path to the copied file
	Path string `json:"path,omitempty"`
	// message if something went wrong
	Error string `json:"error,omitempty"`
}

func (r copyResponse) error() error {
	return errors.New(r.Error)
}

type deleteRequest struct {
	// path to the directory
	DirPath string `json:"dir_path"`
	// name of file to be deleted
	FileName string `json:"file_name,omitempty"`
}

type deleteResponse struct {
	// path to the deleted file
	Path string `json:"path,omitempty"`
	// message if something went wrong
	Error string `json:"error,omitempty"`
}

func (r deleteResponse) error() error {
	return errors.New(r.Error)
}

type downloadRequest struct {
	// path to the directory
	DirPath string `json:"dir_path"`
	// name of file to be deleted
	FileName string `json:"file_name,omitempty"`
}

type downloadResponse struct {
	// path to the deleted file
	Buffer io.ReadCloser `json:"buffer,omitempty"`
	// message if something went wrong
	Error string `json:"error,omitempty"`
}

func (r downloadResponse) error() error {
	return errors.New(r.Error)
}

func (r downloadResponse) ReadCloser() io.ReadCloser {
	return r.Buffer
}

type uploadRequest struct {
	// path to the directory
	DirPath string `json:"dir_path"`
	// name of file to be deleted
	FileName string `json:"file_name,omitempty"`
	//
	Contents io.ReadCloser `json:"contents"`
}

type uploadResponse struct {

	// message if something went wrong
	Error string `json:"error,omitempty"`
}

func (r uploadResponse) error() error {
	return errors.New(r.Error)
}

func (r uploadRequest) ReadCloser() io.ReadCloser {
	return r.Contents
}

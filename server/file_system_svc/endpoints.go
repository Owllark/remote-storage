package file_system_svc

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	fs "server/file_system_svc/repository/filesystem"
)

type Endpoints struct {
	GetStateEndpoint endpoint.Endpoint
	MkDirEndpoint    endpoint.Endpoint
	RenameEndpoint   endpoint.Endpoint
	MoveEndpoint     endpoint.Endpoint
	CopyEndpoint     endpoint.Endpoint
	DeleteEndpoint   endpoint.Endpoint
}

// MakeServerEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the FileSystemService.
func MakeServerEndpoints(s FileSystemService) Endpoints {
	endpoints := Endpoints{
		GetStateEndpoint: makeGetStateEndpoint(s),
		MkDirEndpoint:    makeMkDirEndpoint(s),
		RenameEndpoint:   makeRenameEndpoint(s),
		MoveEndpoint:     makeMoveEndpoint(s),
		CopyEndpoint:     makeCopyEndpoint(s),
		DeleteEndpoint:   makeDeleteEndpoint(s),
	}

	return endpoints
}

func makeGetStateEndpoint(svc FileSystemService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(getStateRequest)
		resp, err := svc.GetState(req.UserRootDir)
		if err != nil {
			return getStateResponse{resp, err.Error()}, nil
		}
		return getStateResponse{resp, ""}, nil
	}
}

func makeMkDirEndpoint(svc FileSystemService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(mkDirRequest)
		resp, err := svc.MkDir(req.DirName, req.Path)
		if err != nil {
			return mkDirResponse{resp, err.Error()}, nil
		}
		return mkDirResponse{resp, ""}, nil
	}
}

func makeRenameEndpoint(svc FileSystemService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(renameRequest)
		resp, err := svc.Rename(req.DirPath, req.OldName, req.NewName)
		if err != nil {
			return renameResponse{resp, err.Error()}, nil
		}
		return renameResponse{resp, ""}, nil
	}
}

func makeMoveEndpoint(svc FileSystemService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(moveRequest)
		resp, err := svc.Move(req.SrcDirPath, req.FileName, req.DestDirPath)
		if err != nil {
			return moveResponse{resp, err.Error()}, nil
		}
		return moveResponse{resp, ""}, nil
	}
}

func makeCopyEndpoint(svc FileSystemService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(copyRequest)
		resp, err := svc.Copy(req.SrcDirPath, req.FileName, req.DestDirPath)
		if err != nil {
			return copyResponse{resp, err.Error()}, nil
		}
		return copyResponse{resp, ""}, nil
	}
}

func makeDeleteEndpoint(svc FileSystemService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteRequest)
		resp, err := svc.Delete(req.DirPath, req.FileName)
		if err != nil {
			return deleteResponse{resp, err.Error()}, nil
		}
		return deleteResponse{resp, ""}, nil
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

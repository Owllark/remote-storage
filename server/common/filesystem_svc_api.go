package common

import (
	fs "server/file_system_svc/repository/filesystem"
)

type GetStateRequest struct {
	UserRootDir string `json:"user_root_dir,omitempty"`
}

type GetStateResponse struct {
	Info  fs.FileInfo `json:"info"`
	Error string      `json:"error,omitempty"`
}

type MkDirRequest struct {
	// path to directory where new directory must be created
	Path string `json:"path"`
	// name of the directory to be created
	DirName string `json:"dir_name,omitempty"`
}

type MkDirResponse struct {
	// path to created directory
	Path string `json:"path,omitempty"`
	// message if something went wrong
	Error string `json:"error,omitempty"`
}

type RenameRequest struct {
	// path to directory
	DirPath string `json:"dir_path"`
	// old name of file
	OldName string `json:"old_name,omitempty"`
	// new name of file
	NewName string `json:"new_name,omitempty"`
}

type RenameResponse struct {
	// path to the renamed file
	Path string `json:"path,omitempty"`
	// message if something went wrong
	Error string `json:"error,omitempty"`
}

type MoveRequest struct {
	// path to source directory
	SrcDirPath string `json:"src_dir_path"`
	// name of file to be moved
	FileName string `json:"file_name,omitempty"`
	// path to destination directory
	DestDirPath string `json:"dest_dir_path"`
}

type MoveResponse struct {
	// path to the moved file
	Path string `json:"path,omitempty"`
	// message if something went wrong
	Error string `json:"error,omitempty"`
}

type CopyRequest struct {
	// path to source directory
	SrcDirPath string `json:"src_dir_path"`
	// name of file to be copied
	FileName string `json:"file_name,omitempty"`
	// path to destination directory
	DestDirPath string `json:"dest_dir_path"`
}

type CopyResponse struct {
	// path to the copied file
	Path string `json:"path,omitempty"`
	// message if something went wrong
	Error string `json:"error,omitempty"`
}

type DeleteRequest struct {
	// path to the directory
	DirPath string `json:"dir_path"`
	// name of file to be deleted
	FileName string `json:"file_name,omitempty"`
}

type DeleteResponse struct {
	// path to the deleted file
	Path string `json:"path,omitempty"`
	// message if something went wrong
	Error string `json:"error,omitempty"`
}

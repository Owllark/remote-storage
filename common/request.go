package common

// all paths must be with separator in the end

type CdRequest struct {
	// path to directory
	Path string `json:"path,omitempty"`
}

type MkDirRequest struct {
	// path to directory where new directory must be created
	Path string `json:"path"`
	// name of the directory to be created
	DirName string `json:"dir_name,omitempty"`
}

type RenameRequest struct {
	// path to directory
	DirPath string `json:"dir_path"`
	// old name of file
	OldName string `json:"old_name,omitempty"`
	// new name of file
	NewName string `json:"new_name,omitempty"`
}

type MoveRequest struct {
	// path to source directory
	SrcDirPath string `json:"src_dir_path"`
	// name of file to be moved
	FileName string `json:"file_name,omitempty"`
	// path to destination directory
	DestDirPath string `json:"dest_dir_path"`
}

type CopyRequest struct {
	// path to source directory
	SrcDirPath string `json:"src_dir_path"`
	// name of file to be copied
	FileName string `json:"file_name,omitempty"`
	// path to destination directory
	DestDirPath string `json:"dest_dir_path"`
}

type DeleteRequest struct {
	// path to the directory
	DirPath string `json:"dir_path"`
	// name of file to be deleted
	FileName string `json:"file_name,omitempty"`
}

type LsRequest struct {
	DirPath string `json:"dir_path,omitempty"`
}

type TreeRequest struct {
	DirPath string `json:"dir_path,omitempty"`
}

type StartUploadRequest struct {
	Location  string `json:"location,omitempty"`
	FileName  string `json:"file_name,omitempty"`
	ChunksNum int    `json:"chunks_num"`
}

type UploadChunkRequest struct {
	Id   int    `json:"id,omitempty"`
	Data []byte `json:"data,omitempty"`
}

type DownloadChunkRequest struct {
}

type CompleteUploadRequest struct {
}

type StartDownloadRequest struct {
	Location string `json:"location,omitempty"`
	FileName string `json:"file_name,omitempty"`
}

type AuthenticateRequest struct {
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
}

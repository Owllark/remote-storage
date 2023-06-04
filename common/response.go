package common

// all paths must be with separator in the end

type CdResponse struct {
	// path to directory
	Path string `json:"path,omitempty"`
	// message if something went wrong
	Message string `json:"message,omitempty"`
}

type MkDirResponse struct {
	// path to created directory
	Path string `json:"path,omitempty"`
	// message if something went wrong
	Message string `json:"message,omitempty"`
}

type RenameResponse struct {
	// path to the renamed file
	Path string `json:"path,omitempty"`
	// message if something went wrong
	Message string `json:"message,omitempty"`
}

type MoveResponse struct {
	// path to the moved file
	Path string `json:"path,omitempty"`
	// message if something went wrong
	Message string `json:"message,omitempty"`
}

type CopyResponse struct {
	// path to the copied file
	Path string `json:"path,omitempty"`
	// message if something went wrong
	Message string `json:"message,omitempty"`
}

type DeleteResponse struct {
	// path to the deleted file
	Path string `json:"path,omitempty"`
	// message if something went wrong
	Message string `json:"message,omitempty"`
}

type LsResponse struct {
	// command output
	CommandOutput []string `json:"command_output,omitempty"`
	// message if something went wrong
	Message string `json:"message,omitempty"`
}

type TreeResponse struct {
	// path to the directory
	DirPath string `json:"dir_path,omitempty"`
	// command output
	CommandOutput []string `json:"command_output,omitempty"`
}

type StartUploadResponse struct {
	// message if something went wrong
	Message string `json:"message,omitempty"`
}

type CompleteUploadResponse struct {
	// array of missed chunks ids
	MissedChunks []int `json:"missed_chunks"`
	// message if something went wrong
	Message string `json:"message,omitempty"`
}

type UploadChunkResponse struct {
}

type DownloadChunkResponse struct {
	Data []byte `json:"data,omitempty"`
}

type StartDownloadResponse struct {
	ChunksNum int `json:"chunks_num"`
	// message if something went wrong
	Message string `json:"message,omitempty"`
}

type AuthenticateResponse struct {
	RootDir string `json:"root_dir,omitempty"`
}

type GetStateResponse struct {
	Info FileInfo `json:"info"`
}

package schemas

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

package schemas

// all paths must be without separator in the end

type CdRequest struct {
	// path to directory without separator at the end
	Path string `json:"path,omitempty"`
}

type MkDirRequest struct {
	// path to directory where new directory must be created, without separator at the end
	Path string `json:"path"`
	// name of the directory to be created
	DirName string `json:"dir_name,omitempty"`
}

type RenameRequest struct {
	// path to directory without separator at the end
	DirPath string `json:"dir_path"`
	// old name of file
	OldName string `json:"old_name,omitempty"`
	// new name of file
	NewName string `json:"new_name,omitempty"`
}

type MoveRequest struct {
	// path to source directory without separator at the end
	SrcDirPath string `json:"src_dir_path"`
	// name of file to be moved
	FileName string `json:"file_name,omitempty"`
	// path to destination directory without separator at the end
	DestDirPath string `json:"dest_dir_path"`
}

type CopyRequest struct {
	// path to source directory without separator at the end
	SrcDirPath string `json:"src_dir_path"`
	// name of file to be copied
	FileName string `json:"file_name,omitempty"`
	// path to destination directory without separator at the end
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

package common

import "time"

type FileInfo struct {
	Name     string     `json:"name"`
	IsDir    bool       `json:"is_dir"`
	Size     int64      `json:"size"`
	Modified time.Time  `json:"modified"`
	Children []FileInfo `json:"children,omitempty"` // Nested files and directories
}

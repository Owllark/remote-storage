package file_system

import (
	"os"
	"strings"
)

const pathSeparator = string(os.PathSeparator)

func (fs *FileSystem) getPath(path string) string {
	var res string
	res += "." + pathSeparator + "storage" + pathSeparator
	res += fs.currentDir
	res += path
	return res
}

func joinPath(path []string) string {
	var res string
	for _, s := range path {
		res += s + string(os.PathSeparator)
	}
	return res
}

func splitPath(path string) []string {
	res := strings.Split(path, string(os.PathSeparator))
	return res
}

func (fs *FileSystem) isExists(path string) bool {
	path = fs.getPath(path)
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (fs *FileSystem) isDirectoryExists(path string) bool {
	path = fs.getPath(path)
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false // Directory does not exist
	}
	return info.IsDir()
}

func (fs *FileSystem) isFileExists(path string) bool {
	path = fs.getPath(path)
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false // File does not exist
	}
	return !info.IsDir()
}

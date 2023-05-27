package file_system

import (
	"errors"
	"os"
	"strings"
)

type FileSystem struct {
	currentDir  string
	currentPath []string
	prevPaths   []string
}

func (fs *FileSystem) CreateFile(dirPath, filename string) (*os.File, error) {
	file, err := os.Create(getPath(fs.currentDir, dirPath, filename))
	return file, err
}

func (fs *FileSystem) WriteFile(filepath string, data []byte) (int, error) {
	file, err := os.Open(getPath(fs.currentDir, filepath))
	if err != nil {
		return 0, err
	}
	writtenBytesNum, err := file.Write(data)
	return writtenBytesNum, err
}

func (fs *FileSystem) DeleteFile(filePath string) error {
	err := os.Remove(getPath(fs.currentDir, filePath))
	return err
}

func (fs *FileSystem) RenameFile(dirPath, oldName, newName string) error {
	err := os.Rename(getPath(fs.currentDir, dirPath, oldName), getPath(fs.currentDir, dirPath, newName))
	return err
}

func (fs *FileSystem) MoveFile(oldPath, newPath string) error {
	err := os.Rename(getPath(fs.currentDir, oldPath), getPath(fs.currentDir, newPath))
	return err
}

func (fs *FileSystem) Cd(path string) error {
	switch path {
	case "..":
		{
			if len(fs.currentPath) == 0 {
				return errors.New("cd cannot move one level above, already in root directory")
			}
			fs.prevPaths = append(fs.prevPaths, joinPath(fs.currentPath))
			fs.currentDir = fs.currentPath[len(fs.currentPath)-1]
			fs.currentPath = fs.currentPath[:len(fs.currentPath)-1]
		}
	case "/":
		{
			fs.prevPaths = append(fs.prevPaths, joinPath(fs.currentPath))
			fs.currentDir = ""
			fs.currentPath = fs.currentPath[:0]
		}
	case "":
		{
			fs.prevPaths = append(fs.prevPaths, joinPath(fs.currentPath))
			fs.currentDir = ""
			fs.currentPath = fs.currentPath[:0]
		}
	case "-":
		{
			if len(fs.prevPaths) != 0 {
				fs.currentPath = splitPath(fs.prevPaths[len(fs.prevPaths)-1])
				fs.currentDir = fs.currentPath[len(fs.currentPath)-1]
				fs.prevPaths = fs.prevPaths[:len(fs.prevPaths)]
			}
		}
	default:
		{
			if IsDirectoryExists(path) {
				fs.prevPaths = append(fs.prevPaths, joinPath(fs.currentPath))
				fs.currentDir = fs.currentPath[len(fs.currentPath)-1]
				fs.currentPath = fs.currentPath[:len(fs.currentPath)-1]
			}

		}

	}
	return nil
}

func IsDirectoryExists(path string) bool {
	_, err := os.Stat(path)
	return os.IsExist(err)
}

func getPath(strings ...string) string {
	path := "storage/"
	for _, s := range strings {
		path += s
	}
	return path
}

func joinPath(path []string) string {
	var res string
	for _, s := range path {
		res += s
	}
	return res
}

func splitPath(path string) []string {
	res := strings.Split(path, string(os.PathSeparator))
	return res
}

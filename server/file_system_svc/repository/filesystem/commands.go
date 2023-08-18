package filesystem

import (
	"io"
	"os"
	"path/filepath"
	"time"
)

type ConfigFileSystem struct {
	RootDir string
}

var rootDir string
var isInitialized = false

type FileInfo struct {
	Name     string     `json:"name"`
	IsDir    bool       `json:"is_dir"`
	Size     int64      `json:"size"`
	Modified time.Time  `json:"modified"`
	Children []FileInfo `json:"children,omitempty"` // Nested files and directories
}

func init() {

}

func checkForInitializing() {
	if !isInitialized {
		panic("package not initialized, root directory must be set")
	}
}

func InitializeFileSystem(config ConfigFileSystem) {
	rootDir = config.RootDir
	isInitialized = true
}

func Create(filepath string) (*os.File, error) {
	checkForInitializing()
	file, err := os.Create(getPath(filepath))
	return file, err
}

func OpenFile(filepath string) (*os.File, error) {
	checkForInitializing()
	file, err := os.OpenFile(getPath(filepath), os.O_RDWR, 0644)
	return file, err
}

func Write(filepath string, data []byte) (int, error) {
	checkForInitializing()
	file, err := os.OpenFile(getPath(filepath), os.O_RDWR, 0644)
	defer file.Close()
	if err != nil {
		return 0, err
	}
	writtenBytesNum, err := file.Write(data)
	return writtenBytesNum, err
}

func ReadAll(filepath string) ([]byte, error) {
	checkForInitializing()
	var contents []byte
	sourceFile, err := os.Open(getPath(filepath))
	if err != nil {
		return contents, err
	}
	defer sourceFile.Close()
	contents, err = io.ReadAll(sourceFile)
	if err != nil {
		return contents, err
	}
	return contents, nil
}

func RemoveAll(filePath string) error {
	checkForInitializing()
	err := os.RemoveAll(getPath(filePath))
	return err
}

func Remove(filePath string) error {
	checkForInitializing()
	err := os.Remove(getPath(filePath))
	return err
}

func Rename(oldPath, newPath string) error {
	checkForInitializing()
	err := os.Rename(getPath(oldPath), getPath(newPath))
	return err
}

func Move(srcPath, destPath string) error {
	checkForInitializing()
	err := os.Rename(getPath(srcPath), getPath(destPath))
	return err
}

func Mkdir(path string, permission os.FileMode) error {
	checkForInitializing()
	err := os.Mkdir(getPath(path), permission)
	return err
}

func Copy(srcPath, destPath string) error {
	checkForInitializing()
	srcPath = getPath(srcPath)
	destPath = getPath(destPath)
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return err
	}

	if !srcInfo.IsDir() {
		return copyFile(srcPath, destPath)
	}

	err = os.MkdirAll(destPath, srcInfo.Mode())
	if err != nil {
		return err
	}

	return filepath.Walk(srcPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(srcPath, filePath)
		if err != nil {
			return err
		}

		dest := filepath.Join(destPath, relPath)

		if info.IsDir() {
			err := os.MkdirAll(dest, info.Mode())
			if err != nil {
				return err
			}
		} else {
			err := copyFile(filePath, dest)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func copyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func TraverseDirectory(dirPath string) (FileInfo, error) {
	checkForInitializing()
	file := FileInfo{
		Name:     filepath.Base(dirPath),
		IsDir:    true,
		Children: []FileInfo{},
	}

	// Open the directory and get its contents
	dir, err := os.Open(dirPath)
	if err != nil {
		return file, err
	}
	defer dir.Close()

	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return file, err
	}

	// Iterate over the contents
	for _, fileInfo := range fileInfos {
		// Create a FileInfo struct for each file or directory
		child := FileInfo{
			Name:     fileInfo.Name(),
			IsDir:    fileInfo.IsDir(),
			Size:     fileInfo.Size(),
			Modified: fileInfo.ModTime(),
		}

		if child.IsDir {
			// Recursively traverse nested directories
			nestedFileInfo, err := TraverseDirectory(filepath.Join(dirPath, child.Name))
			if err != nil {
				return file, err
			}
			child.Children = nestedFileInfo.Children
		}

		file.Children = append(file.Children, child)
	}

	return file, nil
}

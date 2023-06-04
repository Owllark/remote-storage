package file_system

import (
	"fmt"
	"os"
	"path/filepath"
	"remote-storage/common"
	"strings"
)

const pathSeparator = string(os.PathSeparator)

func (fs *FileSystem) getPath(path string) string {
	var res string
	res += fs.rootDir
	//res += fs.currentDir
	res += path
	return res
}

func traverseDirectory(dirPath string) (common.FileInfo, error) {
	file := common.FileInfo{
		Name:     filepath.Base(dirPath),
		IsDir:    true,
		Children: []common.FileInfo{},
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
		child := common.FileInfo{
			Name:     fileInfo.Name(),
			IsDir:    fileInfo.IsDir(),
			Size:     fileInfo.Size(),
			Modified: fileInfo.ModTime(),
		}

		if child.IsDir {
			// Recursively traverse nested directories
			nestedFileInfo, err := traverseDirectory(filepath.Join(dirPath, child.Name))
			if err != nil {
				return file, err
			}
			child.Children = nestedFileInfo.Children
		}

		file.Children = append(file.Children, child)
	}

	return file, nil
}

func (fs *FileSystem) DivideFileIntoChunks(path string, chunkSize int) ([][]byte, error) {
	var chunks [][]byte
	file, err := os.Open(fs.getPath(path))
	if err != nil {
		fmt.Println("Error opening file:", err)
		return chunks, err
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()

	fullChunksAmount := int(fileSize) / chunkSize
	lastChunkSize := int(fileSize) % chunkSize
	chunkAmount := fullChunksAmount
	if lastChunkSize > 0 {
		chunkAmount++
	}
	chunks = make([][]byte, chunkAmount)

	for i := 0; i < fullChunksAmount; i++ {
		chunk := make([]byte, chunkSize)
		_, err := file.Read(chunk)
		if err != nil {
			return chunks, err
		}
		chunks[i] = chunk
	}

	if lastChunkSize > 0 {
		lastChunk := make([]byte, lastChunkSize)
		_, err := file.Read(lastChunk)
		if err != nil {
			return chunks, err
		}
		chunks[chunkAmount-1] = lastChunk
	}

	return chunks, err
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
	if res[len(res)-1] == "" {
		res = res[:len(res)-1]
	}

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

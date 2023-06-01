package file_system

import (
	"fmt"
	"os"
	"strings"
)

const pathSeparator = string(os.PathSeparator)

func (fs *FileSystem) getPath(path string) string {
	var res string
	res += fs.rootDir
	res += fs.currentDir
	res += path
	return res
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

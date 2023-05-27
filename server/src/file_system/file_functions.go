package file_system

import (
	"os"
)

func CreateFile(dirPath, filename string) (*os.File, error) {
	file, err := os.Create(getPath(dirPath, filename))
	return file, err
}

func WriteFile(filepath string, data []byte) (int, error) {
	file, err := os.Open(getPath(filepath))
	if err != nil {
		return 0, err
	}
	writtenBytesNum, err := file.Write(data)
	return writtenBytesNum, err
}

func DeleteFile(filePath string) error {
	err := os.Remove(getPath(filePath))
	return err
}

func RenameFile(dirPath, oldName, newName string) error {
	err := os.Rename(getPath(dirPath, oldName), getPath(dirPath, newName))
	return err
}

func MoveFile(oldPath, newPath string) error {
	err := os.Rename(getPath(oldPath), getPath(newPath))
	return err
}

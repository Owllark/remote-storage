package file_system

import (
	"io"
	"os"
	"path/filepath"
)

func (fs *FileSystem) Create(filepath string) (*os.File, error) {
	file, err := os.Create(fs.getPath(filepath))
	return file, err
}

func (fs *FileSystem) Write(filepath string, data []byte) (int, error) {
	file, err := os.OpenFile(fs.getPath(filepath), os.O_RDWR, 0644)
	defer file.Close()
	if err != nil {
		return 0, err
	}
	writtenBytesNum, err := file.Write(data)
	return writtenBytesNum, err
}

func (fs *FileSystem) RemoveAll(filePath string) error {
	err := os.RemoveAll(fs.getPath(filePath))
	return err
}

func (fs *FileSystem) Remove(filePath string) error {
	err := os.Remove(fs.getPath(filePath))
	return err
}

func (fs *FileSystem) Rename(oldPath, newPath string) error {
	err := os.Rename(fs.getPath(oldPath), fs.getPath(newPath))
	return err
}

func (fs *FileSystem) Move(srcPath, destPath string) error {
	err := os.Rename(fs.getPath(srcPath), fs.getPath(destPath))
	return err
}

func (fs *FileSystem) Mkdir(path string, permission os.FileMode) error {
	err := os.Mkdir(fs.getPath(path), permission)
	return err
}

func (fs *FileSystem) Copy(srcPath, destPath string) error {
	srcPath = fs.getPath(srcPath)
	destPath = fs.getPath(destPath)
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return err
	}

	err = os.MkdirAll(destPath, srcInfo.Mode())
	if err != nil {
		return err
	}

	if !srcInfo.IsDir() {
		return CopyFile(srcPath, destPath)
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
			err := CopyFile(filePath, dest)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func CopyFile(src, dest string) error {
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

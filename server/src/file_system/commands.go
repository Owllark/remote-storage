package file_system

import (
	"errors"
	"fmt"
	"os"
)

type FileSystem struct {
	currentDir  string
	currentPath []string
	prevPaths   []string
}

func (fs *FileSystem) Cd(path string) (string, error) {
	var err error
	switch path {
	case "..":
		{
			if len(fs.currentPath) == 0 {
				err = errors.New("cd cannot move one level above, already in root directory")
			}
			fs.prevPaths = append(fs.prevPaths, joinPath(fs.currentPath))
			fs.currentPath = fs.currentPath[:len(fs.currentPath)-1]
			fs.currentDir = joinPath(fs.currentPath)
		}
	case "/":
		{
			fs.prevPaths = append(fs.prevPaths, joinPath(fs.currentPath))
			fs.currentPath = fs.currentPath[:0]
			fs.currentDir = joinPath(fs.currentPath)
		}
	case "":
		{
			fs.prevPaths = append(fs.prevPaths, joinPath(fs.currentPath))
			fs.currentPath = fs.currentPath[:0]
			fs.currentDir = joinPath(fs.currentPath)
		}
	case "-":
		{
			if len(fs.prevPaths) != 0 {
				fs.currentPath = splitPath(fs.prevPaths[len(fs.prevPaths)-1])
				fs.prevPaths = fs.prevPaths[:len(fs.prevPaths)]
				fs.currentDir = joinPath(fs.currentPath)
			}
		}
	default:
		{
			if fs.isDirectoryExists(path) {
				fs.prevPaths = append(fs.prevPaths, joinPath(fs.currentPath))
				fs.currentPath = append(fs.currentPath, splitPath(path)...)
				fs.currentDir = joinPath(fs.currentPath)
			} else {
				err = errors.New(fmt.Sprintf("directory %s not found", path))
			}

		}

	}
	return joinPath(fs.currentPath), err
}

func (fs *FileSystem) MkDir(path string, dir string) (string, error) {
	var err error
	var newDirPath string
	newDirPath = path + dir

	if !fs.isDirectoryExists(path) {
		err = errors.New("cannot create directory, invalid path")
	} else {
		err = fs.Mkdir(newDirPath, 0644)
	}
	return newDirPath, err
}

func (fs *FileSystem) RenameCmd(dirPath, oldName, newName string) (string, error) {
	var err error
	var oldFilePath = dirPath + string(os.PathSeparator) + oldName
	var newFilePath = dirPath + string(os.PathSeparator) + newName

	if !fs.isExists(oldFilePath) {
		err = errors.New("cannot rename file, file not found")
	} else {
		err = fs.Rename(oldFilePath, newFilePath)
	}
	return newFilePath, err
}

func (fs *FileSystem) MoveCmd(srcDirPath, fileName, destDirPath string) (string, error) {
	var err error
	var oldFilePath = srcDirPath + string(os.PathSeparator) + fileName
	var newFilePath = destDirPath + string(os.PathSeparator) + fileName

	if !fs.isExists(oldFilePath) {
		err = errors.New("cannot move file, file not found")
	} else {
		err = fs.Move(oldFilePath, newFilePath)
	}
	return newFilePath, err
}

func (fs *FileSystem) DeleteCmd(dirPath, fileName string) (string, error) {
	var err error
	var filePath = dirPath + string(os.PathSeparator) + fileName

	if !fs.isExists(filePath) {
		err = errors.New("cannot delete file, file not found")
	} else {
		err = fs.Delete(filePath)
	}
	return filePath, err
}

func (fs *FileSystem) CopyCmd(srcDirPath, fileName, destDirPath string) (string, error) {
	var err error
	var oldFilePath = srcDirPath + string(os.PathSeparator) + fileName
	var newFilePath = destDirPath + string(os.PathSeparator) + fileName

	if !fs.isExists(oldFilePath) {
		err = errors.New("cannot copy file, file not found")
	} else {
		err = fs.Copy(oldFilePath, newFilePath)
	}
	return newFilePath, err
}

func (fs *FileSystem) Ls(dirPath string) ([]string, error) {
	var fileNames []string
	files, err := os.ReadDir(fs.getPath(dirPath))
	if err != nil {
		return fileNames, err
	}
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	return fileNames, nil
}

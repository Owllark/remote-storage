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

	err = fs.Mkdir(newDirPath, 0644)
	if err != nil {
		switch {
		case os.IsNotExist(err):
			err = errors.New(fmt.Sprintf("invalid path"))
		case os.IsExist(err):
			err = errors.New(fmt.Sprintf("directory #{newDirPath} already exist"))
		default:
			err = errors.New(fmt.Sprintf("error creating directory"))
		}
	}

	return newDirPath, err
}

func (fs *FileSystem) RenameCmd(dirPath, oldName, newName string) (string, error) {
	var err error
	var oldFilePath = dirPath + oldName
	var newFilePath = dirPath + newName

	err = fs.Rename(oldFilePath, newFilePath)
	if err != nil {
		switch {
		case os.IsNotExist(err):
			err = errors.New(fmt.Sprintf("#{oldFilePath} not found"))
		case os.IsExist(err):
			err = errors.New(fmt.Sprintf("#{newFilePath} already exist"))
		default:
			err = errors.New(fmt.Sprintf("error renaming file"))
		}
	}

	return newFilePath, err
}

func (fs *FileSystem) MoveCmd(srcDirPath, fileName, destDirPath string) (string, error) {
	var err error
	var oldFilePath = srcDirPath + fileName
	var newFilePath = destDirPath + fileName

	err = fs.Move(oldFilePath, newFilePath)
	if err != nil {
		switch {
		case os.IsNotExist(err):
			err = errors.New(fmt.Sprintf("#{oldFilePath} not found"))
		case os.IsExist(err):
			err = errors.New(fmt.Sprintf("#{newFilePath} already exist"))
		default:
			err = errors.New(fmt.Sprintf("error moving file"))
		}
	}

	return newFilePath, err
}

func (fs *FileSystem) DeleteCmd(dirPath, fileName string) (string, error) {
	var err error
	var filePath = dirPath + fileName

	err = fs.Delete(filePath)
	if err != nil {
		switch {
		case os.IsNotExist(err):
			err = errors.New(fmt.Sprintf("#{filePath} not found"))
		default:
			err = errors.New(fmt.Sprintf("error deleting file"))
		}
	}

	return filePath, err
}

func (fs *FileSystem) CopyCmd(srcDirPath, fileName, destDirPath string) (string, error) {
	var err error
	var oldFilePath = srcDirPath + fileName
	var newFilePath = destDirPath + fileName

	err = fs.Copy(oldFilePath, newFilePath)
	if err != nil {
		switch {
		case os.IsNotExist(err):
			err = errors.New(fmt.Sprintf("#{oldFilePath} not found"))
		case os.IsExist(err):
			err = errors.New(fmt.Sprintf("#{newFilePath} already exist"))
		default:
			err = errors.New(fmt.Sprintf("error copying file"))
		}
	}

	return newFilePath, err
}

func (fs *FileSystem) Ls(dirPath string) ([]string, error) {
	var fileNames []string
	files, err := os.ReadDir(fs.getPath(dirPath))
	if err != nil {
		if err != nil {
			switch {
			case os.IsNotExist(err):
				err = errors.New(fmt.Sprintf("#{oldFilePath} not found"))
			default:
				err = errors.New(fmt.Sprintf("error listing files of directory"))
			}
		}
		return fileNames, err
	}
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	return fileNames, nil
}

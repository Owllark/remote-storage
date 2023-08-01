package file_system_svc

//
//import (
//	"errors"
//	"fmt"
//	"os"
//	"path/filepath"
//	"remote-storage/common"
//	"sort"
//	"strconv"
//	"sync"
//)
//
//type FileSystem struct {
//	rootDir string
//}
//
//var (
//	fileSystemSingleInstance *FileSystem
//	mutex                    = sync.Mutex{}
//)
//
//func New(rootDir string) *FileSystem {
//	if fileSystemSingleInstance == nil {
//		mutex.Lock()
//		defer mutex.Unlock()
//		if fileSystemSingleInstance == nil {
//			fileSystemSingleInstance = &FileSystem{
//				rootDir,
//			}
//		}
//
//	} else {
//		return fileSystemSingleInstance
//	}
//	return fileSystemSingleInstance
//}
//
//func (fs *FileSystem) GetFileSystemState(userRootDir string) (common.FileInfo, error) {
//	res, err := fs.TraverseDirectory(userRootDir)
//	return res, err
//}
//
//func (fs *FileSystem) MkDir(path string, dir string) (string, error) {
//	var err error
//	var newDirPath string
//	newDirPath = path + dir
//
//	err = fs.Mkdir(newDirPath, 0644)
//	if err != nil {
//		switch {
//		case os.IsNotExist(err):
//			err = errors.New(fmt.Sprintf("invalid path"))
//		case os.IsExist(err):
//			err = errors.New(fmt.Sprintf("directory %s already exist", newDirPath))
//		default:
//			err = errors.New(fmt.Sprintf("error creating directory"))
//		}
//	}
//
//	return newDirPath, err
//}
//
//func (fs *FileSystem) RenameCmd(dirPath, oldName, newName string) (string, error) {
//	var err error
//	var oldFilePath = dirPath + oldName
//	var newFilePath = dirPath + newName
//
//	err = fs.Rename(oldFilePath, newFilePath)
//	if err != nil {
//		switch {
//		case os.IsNotExist(err):
//			err = errors.New(fmt.Sprintf("%s not found", oldFilePath))
//		case os.IsExist(err):
//			err = errors.New(fmt.Sprintf("%s already exist", newFilePath))
//		default:
//			err = errors.New(fmt.Sprintf("error renaming file"))
//		}
//	}
//
//	return newFilePath, err
//}
//
//func (fs *FileSystem) MoveCmd(srcDirPath, fileName, destDirPath string) (string, error) {
//	var err error
//	var oldFilePath = srcDirPath + fileName
//	var newFilePath = destDirPath + fileName
//
//	err = fs.Move(oldFilePath, newFilePath)
//	if err != nil {
//		switch {
//		case os.IsNotExist(err):
//			err = errors.New(fmt.Sprintf("%s not found", oldFilePath))
//		case os.IsExist(err):
//			err = errors.New(fmt.Sprintf("%s already exist", newFilePath))
//		default:
//			err = errors.New(fmt.Sprintf("error moving file"))
//		}
//	}
//
//	return newFilePath, err
//}
//
//func (fs *FileSystem) DeleteCmd(dirPath, fileName string) (string, error) {
//	var err error
//	var filePath = dirPath + fileName
//
//	err = fs.RemoveAll(filePath)
//	if err != nil {
//		switch {
//		case os.IsNotExist(err):
//			err = errors.New(fmt.Sprintf("%s not found", filePath))
//		default:
//			err = errors.New(fmt.Sprintf("error deleting file"))
//		}
//	}
//
//	return filePath, err
//}
//
//func (fs *FileSystem) CopyCmd(srcDirPath, fileName, destDirPath string) (string, error) {
//	var err error
//	var oldFilePath = srcDirPath + fileName
//	var newFilePath = destDirPath + fileName
//
//	err = fs.Copy(oldFilePath, newFilePath)
//	if err != nil {
//		switch {
//		case os.IsNotExist(err):
//			err = errors.New(fmt.Sprintf("%s not found", oldFilePath))
//		case os.IsExist(err):
//			err = errors.New(fmt.Sprintf("%s already exist", newFilePath))
//		default:
//			err = errors.New(fmt.Sprintf("error copying file"))
//		}
//	}
//
//	return newFilePath, err
//}
//
//func (fs *FileSystem) AssembleFiles(location, tempDirPath, fileName string) error {
//	var err error
//	tempDirPath = fs.getPath(tempDirPath)
//	files, err := os.ReadDir(tempDirPath)
//	if err != nil {
//		return err
//	}
//
//	// Sort files by name in ascending order
//	sort.Slice(files, func(i, j int) bool {
//		fileNameI, _ := strconv.Atoi(files[i].Name())
//		fileNameJ, _ := strconv.Atoi(files[j].Name())
//		return fileNameI < fileNameJ
//	})
//
//	// Create the output file
//	outputFile, err := fs.Create(filepath.Join(location, fileName))
//	if err != nil {
//		return err
//	}
//
//	// Iterate over the files and append their contents to the output file
//	for _, file := range files {
//		filePath := filepath.Join(tempDirPath, file.Name())
//		fileContent, err := os.ReadFile(filePath)
//		if err != nil {
//			return err
//		}
//
//		_, err = outputFile.Write(fileContent)
//		if err != nil {
//			return err
//		}
//	}
//	defer outputFile.Close()
//	return nil
//}
//
//func (fs *FileSystem) SetRootDir(dir string) {
//	fs.rootDir = dir
//}

package file_system_svc

import (
	"errors"
	"github.com/go-kit/kit/log"
	"io"
	"os"
	"server/authsvc"
	"server/authsvc/client"
	fs "server/file_system_svc/repository/filesystem"
)

type FileSystemService interface {
	GetState(userRootDir string) (fs.FileInfo, error)
	MkDir(path string, dir string) (string, error)
	Rename(dirPath, oldName, newName string) (string, error)
	Move(srcDirPath, fileName, destDirPath string) (string, error)
	Delete(dirPath, fileName string) (string, error)
	Copy(srcDirPath, fileName, destDirPath string) (string, error)
	Download(dirPath, fileName string) (io.ReadCloser, error)
	Upload(dirPath, fileName string, contents io.ReadCloser) error
	getAuthSvc() authsvc.Service
}

var (
	ErrUnknownError  = errors.New("unknown error")
	ErrAlreadyExists = errors.New("already exists")
	ErrNotFound      = errors.New("not found")
	ErrAuthFailed    = errors.New("authentication failed")
)

type Config struct {
	RootDir             string
	ConsulServerAddress string
}

type fileSystemService struct {
	authSvc authsvc.Service
}

func NewFileSystemService(logger log.Logger, config Config) FileSystemService {
	consulAddr := config.ConsulServerAddress
	authSvc, err := client.New(consulAddr, logger)
	if err != nil {
		logger.Log("Error:", "cannot create authentication service client")
	}
	fs.InitializeFileSystem(fs.ConfigFileSystem{
		RootDir: config.RootDir,
	})
	return &fileSystemService{
		authSvc: authSvc,
	}
}

func (svc *fileSystemService) getAuthSvc() authsvc.Service {
	return svc.authSvc
}
func (svc *fileSystemService) GetState(userRootDir string) (fs.FileInfo, error) {
	res, err := fs.TraverseDirectory(userRootDir)
	if err != nil {
		err = getErrorType(err)
	}
	return res, err
}

func (svc *fileSystemService) MkDir(path string, dir string) (string, error) {
	var err error
	var newDirPath string
	newDirPath = path + dir

	err = fs.Mkdir(newDirPath, 0644)
	if err != nil {
		err = getErrorType(err)
	}

	return newDirPath, err
}

func (svc *fileSystemService) Rename(dirPath, oldName, newName string) (string, error) {
	var err error
	var oldFilePath = dirPath + oldName
	var newFilePath = dirPath + newName

	err = fs.Rename(oldFilePath, newFilePath)
	if err != nil {
		err = getErrorType(err)
	}

	return newFilePath, err
}

func (svc *fileSystemService) Move(srcDirPath, fileName, destDirPath string) (string, error) {
	var err error
	var oldFilePath = srcDirPath + fileName
	var newFilePath = destDirPath + fileName

	err = fs.Move(oldFilePath, newFilePath)
	if err != nil {
		err = getErrorType(err)
	}

	return newFilePath, err
}

func (svc *fileSystemService) Delete(dirPath, fileName string) (string, error) {
	var err error
	var filePath = dirPath + fileName

	err = fs.RemoveAll(filePath)
	if err != nil {
		err = getErrorType(err)
	}

	return filePath, err
}

func (svc *fileSystemService) Copy(srcDirPath, fileName, destDirPath string) (string, error) {
	var err error
	var oldFilePath = srcDirPath + fileName
	var newFilePath = destDirPath + fileName

	err = fs.Copy(oldFilePath, newFilePath)
	if err != nil {
		err = getErrorType(err)
	}

	return newFilePath, err
}

func (svc *fileSystemService) Download(dirPath, fileName string) (io.ReadCloser, error) {
	var err error
	var filePath = dirPath + fileName

	file, err := fs.OpenFile(filePath)
	if err != nil {
		err = getErrorType(err)
	}

	return file, err
}

func (svc *fileSystemService) Upload(dirPath, fileName string, contents io.ReadCloser) error {
	var err error
	var filePath = dirPath + fileName

	file, err := fs.Create(filePath)
	if err != nil {
		return getErrorType(err)
	}
	_, err = io.Copy(file, contents)
	if err != nil {
		return getErrorType(err)
	}
	return nil
}

func getErrorType(err error) error {
	switch {
	case os.IsNotExist(err):
		err = ErrNotFound
	case os.IsExist(err):
		err = ErrAlreadyExists
	default:
		err = ErrUnknownError
	}
	return err
}

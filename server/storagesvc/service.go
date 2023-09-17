package storagesvc

import (
	"context"
	"errors"
	"github.com/go-kit/kit/log"
	"io"
	"os"
	"remote-storage/server/authsvc"
	"remote-storage/server/authsvc/client"
	fs "remote-storage/server/storagesvc/repository/filesystem"
)

// Service interface contains methods for managing local file system
type Service interface {
	GetState(ctx context.Context) (fs.FileInfo, error)
	MkDir(ctx context.Context, dir, path string) (string, error)
	Rename(ctx context.Context, dirPath, oldName, newName string) (string, error)
	Move(ctx context.Context, srcDirPath, fileName, destDirPath string) (string, error)
	Delete(ctx context.Context, dirPath, fileName string) (string, error)
	Copy(ctx context.Context, srcDirPath, fileName, destDirPath string) (string, error)
	Download(ctx context.Context, dirPath, fileName string) (io.ReadCloser, error)
	Upload(ctx context.Context, dirPath, fileName string, contents io.ReadCloser) error
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

type service struct {
	authSvc authsvc.Service
}

func NewFileSystemService(logger log.Logger, config Config) Service {
	consulAddr := config.ConsulServerAddress
	authSvc, err := client.New(consulAddr, logger)
	if err != nil {
		logger.Log("Error:", "cannot create authentication service client")
	}
	fs.InitializeFileSystem(fs.ConfigFileSystem{
		RootDir: config.RootDir,
	})
	return &service{
		authSvc: authSvc,
	}
}

func (svc *service) getAuthSvc() authsvc.Service {
	return svc.authSvc
}
func (svc *service) GetState(ctx context.Context) (fs.FileInfo, error) {
	userRootDir := ctx.Value(ctxUserInfKey{}).(string)
	res, err := fs.TraverseDirectory(userRootDir)
	if err != nil {
		err = getErrorType(err)
	}
	return res, err
}

func (svc *service) MkDir(ctx context.Context, path string, dir string) (string, error) {
	var err error
	var newDirPath string
	userRootDir := ctx.Value(ctxUserInfKey{}).(string)
	newDirPath = userRootDir + path + dir

	err = fs.Mkdir(newDirPath, 0644)
	if err != nil {
		err = getErrorType(err)
	}

	return newDirPath, err
}

func (svc *service) Rename(ctx context.Context, dirPath, oldName, newName string) (string, error) {
	var err error
	userRootDir := ctx.Value(ctxUserInfKey{}).(string)
	var oldFilePath = userRootDir + dirPath + oldName
	var newFilePath = userRootDir + dirPath + newName

	err = fs.Rename(oldFilePath, newFilePath)
	if err != nil {
		err = getErrorType(err)
	}

	return newFilePath, err
}

func (svc *service) Move(ctx context.Context, srcDirPath, fileName, destDirPath string) (string, error) {
	var err error
	userRootDir := ctx.Value(ctxUserInfKey{}).(string)
	var oldFilePath = userRootDir + srcDirPath + fileName
	var newFilePath = userRootDir + destDirPath + fileName

	err = fs.Move(oldFilePath, newFilePath)
	if err != nil {
		err = getErrorType(err)
	}

	return newFilePath, err
}

func (svc *service) Delete(ctx context.Context, dirPath, fileName string) (string, error) {
	var err error
	userRootDir := ctx.Value(ctxUserInfKey{}).(string)
	var filePath = userRootDir + dirPath + fileName

	err = fs.RemoveAll(filePath)
	if err != nil {
		err = getErrorType(err)
	}

	return filePath, err
}

func (svc *service) Copy(ctx context.Context, srcDirPath, fileName, destDirPath string) (string, error) {
	var err error
	userRootDir := ctx.Value(ctxUserInfKey{}).(string)
	var oldFilePath = userRootDir + srcDirPath + fileName
	var newFilePath = userRootDir + destDirPath + fileName

	err = fs.Copy(oldFilePath, newFilePath)
	if err != nil {
		err = getErrorType(err)
	}

	return newFilePath, err
}

func (svc *service) Download(ctx context.Context, dirPath, fileName string) (io.ReadCloser, error) {
	var err error
	userRootDir := ctx.Value(ctxUserInfKey{}).(string)
	var filePath = userRootDir + dirPath + fileName

	file, err := fs.OpenFile(filePath)
	if err != nil {
		err = getErrorType(err)
	}

	return file, err
}

func (svc *service) Upload(ctx context.Context, dirPath, fileName string, contents io.ReadCloser) error {
	var err error
	userRootDir := ctx.Value(ctxUserInfKey{}).(string)
	var filePath = userRootDir + dirPath + fileName

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

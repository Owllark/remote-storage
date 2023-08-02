package file_system_svc

import (
	"github.com/go-kit/kit/log"
	"io"
	fs "server/file_system_svc/repository/filesystem"
	"time"
)

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next FileSystemService) FileSystemService {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   FileSystemService
	logger log.Logger
}

func (mw loggingMiddleware) GetState(userRootDir string) (info fs.FileInfo, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetState", "user root dir", userRootDir, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetState(userRootDir)
}

func (mw loggingMiddleware) MkDir(path string, dir string) (s string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "MkDir", "path", path, "dir", dir, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.MkDir(path, dir)
}

func (mw loggingMiddleware) Rename(dirPath, oldName, newName string) (s string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Rename", "dirPath", dirPath, "oldName", oldName, "newName", newName, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Rename(dirPath, oldName, newName)
}

func (mw loggingMiddleware) Move(srcDirPath, fileName, destDirPath string) (s string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Move", "dirPath", srcDirPath, "fileName", fileName, "destDirPath", destDirPath, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Move(srcDirPath, fileName, destDirPath)
}

func (mw loggingMiddleware) Delete(dirPath, fileName string) (s string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Delete", "dirPath", dirPath, "fileName", fileName, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Delete(dirPath, fileName)
}

func (mw loggingMiddleware) Copy(srcDirPath, fileName, destDirPath string) (s string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Copy", "srcDirPath", srcDirPath, "fileName", fileName, "destDirPath", destDirPath, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Copy(srcDirPath, fileName, destDirPath)
}

func (mw loggingMiddleware) Download(dirPath, fileName string) (buffer io.ReadCloser, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Download", "dirPath", dirPath, "fileName", fileName, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Download(dirPath, fileName)
}

func (mw loggingMiddleware) Upload(dirPath, fileName string, contents io.ReadCloser) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Upload", "dirPath", dirPath, "fileName", fileName, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Upload(dirPath, fileName, contents)
}

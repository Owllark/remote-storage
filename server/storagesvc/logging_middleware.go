package storagesvc

import (
	"context"
	"github.com/go-kit/kit/log"
	"io"
	"remote-storage/server/authsvc"
	fs "remote-storage/server/storagesvc/repository/filesystem"
	"time"
)

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) getAuthSvc() authsvc.Service {
	return mw.next.getAuthSvc()
}

func (mw loggingMiddleware) GetState(ctx context.Context) (info fs.FileInfo, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetState", "user root dir", userRootDir, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetState(nil)
}

func (mw loggingMiddleware) MkDir(ctx context.Context, dir, path string) (s string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "MkDir", "path", path, "dir", dir, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.MkDir(nil, dir, path)
}

func (mw loggingMiddleware) Rename(ctx context.Context, dirPath, oldName, newName string) (s string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Rename", "dirPath", dirPath, "oldName", oldName, "newName", newName, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Rename(ctx, dirPath, oldName, newName)
}

func (mw loggingMiddleware) Move(ctx context.Context, srcDirPath, fileName, destDirPath string) (s string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Move", "dirPath", srcDirPath, "fileName", fileName, "destDirPath", destDirPath, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Move(ctx, srcDirPath, fileName, destDirPath)
}

func (mw loggingMiddleware) Delete(ctx context.Context, dirPath, fileName string) (s string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Delete", "dirPath", dirPath, "fileName", fileName, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Delete(ctx, dirPath, fileName)
}

func (mw loggingMiddleware) Copy(ctx context.Context, srcDirPath, fileName, destDirPath string) (s string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Copy", "srcDirPath", srcDirPath, "fileName", fileName, "destDirPath", destDirPath, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Copy(ctx, srcDirPath, fileName, destDirPath)
}

func (mw loggingMiddleware) Download(ctx context.Context, dirPath, fileName string) (buffer io.ReadCloser, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Download", "dirPath", dirPath, "fileName", fileName, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Download(ctx, dirPath, fileName)
}

func (mw loggingMiddleware) Upload(ctx context.Context, dirPath, fileName string, contents io.ReadCloser) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Upload", "dirPath", dirPath, "fileName", fileName, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Upload(ctx, dirPath, fileName, contents)
}

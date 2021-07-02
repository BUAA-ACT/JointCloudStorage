package util

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

type TempFileStorage struct {
	TempPath       string
	TempFiles      []*TempFile
	ExpireDuration time.Duration
}

func NewTempFileStorage(tempPath string, expireDuration time.Duration) (*TempFileStorage, error) {
	return &TempFileStorage{
		TempPath:       tempPath,
		TempFiles:      []*TempFile{},
		ExpireDuration: expireDuration,
	}, nil
}

type TempFile struct {
	FilePath   string
	File       *os.File
	ExpireTime time.Time
}

// CreateTmpFile 创建一个临时文件，若指定了 key，则文件名以 key 作为结尾
func (ts *TempFileStorage) CreateTmpFile(key string) (*os.File, *TempFile) {
	prefix := uuid.New().String()
	fileName := "/" + prefix
	if key != "" {
		fileName += "_" + key
	}
	filePath := path.Join(ts.TempPath, fileName)
	tmpFile := TempFile{
		FilePath:   filePath,
		ExpireTime: time.Now().Add(ts.ExpireDuration),
	}
	ts.TempFiles = append(ts.TempFiles, &tmpFile)
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666) // 此处假设当前目录下已存在test目录
	if err != nil {
		Log(logrus.ErrorLevel, "create temp file", "cannot open temp file",
			"", "err", err.Error())
	}
	return f, &tmpFile
}

package model

import (
	"errors"
	"time"
)

type File struct {
	Id                string
	Filename          string
	Owner             string
	Size              int64
	LastChange        time.Time
	SyncStatus        string // 同步状态 Pending/Deleting/Done
	ReconstructStatus string // 重建状态
	DownloadUrl       string
	ReconstructTime   time.Time
}

type FileDatabase interface {
	CreateFileInfo(file *File) (err error)
	DeleteFileInfo(file *File) (err error)
	UpdateFileInfo(file *File) (err error)
	GetFileInfo(Id string) (file *File, err error)
}

type InMemoryFileDatabase struct {
	db map[string]File
}

func (fd *InMemoryFileDatabase) CreateFileInfo(file *File) (err error) {
	fd.db[file.Id] = *file
	return nil
}

func (fd *InMemoryFileDatabase) DeleteFileInfo(file *File) (err error) {
	delete(fd.db, file.Id)
	return nil
}

func (fd *InMemoryFileDatabase) UpdateFileInfo(file *File) (err error) {
	if _, ok := fd.db[file.Id]; ok {
		return errors.New("file info not exist")
	}
	fd.db[file.Id] = *file
	return nil
}

func (fd *InMemoryFileDatabase) GetFileInfo(Id string) (file *File, err error) {
	f, ok := fd.db[Id]
	if ok {
		return &f, nil
	}
	return nil, errors.New("file info not exist")
}

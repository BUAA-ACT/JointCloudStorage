package controller

import "act.buaa.edu.cn/jcspan/transporter/model"

type Scheduler interface {
	UploadFileMetadata(clouds []string, uid string, file model.File) error
	DeleteFileMetadata() error
	MigrateFileMetadata() error
}

type JcsPanScheduler struct {
	localCloudID string
}

type JcsPanSchedulerType string

const (
	MetadataUpload  = "Upload"
	MetadataDelete  = "Delete"
	MetadataMigrate = "Migrate"
)

type Metadata struct {
	CloudID string
	UserID  string
	Type    JcsPanSchedulerType
	Files   []model.File
}

func (s *JcsPanScheduler) UploadFileMetadata(clouds []string, uid string, file model.File) error {
	return nil
}

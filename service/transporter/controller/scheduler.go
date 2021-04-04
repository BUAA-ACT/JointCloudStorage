package controller

import (
	"act.buaa.edu.cn/jcspan/transporter/model"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type Scheduler interface {
	UploadFileMetadata(clouds []*model.Cloud, uid string, file *model.File) error
	DeleteFileMetadata(clouds []*model.Cloud, uid string, file *model.File) error
	MigrateFileMetadata(clouds []*model.Cloud, uid string, file *model.File) error
}

type JcsPanScheduler struct {
	LocalCloudID     string
	SchedulerHostUrl string
}

func NewJcsPanScheduler(localCloudID string, schedulerHostUrl string) *JcsPanScheduler {
	return &JcsPanScheduler{
		LocalCloudID:     localCloudID,
		SchedulerHostUrl: schedulerHostUrl,
	}
}

type MetaDataProcessType string

const (
	MetadataUpload  MetaDataProcessType = "Upload"
	MetadataDelete  MetaDataProcessType = "Delete"
	MetadataMigrate MetaDataProcessType = "Migrate"
)

type Metadata struct {
	CloudID string
	UserID  string
	Type    MetaDataProcessType
	Files   []*model.File
	Clouds  []*model.Cloud
}

type MetadataReply struct {
	Code      int
	Msg       string
	RequestID string
}

func (s *JcsPanScheduler) sendFileMetadata(mType MetaDataProcessType,
	clouds []*model.Cloud, uid string, file *model.File) error {
	metaData := Metadata{
		CloudID: s.LocalCloudID,
		UserID:  uid,
		Type:    mType,
		Files: []*model.File{
			file,
		},
		Clouds: clouds,
	}
	b, err := json.Marshal(metaData)
	if err != nil {
		return err
	}
	resp, err := http.Post(s.SchedulerHostUrl+"/metadata", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	var reply MetadataReply
	err = json.NewDecoder(resp.Body).Decode(&reply)
	if err != nil {
		return err
	}
	if reply.Code != http.StatusOK {
		return errors.New(reply.Msg)
	}
	return nil
}

func (s *JcsPanScheduler) UploadFileMetadata(clouds []*model.Cloud, uid string, file *model.File) error {
	return s.sendFileMetadata(MetadataUpload, clouds, uid, file)
}

func (s *JcsPanScheduler) DeleteFileMetadata(clouds []*model.Cloud, uid string, file *model.File) error {
	return s.sendFileMetadata(MetadataDelete, clouds, uid, file)
}

func (s *JcsPanScheduler) MigrateFileMetadata(clouds []*model.Cloud, uid string, file *model.File) error {
	return s.sendFileMetadata(MetadataMigrate, clouds, uid, file)
}

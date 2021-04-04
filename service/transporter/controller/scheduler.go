package controller

import (
	"act.buaa.edu.cn/jcspan/transporter/model"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Scheduler interface {
	UploadFileMetadata(clouds []string, uid string, file *model.File) error
	DeleteFileMetadata(clouds []string, uid string, file *model.File) error
	MigrateFileMetadata(clouds []string, uid string, file *model.File) error
}

type JcsPanScheduler struct {
	LocalCloudID     string
	SchedulerHostUrl string
	ReloadCloudInfo  bool
	CloudDatabase    model.CloudDatabase
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

func (s *JcsPanScheduler) getClouds(cloudsID []string) []*model.Cloud {
	var clouds []*model.Cloud
	for _, cloudID := range cloudsID {
		if s.ReloadCloudInfo && s.CloudDatabase != nil {
			cloud, err := s.CloudDatabase.GetCloudInfoFromCloudID(cloudID)
			if err != nil {
				logrus.Errorf("can't get cloud info: %v", cloudID)
				continue
			}
			clouds = append(clouds, cloud)
		} else {
			clouds = append(clouds, &model.Cloud{
				CloudID: cloudID,
			})
		}
	}
	return clouds
}

func (s *JcsPanScheduler) UploadFileMetadata(cloudsID []string, uid string, file *model.File) error {
	clouds := s.getClouds(cloudsID)
	return s.sendFileMetadata(MetadataUpload, clouds, uid, file)
}

func (s *JcsPanScheduler) DeleteFileMetadata(cloudsID []string, uid string, file *model.File) error {
	clouds := s.getClouds(cloudsID)
	return s.sendFileMetadata(MetadataDelete, clouds, uid, file)
}

func (s *JcsPanScheduler) MigrateFileMetadata(cloudsID []string, uid string, file *model.File) error {
	clouds := s.getClouds(cloudsID)
	return s.sendFileMetadata(MetadataMigrate, clouds, uid, file)
}

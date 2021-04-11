package model

import "time"

type File struct {
	FileID            string    `json:"FileID" bson:"file_id"`
	FileName          string    `json:"FileName" bson:"file_name"`
	Owner             string    `json:"Owner" bson:"owner"`
	Size              uint64    `json:"Size" bson:"size"`
	LastModified      time.Time `json:"LastModified" bson:"last_modified"`
	LastReconstructed time.Time `json:"LastReconstructed" bson:"last_reconstructed"`
	ReconstructStatus string    `json:"ReconstructStatus" bson:"reconstruct_status"`
	DownloadUrl       string    `json:"DownloadUrl,omitempty" bson:"download_url"`
}

type FileAndDir struct {
	FileType string `json:"FileType" bson:"file_type"`
	FileInfo File   `json:"FileInfo" json:"file_info"`
}

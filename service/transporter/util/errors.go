package util

const (
	ErrorCodeInternalErr        = 3400
	ErrorCodeWrongRequestFormat = 3401
	ErrorCodeWrongTaskType      = 3402
	ErrorCodeWrongStorageType   = 3403
	ErrorCodeGetFileLockErr     = 3404
)

const (
	ErrorMsgWrongRequestFormat        = "wrong request format"
	ErrorMsgWrongTaskType             = "task type not implement"
	ErrorMsgWrongStorageType          = "storage type not implement"
	ErrorMsgGetFileLockErr            = "can't get file lock"
	ErrorMsgCantGetFileInfo           = "can't get file info"
	ErrorMsgWrongCloudNum             = "clouds num miss match"
	ErrorMsgProcessMigrateDownloadErr = "process migrate download err"
	ErrorMsgProcessMigrateUploadErr   = "process migrate Upload err"
	ErrorMsgEmptyFilename             = "empty file name"
)

package args

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

/* error code*/
const (
	// OK
	CodeOK = 200
	// BadRequest
	CodeUploadError            = 1400
	CodeVerifyFail             = 1401
	CodeInvalidToken           = 1402
	CodeFileNotExists          = 1403
	CodeEmailNotExist          = 1404
	CodeSameEmail              = 1405
	CodeFieldNotExist          = 1406
	CodeRegexWrong             = 1407
	CodePasswordNotRight       = 1408
	CodeTokenNotValid          = 1409
	CodePreferenceNotExist     = 1410
	CodeStoragePlanNotExist    = 1411
	CodeJsonError              = 1412
	CodeDifferentUser          = 1413
	CodeFileNotExist           = 1414
	CodeAlreadyHaveStoragePlan = 1415
	CodeForbiddenTransport     = 1416

	// user status code
	CodeStatusNormal       = 1501
	CodeStatusForbidden    = 1502
	CodeStatusVerify       = 1503
	CodeStatusTransporting = 1504

	// InternalError
	CodeInternalError    = 1000
	CodeDatabaseError    = 2000
	CodeSchedulerError   = 1300
	CodeTransporterError = 1500
)

/* http request method */
const (
	HttpMethodGet             = "GET"
	HttpMethodPost            = "POST"
	HttpContentTypeUrlEncoded = "application/x-www-form-urlencoded"
	HttpContentTypeJson       = "application/json"
	HttpContentTypeRaw        = "text/plain"
	HttpContentTypeDataForm   = "multipart/form-data"
)

/* http body field const */
const (
	FieldWordAccessToken    = "AccessToken"
	FieldWordEmail          = "Email"
	FieldWordPassword       = "Password"
	FieldWordNickname       = "Nickname"
	FieldWordVerifyCode     = "VerifyCode"
	FieldWordNewEmail       = "NewEmail"
	FieldWordNewNickname    = "NewNickname"
	FieldWordOriginPassword = "OriginPassword"
	FieldWordNewPassword    = "NewPassword"
	FieldWordVendor         = "Vendor"
	FieldWordStoragePrice   = "StoragePrice"
	FieldWordTrafficPrice   = "TrafficPrice"
	FieldWordAvailability   = "Availability"
	FieldWordLatency        = "Latency"
	FieldWordFilePath       = "FilePath"
	FieldWordOriginFilePath = "OriginFilePath"
	FieldWordNewFilePath    = "NewFilePath"
	FieldWordFileID         = "FileID"
	FieldWordOriginFileName = "OriginFileName"
	FieldWordNewFileName    = "NewFileName"
	FieldWordStoragePlan    = "StoragePlan"
	FieldWordTaskID         = "TaskID"
)

/* user const */
const (
	/* user role */
	UserAdminRole    = "Admin"
	UserSuperRole    = "Super"
	UserOrdinaryRole = "Ordinary"
	UserHostRole     = "HOST"
	UserGuestRole    = "GUEST"

	/* user status */
	UserVerifyStatus       = "VERIFYING"
	UserForbiddenStatus    = "FORBIDDEN"
	UserNormalStatus       = "NORMAL"
	UserTransportingStatus = "Transporting"
)

/* task const*/
const (
	TaskTypeUpload   = "Upload"
	TaskTypeDownload = "Download"
	TaskTypeSync     = "Migrate"
	TaskTypeDelete   = "Delete"

	TaskStateCreated   = "Created"
	TaskStatePending   = "Pending"
	TaskStateExecuting = "Executing"
	TaskStateFailed    = "Failed"
	TaskStateDone      = "Done"
)

/* file const */
const (
	FileTypeDir  = "DIR"
	FileTypeFile = "FILE"

	FileReconstructStatusPending = "Pending"
	FileReconstructStatusWorking = "Working"
	FileReconstructStatusDone    = "Done"
)

/* properties */
var (
	properties                = make(map[string]string)
	MongoTitle                = flag.String("MongoTitle", "", "mongo prefix")
	MongoUsername             = flag.String("MongoUsername", "", "mongo database admin username")
	MongoPassword             = flag.String("MongoPassword", "", "mongo database admin password")
	MongoAddr                 = flag.String("MongoAddr", "", "mongo database ip address")
	MongoPort                 = flag.Uint64("MongoPort", 0, "mongodb server port")
	MongoURL                  = flag.String("Mongo", "", "mongodb server address")
	SchedulerAddr             = flag.String("SchedulerAddress", "", "scheduler address")
	SchedulerPort             = flag.String("SchedulerPort", "", "scheduler port")
	SchedulerUrl              = flag.String("SchedulerUrl", "", "whole scheduler url")
	TransporterAddr           = flag.String("TransporterAddress", "", "transporter address")
	TransporterPort           = flag.String("TransporterPort", "", "transporter port")
	TransporterUrl            = flag.String("TransporterUrl", "", "whole transporter url")
	HttpserverPort            = flag.Uint64("HttpserverPort", 0, "http server port")
	Debug                     = flag.Bool("Debug", false, "debug mode")
	TestMode                  = flag.Bool("Test", false, "enable test mode")
	DataBase                  = flag.String("DatabaseName", "", "httpserver's database name")
	UserCollection            = flag.String("UserCollection", "", "users collection name")
	FileCollection            = flag.String("FileCollection", "", "file collection name")
	AccessTokenCollection     = flag.String("AccessTokenCollection", "", "token collection name")
	TaskCollection            = flag.String("TaskCollection", "", "task collection name")
	VerifyCodeCollection      = flag.String("VerifyCodeCollection", "", "verify code collection name")
	MigrationAdviceCollection = flag.String("MigrationAdviceCollection", "", "migration advice collection name")
	EncryptKey                = flag.String("EncryptKey", "", "password encrypt key for AES")
	CloudID                   = flag.String("CloudID", "", "server's cloud id")
)

func LoadProperties(propertiesFilePath string) {
	if propertiesFilePath == "" {
		propertiesFilePath = "./httpserver.properties"
	}
	srcFile, err := os.OpenFile(propertiesFilePath, os.O_RDONLY, 0666)
	defer srcFile.Close()
	if err != nil {
		fmt.Println("The file not exits.")
	} else {
		reg := regexp.MustCompile("\\s+")
		srcReader := bufio.NewReader(srcFile)
		for {
			str, err := srcReader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
			}
			str = reg.ReplaceAllString(str, "")
			if len(str) == 0 {
				continue
			}
			pair := strings.Split(str, "=")
			key := pair[0]
			value := pair[1]
			properties[key] = value
		}
	}
	*MongoTitle = properties["MongoTitle"]
	*MongoUsername = properties["MongoUsername"]
	*MongoPassword = properties["MongoPassword"]
	*MongoAddr = properties["MongoAddr"]
	*MongoPort, err = strconv.ParseUint(properties["MongoPort"], 10, 64)
	if err != nil {
		*MongoPort = 27017
	}
	*Debug, err = strconv.ParseBool(properties["Debug"])
	if err != nil {
		*Debug = false
	}
	*TestMode, err = strconv.ParseBool(properties["TestMode"])
	if err != nil {
		*TestMode = false
	}
	if *TestMode {
		*MongoURL = *MongoTitle + *MongoAddr + ":" + strconv.FormatUint(*MongoPort, 10)
	} else {
		*MongoURL = *MongoTitle + *MongoUsername + ":" + *MongoPassword + "@" + *MongoAddr + ":" + strconv.FormatUint(*MongoPort, 10)
	}

	*SchedulerAddr = properties["SchedulerAddr"]
	*SchedulerPort = properties["SchedulerPort"]
	*SchedulerUrl = "http://" + *SchedulerAddr + ":" + *SchedulerPort
	*TransporterAddr = properties["TransporterAddr"]
	*TransporterPort = properties["TransporterPort"]
	*TransporterUrl = "http://" + *TransporterAddr + ":" + *TransporterPort
	*HttpserverPort, err = strconv.ParseUint(properties["HttpserverPort"], 10, 64)
	if err != nil {
		*HttpserverPort = 8081
	}
	*DataBase = properties["DataBase"]
	*UserCollection = properties["UserCollection"]
	*FileCollection = properties["FileCollection"]
	*AccessTokenCollection = properties["AccessTokenCollection"]
	*VerifyCodeCollection = properties["VerifyCodeCollection"]
	*TaskCollection = properties["TaskCollection"]
	*MigrationAdviceCollection = properties["MigrationAdviceCollection"]
	*EncryptKey = properties["EncryptKey"]
	*CloudID = properties["CloudID"]
}

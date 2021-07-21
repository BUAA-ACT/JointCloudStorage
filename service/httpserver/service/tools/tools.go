package tools

import (
	"cloud-storage-httpserver/args"
	"log"
	"strings"
)

const (
	CONNECT = 0
	DIGITAL = 1
	LOWER   = 2
	UPPER   = 3
	SPECIAL = 4
)

/*
	Function to turn string in MIMEHeader form
	accesstoken -> Accesstoken
	accessToken -> Access-Token
	access-token -> Access-Token
	AccessToken -> Access-Token
	u1s23e4r5I67D89 -> U1s23e4r5-I67-D89
	CloudID -> Cloud-ID
	a_B&f^D--cn -> a_B&f^D--Cn
	a-fWord -> A-F-Word
*/

func classifyChar(c uint8) int {
	if '0' <= c && c <= '9' {
		return DIGITAL
	} else if 'a' <= c && c <= 'z' {
		return LOWER
	} else if 'A' <= c && c <= 'Z' {
		return UPPER
	} else if c == '-' {
		return CONNECT
	} else {
		return SPECIAL
	}
}

func changeCharUpLow(c uint8) string {
	if classifyChar(c) == LOWER {
		return string(c - 'a' + 'A')
	} else if classifyChar(c) == UPPER {
		return string(c - 'A' + 'a')
	} else {
		return string(c)
	}
}

func CanonicalMIMEHeaderKey(key string) string {
	var oldKind, newKind int
	oldKind = CONNECT
	for i := 0; i < len(key); i++ {
		newKind = classifyChar(key[i])
		if oldKind == DIGITAL && newKind == UPPER {
			// 5I67D89 -> 5-I67-D89
			key = key[:i] + "-" + key[i:]
		} else if oldKind == LOWER && newKind == UPPER {
			// AccessToken -> Access-Token
			key = key[:i] + "-" + key[i:]
		} else if oldKind == CONNECT && newKind == LOWER {
			// access-token -> Access-Token
			key = key[:i] + changeCharUpLow(key[i]) + key[i+1:]
		}
		oldKind = newKind
	}
	return key
}

/*
	Function of File&DIR Name
*/
func FileIdToUserAndFileName(fileId string) (string, string) {
	segments := strings.Split(fileId, "/")
	fileName := strings.TrimPrefix(fileId, segments[0])
	return segments[0], fileName
}

func IsDir(path string) bool {
	return path[len(path)-1:] == "/"
}

// RequiredFieldNotExist Field of the http request
func RequiredFieldNotExist(requiredMap *map[string]bool, existMap *map[string]bool) bool {
	if requiredMap == nil || existMap == nil {
		return true
	}
	for field, required := range *requiredMap {
		if required && !(*existMap)[field] {
			return true
		}
	}
	return false
}

// PrintError put out the err out if has error
func PrintError(err error) bool {
	if err != nil {
		log.Println("出错啦！错误信息为：")
		log.Println(err)
		return true
	}
	return false
}

// UserStatusMessageCode turn user status into error code
func UserStatusMessageCode(status string) (int, string) {
	switch status {
	case args.UserForbiddenStatus:
		return args.CodeStatusForbidden, "用户正在迁移"
	case args.UserNormalStatus:
		return args.CodeStatusNormal, "用户状态正常,但就是有点不正常"
	case args.UserVerifyStatus:
		return args.CodeStatusVerify, "用户邮箱未进行验证"
	case args.UserTransportingStatus:
		return args.CodeStatusTransporting, "暂未定义"
	default:
		return args.CodeOK, "没毛病"
	}
}

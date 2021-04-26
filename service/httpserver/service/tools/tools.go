package tools

import (
	"cloud-storage-httpserver/args"
	"fmt"
	"strings"
)

func FileIdToUserAndFileName(fileId string) (string, string) {
	segments := strings.Split(fileId, "/")
	fileName := strings.TrimPrefix(fileId, segments[0])
	return segments[0], fileName
}

func IsDir(path string) bool {
	return path[len(path)-1:] == "/"
}

func RequiredFieldNotExist(requiredMap *map[string]bool, existMap *map[string]bool) bool {
	for field, required := range *requiredMap {
		if required && !(*existMap)[field] {
			return true
		}
	}
	return false
}

func PrintError(err error) bool {
	if err != nil {
		fmt.Println("出错啦！错误信息为：")
		fmt.Println(err)
		return true
	}
	return false
}

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

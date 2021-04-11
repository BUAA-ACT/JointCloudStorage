package tools

import (
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

func RequiredFieldNotExist(requiredMap map[string]bool, existMap map[string]bool) bool {
	for field, required := range requiredMap {
		if required && !existMap[field] {
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

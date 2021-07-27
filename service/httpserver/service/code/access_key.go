package code

import "strings"

func GenAccessKey() string {
	newId := genUUIDv4()
	return strings.ReplaceAll(newId.String(), "-", "")
}

func GenSecretKey() string {
	newId := genUUIDv4()
	return strings.ReplaceAll(newId.String(), "-", "")
}

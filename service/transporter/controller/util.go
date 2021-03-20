package controller

import "github.com/sirupsen/logrus"

func CheckErr(err error, label string) bool {
	if err != nil {
		logrus.Warnf("%v ERR: %v", label, err)
		return true
	}
	return false
}

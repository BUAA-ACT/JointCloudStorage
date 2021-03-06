package controller

import (
	"act.buaa.edu.cn/jcspan/transporter/model"
	"act.buaa.edu.cn/jcspan/transporter/util"
	"github.com/alibaba/pouch/pkg/kmutex"
	"github.com/sirupsen/logrus"
	"os"
)

type TrafficMonitor struct {
	UserDB model.UserDatabase
	kMutex *kmutex.KMutex
}

func NewTrafficMonitor(userDB model.UserDatabase) *TrafficMonitor {
	return &TrafficMonitor{
		UserDB: userDB,
		kMutex: kmutex.New(),
	}
}

func (monitor TrafficMonitor) getUserAndLock(uid string) (user *model.User, err error) {
	monitor.kMutex.Lock(uid)
	user, err = monitor.UserDB.GetUserFromID(uid)
	return
}

func (monitor TrafficMonitor) updateUserAndUnlock(user *model.User) (err error) {
	err = monitor.UserDB.UpdateUserInfo(user)
	monitor.kMutex.Unlock(user.UserId)
	return
}

func (monitor TrafficMonitor) unlock(uid string) {
	monitor.kMutex.Unlock(uid)
}

func (monitor TrafficMonitor) AddVolume(uid string, delta int64) (size int64, err error) {
	user, err := monitor.getUserAndLock(uid)
	if err != nil {
		util.Log(logrus.ErrorLevel, "TrafficMonitor AddVolume", "can't get user",
			"a user info", "nil", err.Error())
		monitor.unlock(uid)
		return 0, err
	}
	user.DataStats.Volume = user.DataStats.Volume + delta
	err = monitor.updateUserAndUnlock(user)
	if err != nil {
		util.Log(logrus.ErrorLevel, "TrafficMonitor AddVolume", "can't set and unlock user",
			"", "err", err.Error())
		return 0, err
	}
	return user.DataStats.Volume, nil
}

func (monitor TrafficMonitor) AddUploadTrafficFromFile(uid string, path string, cloudID string) {
	go func() {
		fi, err := os.Stat(path)
		if err != nil {
			util.Log(logrus.ErrorLevel,
				"TrafficMonitor AddUploadFromFile",
				"can't get file state", "", "", err.Error())
			return
		}
		_, _ = monitor.AddUploadTraffic(uid, fi.Size(), cloudID)
	}()
}

func (monitor TrafficMonitor) AddUploadTraffic(uid string, delta int64, cloudID string) (size int64, err error) {
	user, err := monitor.getUserAndLock(uid)
	if err != nil {
		util.Log(logrus.ErrorLevel, "TrafficMonitor AddUploadTraffic", "can't get user",
			"a user info", "nil", err.Error())
		monitor.unlock(uid)
		return 0, err
	}
	if user.DataStats.UploadTraffic == nil {
		user.DataStats.UploadTraffic = map[string]int64{}
	}
	uploadTraffic := user.DataStats.UploadTraffic[cloudID]
	user.DataStats.UploadTraffic[cloudID] = uploadTraffic + delta
	err = monitor.updateUserAndUnlock(user)
	if err != nil {
		util.Log(logrus.ErrorLevel, "TrafficMonitor AddUploadTraffic", "can't set and unlock user",
			"", "err", err.Error())
		return 0, err
	}
	return user.DataStats.UploadTraffic[cloudID], nil
}
func (monitor TrafficMonitor) AddDownloadTrafficFromFile(uid string, path string, cloudID string) {
	go func() {
		fi, err := os.Stat(path)
		if err != nil {
			util.Log(logrus.ErrorLevel,
				"TrafficMonitor AddDownloadFromFile",
				"can't get file state", "", "", err.Error())
			return
		}
		_, _ = monitor.AddDownloadTraffic(uid, fi.Size(), cloudID)
	}()
}
func (monitor TrafficMonitor) AddDownloadTraffic(uid string, delta int64, cloudID string) (size int64, err error) {
	user, err := monitor.getUserAndLock(uid)
	if err != nil {
		util.Log(logrus.ErrorLevel, "TrafficMonitor AddDownloadTraffic", "can't get user",
			"a user info", "nil", err.Error())
		defer monitor.unlock(uid)
		return 0, err
	}
	if user.DataStats.DownloadTraffic == nil {
		user.DataStats.DownloadTraffic = map[string]int64{}
	}
	downloadTraffic := user.DataStats.DownloadTraffic[cloudID]
	user.DataStats.DownloadTraffic[cloudID] = downloadTraffic + delta
	err = monitor.updateUserAndUnlock(user)
	if err != nil {
		util.Log(logrus.ErrorLevel, "TrafficMonitor AddDownloadTraffic", "can't set and unlock user",
			"", "err", err.Error())
		return 0, err
	}
	return user.DataStats.DownloadTraffic[cloudID], nil
}

func (monitor TrafficMonitor) ReduceVolume(uid string, delta int64) (size int64, err error) {
	user, err := monitor.getUserAndLock(uid)
	if err != nil {
		util.Log(logrus.ErrorLevel, "TrafficMonitor ReduceVolume", "can't get user",
			"a user info", "nil", err.Error())
		defer monitor.unlock(uid)
		return 0, err
	}
	user.DataStats.Volume = user.DataStats.Volume - delta
	err = monitor.updateUserAndUnlock(user)
	if err != nil {
		util.Log(logrus.ErrorLevel, "TrafficMonitor ReduceVolume", "can't set and unlock user",
			"", "err", err.Error())
		return 0, err
	}
	return user.DataStats.Volume, nil
}

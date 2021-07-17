package dao

import (
	"github.com/sirupsen/logrus"
	"shaoliyin.me/jcspan/config"
	"sync"
)

var globalDao *Dao
var once sync.Once

type Database struct {
	*Dao
}

func GetDatabaseInstance() Database {
	if globalDao == nil {
		conf := config.GetConfig()
		once.Do(func() {
			dao, err := NewDao(conf.FlagMongo, conf.FlagEnv,
				conf.CollectionNames.CollectionCloud, conf.CollectionNames.CollectionUser,
				conf.CollectionNames.CollectionFile, conf.CollectionNames.MigrationAdvice,
				conf.CollectionNames.CollectionAK)
			if err != nil {
				logrus.Errorf("创建 Dao 失败： %v", err)
			}
			globalDao = dao
			logrus.Infof("创建全局 Dao 实例成功")
		})
	}
	return Database{globalDao}
}

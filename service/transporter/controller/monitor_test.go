package controller

import (
	"act.buaa.edu.cn/jcspan/transporter/model"
	"act.buaa.edu.cn/jcspan/transporter/util"
	"sync"
	"testing"
)

func TestTrafficMonitor(t *testing.T) {
	var monitor *TrafficMonitor
	var db model.UserDatabase
	//util.Config.Database.Driver = util.InMemoryDB
	if util.Config.Database.Driver == util.InMemoryDB {
		db = model.NewInMemoryUserDatabase()
		monitor = NewTrafficMonitor(db)
	} else if util.Config.Database.Driver == util.MongoDB {
		var err error
		db, err = model.NewMongoUserDatabase()
		if err != nil {
			t.Fatalf("create db error: %v", err)
		}
		user, _ := db.GetUserFromID("tester")
		user.DataStats.Volume = 0
		user.DataStats.UploadTraffic = map[string]int64{}
		user.DataStats.DownloadTraffic = map[string]int64{}
		db.UpdateUserInfo(user)
		monitor = NewTrafficMonitor(db)
	}

	t.Run("test AddVolume", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < 20; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := monitor.AddVolume("tester", 100)
				_, err = monitor.AddVolume("tester", 100)
				_, err = monitor.AddVolume("tester", 120)
				if err != nil {
					t.Fatalf("test AddVolum Fail: %v", err)
				}
			}()
		}
		wg.Wait()
		user, _ := db.GetUserFromID("tester")
		if user.DataStats.Volume != 320*20 {
			t.Fatalf("teat AddVolum fail, expect 32000, got %v", user.DataStats.Volume)
		}
	})
	t.Run("test AddUpload", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < 20; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := monitor.AddUploadTraffic("tester", 100, "test")
				_, err = monitor.AddUploadTraffic("tester", 100, "test")
				_, err = monitor.AddUploadTraffic("tester", 120, "test")
				if err != nil {
					t.Fatalf("test AddUpload Fail: %v", err)
				}
			}()
		}
		wg.Wait()
		user, _ := db.GetUserFromID("tester")
		if user.DataStats.UploadTraffic["test"] != 320*20 {
			t.Fatalf("teat AddVolum fail, expect 3200, got %v", user.DataStats.UploadTraffic["test"])
		}
	})
}

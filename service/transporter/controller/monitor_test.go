package controller

import (
	"act.buaa.edu.cn/jcspan/transporter/model"
	"act.buaa.edu.cn/jcspan/transporter/util"
	"sync"
	"testing"
)

func TestTrafficMonitor(t *testing.T) {
	db := model.NewInMemoryUserDatabase()
	monitor := NewTrafficMonitor(db)
	t.Run("test AddVolume", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				_, err := monitor.AddVolume("tester", 100)
				_, err = monitor.AddVolume("tester", 100)
				_, err = monitor.AddVolume("tester", 120)
				if err != nil {
					t.Fatalf("test AddVolum Fail: %v", err)
				}
				wg.Done()
			}()
		}
		wg.Wait()
		user, _ := db.GetUserFromID("tester")
		if user.DataStats.Volume != 320*100 {
			t.Fatalf("teat AddVolum fail, expect 200, got %v", user.DataStats.Volume)
		}
	})
	t.Run("test AddUpload", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func() {
				_, err := monitor.AddUploadTraffic("tester", 100)
				_, err = monitor.AddUploadTraffic("tester", 100)
				_, err = monitor.AddUploadTraffic("tester", 120)
				if err != nil {
					t.Fatalf("test AddUpload Fail: %v", err)
				}
				wg.Done()
			}()
		}
		wg.Wait()
		user, _ := db.GetUserFromID("tester")
		if user.DataStats.UploadTraffic[util.Config.LocalCloudID] != 320*1000 {
			t.Fatalf("teat AddVolum fail, expect 200, got %v", user.DataStats.UploadTraffic[util.Config.LocalCloudID])
		}
	})
}

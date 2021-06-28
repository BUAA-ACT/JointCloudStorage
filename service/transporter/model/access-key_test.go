package model

import "testing"

func TestAccessKeyDB(t *testing.T) {
	var db *AccessKeyDB
	var key AccessKey
	var err error
	t.Run("test New accessKeyDB", func(t *testing.T) {
		dao, err := InitDao()
		if err != nil {
			t.Fatalf("init Dao fail, err: %s", err.Error())
		}
		db, _ = NewAccessKeyDB(dao)
		t.Log(db)
	})
	t.Run("test Gen key for tester", func(t *testing.T) {
		key, err = db.GenerateKeys("tester")
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("test Certificate", func(t *testing.T) {
		uid, err := db.Certificate(key.AccessKey, key.SecretKey)
		if err != nil {
			t.Error(err)
		}
		t.Logf("certificate success: %v", uid)
	})
}

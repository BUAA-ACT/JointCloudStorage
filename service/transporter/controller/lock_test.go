package controller

import (
	"act.buaa.edu.cn/jcspan/transporter/util"
	"testing"
)

func TestLock(t *testing.T) {
	l, err := NewLock(util.Config.ZookeeperHost)
	if err != nil {
		t.Fatalf("create Lock connect error: %v", err)
	}
	l.UnLockAll("/")
	t.Run("Lock", func(t *testing.T) {
		err = l.Lock("/tester/lockTest")
		if err != nil {
			t.Fatalf("Lock error: %v", err)
		}
	})
	t.Run("unlock", func(t *testing.T) {
		err = l.UnLock("/tester/lockTest")
		if err != nil {
			t.Fatalf("unlock error: %v", err)
		}
	})
}

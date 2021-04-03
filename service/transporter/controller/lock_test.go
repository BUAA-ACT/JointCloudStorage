package controller

import (
	"act.buaa.edu.cn/jcspan/transporter/util"
	"testing"
)

func TestLock(t *testing.T) {
	l, err := NewLock(util.CONFIG.ZookeeperHost)
	if err != nil {
		t.Fatalf("create lock connect error: %v", err)
	}
	t.Run("lock", func(t *testing.T) {
		err = l.Lock("/tester/lockTest")
		if err != nil {
			t.Fatalf("lock error: %v", err)
		}
	})
	t.Run("unlock", func(t *testing.T) {
		err = l.UnLock("/tester/lockTest")
		if err != nil {
			t.Fatalf("unlock error: %v", err)
		}
	})
}

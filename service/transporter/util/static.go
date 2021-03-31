package util

import (
	"fmt"
	"time"
)

var (
	Version   = "0.3.3"
	GitRev    = "devRun"
	BuildTime = time.Now().Format("2006-01-02_15:04:05")
)

func GetVersionStr() string {
	return fmt.Sprintf("%s-%s.%s", Version, GitRev, BuildTime)
}

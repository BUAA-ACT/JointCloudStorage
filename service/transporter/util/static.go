package util

import (
	"fmt"
	"time"
)

var (
	Version   string
	GitRev    string
	BuildTime string
)

func GetVersionStr() string {
	if BuildTime == "" {
		BuildTime = time.Now().Format("2006-01-02_15:04:05")
	}
	if GitRev == "" {
		GitRev = "devRun"
	}
	return fmt.Sprintf("%s-%s.%s", Version, GitRev, BuildTime)
}

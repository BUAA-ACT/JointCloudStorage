package sweeper

import (
	"cloud-storage-httpserver/args"
	"cloud-storage-httpserver/dao"
	"time"
)

func CleanVerifyCode() {
	var verifyCodeSweeper = dao.VerifyCodeDao
	for {
		verifyCodeSweeper.CleanVerifyCode()
		time.Sleep(args.VerifyCodeCleanTime)
	}
}

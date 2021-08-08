package sweeper

import (
	"cloud-storage-httpserver/args"
	"cloud-storage-httpserver/dao"
	"time"
)

func CleanAccessToken() {
	var accessTokenSweeper = dao.AccessTokenDao
	var cc int64 = 0
	for {
		cleanResult, _ := accessTokenSweeper.CleanAccessToken()
		cc += cleanResult.DeletedCount
		time.Sleep(args.AccessTokenCleanTime)
	}
}

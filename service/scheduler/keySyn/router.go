package keySyn

import (
	"github.com/gin-gonic/gin"
	"shaoliyin.me/jcspan/dao"
)

var (
	endPointAddKey    = "/add_key"
	endPointDeleteKey = "/delete_key"
)

func KeySyncInit(cid string, r *gin.Engine) {
	localCid = cid

	r.POST(endPointAddKey, PostKeyUpsert)
	r.POST(endPointDeleteKey, PostKeyDelete)
}

func KeySyncDaoInit(map[string]dao.DatabaseConfig) {

	dao.NewDao(mongoURI)
	keyDao = dao.GetDatabaseInstance()
}

func KeySyncRouteInit() {

}

package keySyn

import (
	"github.com/gin-gonic/gin"
	"shaoliyin.me/jcspan/dao"
)

func KeySynInit(cid string, r *gin.Engine) {
	keyDao = dao.GetDatabaseInstance()
	localCid = cid

	r.POST("/add_key", PostKeyUpsert)
	r.POST("/delete_key", PostKeyDelete)
}

package keySyn

import (
	"github.com/gin-gonic/gin"
	"shaoliyin.me/jcspan/dao"
)

func KeySynInit(cid string,dao *dao.Dao,r *gin.Engine){
	keyDao=dao
	localCid=cid

	r.POST("/key_upsert",PostKeyUpsert)
	r.POST("/key_delete",PostKeyDelete)
}
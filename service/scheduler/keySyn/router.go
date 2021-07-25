package keySyn

import (
	"github.com/gin-gonic/gin"
	"shaoliyin.me/jcspan/dao"
)


var (
	endPointAddKey		=	"/add_key"
	endPointDeleteKey	=	"/delete_key"
)
func KeySynInit(cid string, r *gin.Engine) {
	keyDao = dao.GetDatabaseInstance()
	localCid = cid

	r.POST(endPointAddKey, PostKeyUpsert)
	r.POST(endPointDeleteKey, PostKeyDelete)
}

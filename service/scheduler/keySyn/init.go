package keySyn

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"shaoliyin.me/jcspan/dao"
)

const (
	endPointAddKey    = "/add_key"
	endPointDeleteKey = "/delete_key"
)

const (
	CallerHttpServer = "http-server"
	CallerScheduler  = "scheduler"
	CallerHeaderName = "Caller"

	SynTypeUpsert = "upsert"
	SynTypeDelete = "delete"
)

var (
	keyCol   *mongo.Collection
	cloudCol *mongo.Collection
	localCid string
)

func IDInit(cid string) {
	localCid = cid
}

func DaoInit(mongoURI string, databaseMap map[string]map[string]*dao.CollectionConfig) error {
	return dao.NewDao(mongoURI, databaseMap)
}

func RouteInit(r *gin.Engine) {
	r.POST(endPointAddKey, PostKeyUpsert)
	r.POST(endPointDeleteKey, PostKeyDelete)
}

func SetKeyCol(thisCol *mongo.Collection) {
	keyCol = thisCol
}

func SetCloudCol(thisCol *mongo.Collection) {
	cloudCol = thisCol
}

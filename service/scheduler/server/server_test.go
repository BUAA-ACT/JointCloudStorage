package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shaoliyin.me/jcspan/config"
	"shaoliyin.me/jcspan/dao"
	"strings"
	"testing"
)

func serverPlugIn(t *testing.T, r *gin.Engine) {
	IDInit(*config.FlagCloudID)
	databaseMap := map[string]map[string]*dao.CollectionConfig{
		*config.FlagEnv: map[string]*dao.CollectionConfig{
			config.CloudCollectionName:           nil,
			config.UserCollectionName:            nil,
			config.FileCollectionName:            nil,
			config.MigrationAdviceCollectionName: nil,
		},
	}
	err := DaoInit(*config.FlagMongo, databaseMap)
	if err != nil {
		t.Fatal("dao init failed")
	}
	cloudCol = databaseMap[*config.FlagEnv][config.CloudCollectionName].CollectionHandler
	userCol = databaseMap[*config.FlagEnv][config.UserCollectionName].CollectionHandler
	fileCol = databaseMap[*config.FlagEnv][config.FileCollectionName].CollectionHandler
	adviceCol = databaseMap[*config.FlagEnv][config.MigrationAdviceCollectionName].CollectionHandler
	RouteInit(r)
	r.Run(":8082")
}

func TestAllCloudStatus(t *testing.T) {
	t.Log("start testing getallcloudstatus.")
	r := gin.Default()
	serverPlugIn(t, r)

	req, err := http.NewRequest("GET", "http://0.0.0.0:8082/all_clouds_status", nil)
	if err != nil {
		t.Fatal("\tShould be able to create a request.", err)
	}

	//rw := httptest.NewRecorder()
	//http.DefaultServeMux.ServeHTTP(rw, req)
	//fmt.Println(rw)
	rw, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal("\tcan't do http Get request.", err)
	}
	if rw.StatusCode != 200 {
		t.Fatal("\tShould receive 200.", rw.Status)
	}
	t.Log("\t Should receive 200.")

}

func TestPostUpdateClouds(t *testing.T) {
	t.Log("start testing PostUpdateClouds.")
	r := gin.Default()
	go serverPlugIn(t, r)

	data := `{"CloudID":"aliyun-beijing","Endpoint":"oss-cn-beijing.aliyuncs.com",
			"StoragePrice":0.5,"TrafficPrice":0.5,"Availability":0.9999,
			"Status":"UP","Location":"116.381252,20.0","Address":""}`
	req, err := http.NewRequest("POST", "http://0.0.0.0:8082/update_clouds", strings.NewReader(data))
	if err != nil {
		t.Log("\tCan't creat a request.", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal("\tCan't send http request.", err)
	}

	if resp.StatusCode != 200 {
		t.Fatal("\tShould receive 200.", resp.StatusCode, resp)
	}

	t.Log("\tShould receive 200.", resp.StatusCode)
}

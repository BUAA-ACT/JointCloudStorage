package newcloud

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"shaoliyin.me/jcspan/dao"
	"strings"
	"testing"
)

func newCloudPlugIn(t *testing.T, r *gin.Engine, mongoURI string, databaseName string, cid string, envMod string) {
	IDInit(cid, envMod)

	databaseMap := map[string]map[string]*dao.CollectionConfig{
		databaseName: map[string]*dao.CollectionConfig{
			CollectionCloud:     nil,
			CollectionTempCloud: nil,
			CollectionVoteCloud: nil,
		},
	}
	err := DaoInit(mongoURI, databaseMap)
	if err != nil {
		t.Fatalf("dao init failed")
	}
	cloudCol = databaseMap[databaseName][CollectionCloud].CollectionHandler
	tempCloudCol = databaseMap[databaseName][CollectionTempCloud].CollectionHandler
	voteCloudCol = databaseMap[databaseName][CollectionVoteCloud].CollectionHandler
	RouteInit(r)
}

func TestPostNewCloudVote(t *testing.T) {
	t.Log("Start testing endpoint /new_cloud_vote, func:PostNewCloudVote.")

	//初始化服务
	r := gin.Default()
	newCloudPlugIn(t, r, "mongodb://localhost:27017", "test", "test-wanggj", "localDebug")
	data := `{"id":"test_cloud_id","vote_num":1,"address":"aliyun-beijing","cloud":{"CloudID":"test_cloud_id","Endpoint":"asdfasdasd",
			"StoragePrice":0.5,"TrafficPrice":0.5,"Availability":0.9999,
			"Status":"UP","Location":"116.381252,20.0","Address":""}}`
	req, _ := http.NewRequest("POST", "/new_cloud_vote", strings.NewReader(data))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	t.Log(w.Body.String())
}

func TestPostNewCloud(t *testing.T) {
	t.Log("Start testing endpoint /new_cloud, func:PostNewCloud.")

	//初始化服务
	r := gin.Default()
	newCloudPlugIn(t, r, "mongodb://localhost:27017", "dev", "aliyun-beijing", "localDebug")
	data := `{"CloudID":"test-wanggj","Endpoint":"asdfasdasd",
			"StoragePrice":0.5,"TrafficPrice":0.5,"Availability":0.9999,
			"Status":"UP","Location":"116.381252,20.0","Address":""}`

	req, _ := http.NewRequest("POST", "/new_cloud", strings.NewReader(data))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	t.Log(w.Body.String())
}

func TestGetVoteRequest(t *testing.T) {
	t.Log("Start testing endpoint /new_cloud, func:PostNewCloud.")

	//初始化服务
	r := gin.Default()
	newCloudPlugIn(t, r, "mongodb://localhost:27017", "dev", "test-wanggj", "localDebug")

	req, _ := http.NewRequest("GET", "/vote_request", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	t.Log(w.Body.String())
}

func TestPostCloudVote(t *testing.T) {
	t.Log("Start testing endpoint /new_cloud, func:PostNewCloud.")

	//初始化服务
	r := gin.Default()
	newCloudPlugIn(t, r, "mongodb://localhost:27017", "dev", "aliyun-beijing", "localDebug")
	data := `{"CloudID":"test-wanggj"}`
	data2 := `{"CloudID":"test_cloud_id"}`
	//testCloud:=`{"id":"test_cloud_id","vote_num":1,"address":"aliyun-beijing","cloud":{"CloudID":"test_cloud_id","Endpoint":"asdfasdasd",
	//		"StoragePrice":0.5,"TrafficPrice":0.5,"Availability":0.9999,
	//		"Status":"UP","Location":"116.381252,20.0","Address":""}}`

	req, _ := http.NewRequest("POST", "/cloud_vote", strings.NewReader(data))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	t.Log(w.Body.String())

	req, _ = http.NewRequest("POST", "/cloud_vote", strings.NewReader(data2))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	t.Log(w.Body.String())
}

func TestPostMasterCloudVote(t *testing.T) {
	t.Log("Start testing endpoint /new_cloud, func:PostNewCloud.")

	//初始化服务
	r := gin.Default()
	newCloudPlugIn(t, r, "mongodb://localhost:27017", "dev", "aliyun-beijing", "localDebug")
	data := `{"id":"test-wanggj","vote":1}`

	req, _ := http.NewRequest("POST", "/master_cloud_vote", strings.NewReader(data))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	t.Log(w.Body.String())
}

func TestCloudSyn(t *testing.T) {
	t.Log("Start testing endpoint /new_cloud, func:PostNewCloud.")

	//初始化服务
	r := gin.Default()
	newCloudPlugIn(t, r, "mongodb://localhost:27017", "dev", "aliyun-beijing", "localDebug")
	data := `[{"CloudID":"aliyun-beijing","Endpoint":"asdfasdasd",
			"StoragePrice":0.5,"TrafficPrice":0.5,"Availability":0.9999,
			"Status":"UP","Location":"116.381252,20.0","Address":""}]`
	data2 := `[{"CloudID":"test-wanggj","Endpoint":"asdfasdasd",
			"StoragePrice":0.5,"TrafficPrice":0.5,"Availability":0.9999,
			"Status":"UP","Location":"116.381252,20.0","Address":""}]`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/cloud_syn", strings.NewReader(data))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	t.Log(w.Body.String())

	req2, _ := http.NewRequest("POST", "/cloud_syn", strings.NewReader(data2))
	r.ServeHTTP(w, req2)
	assert.Equal(t, 200, w.Code)
	t.Log(w.Body.String())
	cloud, err := dao.GetCloud(cloudCol, "test-wanggj")
	if err != nil {
		t.Fatal(err)
	}
	if cloud.CloudID != "test-wanggj" {
		t.Fatal("Cloud syn fail!")
	}
}

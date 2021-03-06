package keySyn

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	r *gin.Engine
)

func InitTest() error {
	r = gin.Default()
	KeySynInit("aliyun-beijing", r)
	return nil
}

//type AccessKey struct {
//	UserID     string    `json:"UserID" bson:"user_id"`
//	AccessKey  string    `json:"AccessKey" bson:"access_key"`
//	SecretKey  string    `json:"SecretKey" bson:"secret_key"`
//	Comment    string    `json:"Comment" bson:"comment"`
//	CreateTime time.Time `json:"CreateTime" bson:"create_time"`
//	Available  bool      `json:"Available" bson:"available"`
//}

func TestPostKeyUpsert(t *testing.T) {
	data1 := `{"user_id":"wanggj","access_key":"wanggj_ak1","secret_key":"wanggj_sk1","comment":"this is test data1"}`
	data2 := `{"user_id":"wanggj2","access_key":"wanggj_ak2","secret_key":"wanggj_sk2","comment":"this is test data2"}`
	t.Log("Start testing endpoint PostKeyUpsert.")
	if err := InitTest(); err != nil {
		t.Error("Init test failed! ", err.Error())
	}

	req1, _ := http.NewRequest("POST", "/add_key", strings.NewReader(data1))
	req2, _ := http.NewRequest("POST", "/add_key", strings.NewReader(data2))
	req1.Header.Set(CallerHeaderName, CallerHttpServer)
	req2.Header.Set(CallerHeaderName, CallerHttpServer)

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req1)
	if resp.Code != 200 {
		t.Error("Status should be 200!", resp.Code, resp.Body)
	}
	r.ServeHTTP(resp, req2)
	if resp.Code != 200 {
		t.Error("Status should be 200!", resp.Code)
	}

	t.Log("PostKeyUpsert succeed!")
}

func TestPostKeyDelete(t *testing.T) {
	data1 := `{"user_id":"wanggj","access_key":"wanggj_ak1","secret_key":"wanggj_sk1","comment":"this is test data1"}`
	data2 := `{"user_id":"wanggj2","access_key":"wanggj_ak2","secret_key":"wanggj_sk2","comment":"this is test data2"}`
	t.Log("Start testing endpoint PostKeyDelete.")
	if err := InitTest(); err != nil {
		t.Error("Init test failed! ", err.Error())
	}

	req1, _ := http.NewRequest("POST", "/delete_key", strings.NewReader(data1))
	req2, _ := http.NewRequest("POST", "/delete_key", strings.NewReader(data2))
	req1.Header.Set(CallerHeaderName, CallerHttpServer)
	req2.Header.Set(CallerHeaderName, CallerHttpServer)

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req1)
	if resp.Code != 200 {
		t.Error("Status should be 200!", resp.Code, resp.Body)
	}
	r.ServeHTTP(resp, req2)
	if resp.Code != 200 {
		t.Error("Status should be 200!", resp.Code)
	}

	t.Log("PostKeyDelete succeed!")
}

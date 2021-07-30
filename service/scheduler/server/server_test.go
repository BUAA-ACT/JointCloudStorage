package server

import (
	"net/http"
	"strings"
	"testing"

	_ "shaoliyin.me/jcspan/testinit"
)

func init() {

}

func TestAllCloudStatus(t *testing.T) {
	t.Log("start testing getallcloudstatus.")
	//r := gin.Default()
	//NewRouter(r)
	//r.Run()

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

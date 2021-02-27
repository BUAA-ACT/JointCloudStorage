package transporter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestNewRouter(t *testing.T) {
	// 初始化任务数据库
	storage := NewInMemoryTaskStorage()
	processor := TaskProcessor{}
	processor.SetTaskStorage(storage)
	// 初始化存储数据库
	processor.SetStorageDatabase(NewSimpleInMemoryStorageDatabase())
	// 初始化路由
	router := NewTestRouter(processor)
	// 启动 processor
	processor.StartProcessTasks(context.Background())

	req, _ := http.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	t.Logf("Recive: %v", recorder.Body.String())

	req, _ = http.NewRequest("GET", "/test/setcookie", nil)
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	t.Logf("Recive: %v", recorder.Body.String())
	cookies := recorder.Result().Cookies()
	if len(cookies) == 0 {
		t.Error("No Cookie Set")
	}
	t.Logf("Cookie: %v", cookies[0])

	t.Run("simple upload", func(t *testing.T) {
		data := url.Values{"fileName": {"test.txt"}}
		body := strings.NewReader(data.Encode())
		req, _ = http.NewRequest("POST", "/upload/jcspan/path/to/file", body)
		req.AddCookie(cookies[0])
		recorder = httptest.NewRecorder()
		router.ServeHTTP(recorder, req)
		time.Sleep(time.Second * 5)
	})
}

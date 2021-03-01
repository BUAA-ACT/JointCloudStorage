package controller

import (
	"act.buaa.edu.cn/jcspan/transporter/model"
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestNewRouter(t *testing.T) {
	// 初始化任务数据库
	storage := model.NewInMemoryTaskStorage()
	processor := TaskProcessor{}
	processor.SetTaskStorage(storage)
	// 初始化存储数据库
	processor.SetStorageDatabase(model.NewSimpleInMemoryStorageDatabase())
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
		filename := "../test/tmp/test.txt"
		f, err := os.Open(filename)
		if err != nil {
			t.Error("Open test file Fail")
		}
		defer f.Close()

		//req,_ http.Post("/upload/jcspan/path/to/file", "multipart/form-data", body)
		//	req, _ = http.NewRequest("POST", "/upload/jcspan/path/to/file", body)
		req, _ = postFile("test.txt", "../test/tmp/test.txt", "/upload/path/to/jcspantest.txt")
		req.AddCookie(cookies[0])
		recorder = httptest.NewRecorder()
		router.ServeHTTP(recorder, req)
		time.Sleep(time.Second * 5)
	})
}

func postFile(filename string, filepath string, target_url string) (*http.Request, error) {
	body_buf := bytes.NewBufferString("")
	body_writer := multipart.NewWriter(body_buf)

	// use the body_writer to write the Part headers to the buffer
	_, err := body_writer.CreateFormFile("file", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return nil, err
	}

	// the file data will be the second part of the body
	fh, err := os.Open(filepath)
	if err != nil {
		fmt.Println("error opening file")
		return nil, err
	}
	// need to know the boundary to properly close the part myself.
	boundary := body_writer.Boundary()
	//close_string := fmt.Sprintf("\r\n--%s--\r\n", boundary)
	close_buf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))

	// use multi-reader to defer the reading of the file data until
	// writing to the socket buffer.
	request_reader := io.MultiReader(body_buf, fh, close_buf)
	fi, err := fh.Stat()
	if err != nil {
		fmt.Printf("Error Stating file: %s", filename)
		return nil, err
	}
	req, err := http.NewRequest("POST", target_url, request_reader)
	if err != nil {
		return nil, err
	}

	// Set headers for multipart, and Content Length
	req.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	req.ContentLength = fi.Size() + int64(body_buf.Len()) + int64(close_buf.Len())
	return req, nil
}

func multipartUpload(f io.Reader, fields map[string]string) (*bytes.Buffer, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("file", fields["filename"])
	if err != nil {
		return nil, fmt.Errorf("CreateFormFile %v", err)
	}

	_, err = io.Copy(fw, f)
	if err != nil {
		return nil, fmt.Errorf("copying fileWriter %v", err)
	}

	for k, v := range fields {
		_ = writer.WriteField(k, v)
	}

	err = writer.Close() // close writer before POST request
	if err != nil {
		return nil, fmt.Errorf("writerClose: %v", err)
	}

	return body, nil
}

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
		req, _ = postFile("test.txt", "../test/tmp/test.txt", "/upload/path/to/jcspantest.txt", "1")
		req.AddCookie(cookies[0])
		recorder = httptest.NewRecorder()
		router.ServeHTTP(recorder, req)
		fmt.Print(req.Header)
		fmt.Print(req.Body)
		time.Sleep(time.Second * 5)
	})

	t.Run("simple download", func(t *testing.T) {
		req, _ = http.NewRequest("GET", "/jcspan/path/to/jcspantest.txt", nil)
		req.AddCookie(cookies[0])
		recorder = httptest.NewRecorder()
		router.ServeHTTP(recorder, req)
		if recorder.Code != http.StatusOK {
			t.Errorf("except code %v, get %v", http.StatusOK, recorder.Code)
		}
		t.Logf("download url: %v", recorder.Body)
	})

	t.Run("simple download without auth", func(t *testing.T) {
		req, _ = http.NewRequest("GET", "/jcspan/path/to/jcspantest.txt", nil)
		recorder = httptest.NewRecorder()
		router.ServeHTTP(recorder, req)
		if recorder.Code != http.StatusUnauthorized {
			t.Errorf("except code %v, get %v", http.StatusOK, recorder.Code)
		}
		t.Logf("UnAuth download url: %v", recorder.Body)
	})

	t.Run("index path", func(t *testing.T) {
		req, _ = http.NewRequest("GET", "/index/path/to/", nil)
		req.AddCookie(cookies[0])
		recorder = httptest.NewRecorder()
		router.ServeHTTP(recorder, req)
		t.Logf("index: %v", recorder.Body)
	})

	t.Run("simple sync", func(t *testing.T) {
		body := new(bytes.Buffer)
		w := multipart.NewWriter(body)
		contentType := w.FormDataContentType()

		w.WriteField("srcpath", "path/to/")
		w.WriteField("dstpath", "dst/to/")
		w.WriteField("sid", "tttteeeesssstttt")
		w.Close()
		req, _ = http.NewRequest("POST", "/task/simplesync", body)
		req.Header.Set("Content-Type", contentType)

		recorder = httptest.NewRecorder()
		router.ServeHTTP(recorder, req)
		if recorder.Code != http.StatusOK {
			t.Errorf("sync fail")
		}
		t.Logf("%v", recorder.Body)
		time.Sleep(time.Second * 5)
	})

	t.Run("Create Upload Task and process", func(t *testing.T) {
		jsonStr := []byte(`{
  "TaskType": "Upload",
   "Uid": "12",
   "Sid": "tttteeeesssstttt",
   "DestinationPath":"/path/to/upload/",
   "StoragePlan":{
      "StorageMode": "Replica",
      "Clouds": [
         {
            "ID": "txyun-chongqing"
         },
         {
            "ID": "aliyun-beijing"
         }
      ]
   }
}`)
		req, _ = http.NewRequest("POST", "/task", bytes.NewBuffer(jsonStr))
		recorder = httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		tid := recorder.Body.String()

		filename := "../test/tmp/test.txt"
		f, err := os.Open(filename)
		if err != nil {
			t.Error("Open test file Fail")
		}
		defer f.Close()

		//req,_ http.Post("/upload/jcspan/path/to/file", "multipart/form-data", body)
		//	req, _ = http.NewRequest("POST", "/upload/jcspan/path/to/file", body)
		req, _ = postFile("test.txt", "../test/tmp/test.txt", "/upload/path/to/jcspantest.txt", tid)
		req.AddCookie(cookies[0])
		recorder = httptest.NewRecorder()
		router.ServeHTTP(recorder, req)
		fmt.Print(req.Header)
		fmt.Print(req.Body)
		time.Sleep(time.Second * 50)
	})
}

func postFile(filename string, filepath string, target_url string, tid string) (*http.Request, error) {
	body_buf := bytes.NewBufferString("")
	body_writer := multipart.NewWriter(body_buf)

	// use the body_writer to write the Part headers to the buffer
	writer, _ := body_writer.CreateFormField("tid")
	writer.Write([]byte(tid))
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

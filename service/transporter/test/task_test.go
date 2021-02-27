package transporter

import (
	"act.buaa.edu.cn/jcspan/transporter"
	"io"
	"os"
	"testing"
	"time"
)

func TestTask(t *testing.T) {
	storage := transporter.NewInMemoryTaskStorage()
	processor := transporter.TaskProcessor{}
	processor.SetTaskStorage(storage)
	processor.SetStorageDatabase(transporter.NewSimpleInMemoryStorageDatabase())
	t.Run("add user upload task", func(t *testing.T) {
		fileName := time.Now().Format("2006-01-02-15-04-05-UserUploadTest.txt")
		os.MkdirAll("./tmp", os.ModePerm)
		filePath := "./tmp/" + fileName
		f, err := os.Create(filePath)
		defer f.Close()
		if err != nil {
			t.Errorf("local File creat Fail")
		} else {
			io.WriteString(f, fileName)
		}
		processor.CreateTask(
			transporter.USER_UPLOAD_SIMPLE,
			"1",
			filePath,
			fileName,
		)
		processor.ProcessTasks()
	})
}

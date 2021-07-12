package transporter

import (
	"act.buaa.edu.cn/jcspan/transporter/controller"
	"act.buaa.edu.cn/jcspan/transporter/model"
	"io"
	"os"
	"testing"
	"time"
)

func TestTask(t *testing.T) {
	// 初始化任务数据库
	storage := model.NewInMemoryTaskStorage()
	processor := controller.TaskProcessor{}
	processor.SetTaskStorage(storage)
	// 初始化存储数据库
	processor.SetStorageDatabase(model.NewSimpleInMemoryStorageDatabase())
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
			model.UPLOAD,
			"1",
			filePath,
			fileName,
		)
		processor.ProcessTasks()
	})
}

package transporter

import (
	"act.buaa.edu.cn/jcspan/transporter"
	"testing"
)

func TestTask(t *testing.T) {
	storage := transporter.NewInMemoryTaskStorage()
	processor := transporter.TaskProcessor{}
	processor.SetTaskStorage(storage)
	t.Run("add user upload task", func(t *testing.T) {
		processor.CreateTask(transporter.USER_UPLOAD_SIMPLE, 1, "/tmp/test.txt", "/jcsPan/test.txt")
		processor.ProcessTasks()
	})
}

package newcloud

import (
	"encoding/json"

	"shaoliyin.me/jcspan/dao"
)

//CloudReader implement the interface io.Reader
//return json format of Cloud struct
//used to send http request
type CloudReader struct {
	json dao.Cloud
	i    int64
}

func (cr *CloudReader) Read(b []byte) (int, error) {
	b, err := json.Marshal(cr)
	if err != nil {
		return len(b), err
	}
	return len(b), nil
}

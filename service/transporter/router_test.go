package transporter

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRouter(t *testing.T) {
	router := NewRouter()

	req, _ := http.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	t.Logf("Recive: %v", recorder.Body.String())

}

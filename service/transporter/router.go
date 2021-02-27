package transporter

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func NewRouter() *httprouter.Router {
	router := httprouter.New()
	router.GET("/", Index)
	return router
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "JcsPan Transporter")
}

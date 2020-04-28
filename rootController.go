package myhouse

import (
	"net/http"

	"github.com/gorilla/mux"
)

//NewRouter produces the router with configured handlers
func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("./static"))))
}

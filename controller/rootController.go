package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func router(r *mux.Router) {
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("./static"))))
}

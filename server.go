//Package myhouse defines the web server that  runs household uitlites
package myhouse

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func assertString(i interface{}) string {
	s, ok := i.(string)
	if ok {
		return s
	}
	return ""
}

// Server provides the http handler for the webserver which routes and muxes
type Server struct {
	ListenAddress string
	cfg           *myhouseConfig
}

// ListenAndServe returns  when the server errors-out but blocks while the server runs
func (s *Server) ListenAndServe() error {
	r := mux.NewRouter()
	fmt.Printf("Starting up on %s]n", s.ListenAddress)
	return http.ListenAndServe(s.ListenAddress, r)
}

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello world!")
}

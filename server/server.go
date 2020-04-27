//Package myhouse defines the web server that  runs household uitlites
package myhouse

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Server provides the http handler for the webserver which routes and muxes
type Server struct {
	ListenAddress string
}

func assertString(i interface{}) string {
	s, ok := i.(string)
	if ok {
		return s
	}
	return ""
}

//NewServer inits the server from the default config or the passed-in Reader
func NewServer(addr string) (*Server, error) {
	s := Server{ListenAddress: addr}
	return &s, nil
}

func (s *Server) ListenAndServe() error {
	r := mux.NewRouter()
	r.HandleFunc("/", Hello)
	http.Handle("/", r)
	fmt.Printf("Starting up on %s]n", s.ListenAddress)

	return http.ListenAndServe(s.ListenAddress, r)
}

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello world!")
}

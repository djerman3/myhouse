//Package myhouse defines the web server that  runs household uitlites
package myhouse

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Server provides the http handler for the webserver which routes and muxes
type Server struct {
	ListenAddress string
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
}

// ListenAndServe returns  when the server errors-out but blocks while the server runs
func (s *Server) ListenAndServe() error {
	r := MyNewRouter()
	log.Printf("Starting up on %s\n", s.ListenAddress)
	srv := &http.Server{
		Addr:           s.ListenAddress,
		Handler:        r,
		ReadTimeout:    s.ReadTimeout,
		WriteTimeout:   s.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	return srv.ListenAndServe()
}

// NewServer loads the config and makes a new web server
func NewServer(cfgFileName *string) (*Server, error) {
	cfg, err := GetConfig(cfgFileName)
	if err != nil {
		return nil, fmt.Errorf("Server Init: %v", err)
	}
	s := Server{
		ListenAddress: cfg.Server.ListenAddress,
		ReadTimeout:   5 * time.Second,
		WriteTimeout:  5 * time.Second,
	}
	return &s, nil
}

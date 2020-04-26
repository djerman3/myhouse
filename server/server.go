//Package control defines the control server that  runs household uitlites
package control

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// Server provides the http handler for the webserver which routes and muxes
type Server struct {
	ListenAddress string `json:"listenAddress,omitempty"`
	ConfigFile    string `json:"configFile,omitempty"`
}

func assertString(i interface{}) string {
	s, ok := i.(string)
	if ok {
		return s
	}
	return ""
}

//NewServer inits the server from the default config or the passed-in Reader
func NewServer(cfgFileName string) (*Server, error) {
	cfgFile, err := os.Open(cfgFileName)
	if err != nil {
		return nil, fmt.Errorf("Can't open config file %s:%v", cfgFileName, err)
	}

	cfgJSON, err := ioutil.ReadAll(cfgFile)
	if err != nil {
		return nil, fmt.Errorf("Server Init: %v", err)
	}
	result := make(map[string]interface{})
	err = json.Unmarshal(cfgJSON, &result)
	if err != nil {
		return nil, fmt.Errorf("Server Init: %v", err)
	}
	cfg, ok := result["server"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("This is not a valid server cfg:%#v", result["server"])
	}
	s := Server{}
	s.ListenAddress = assertString(cfg["listenAdddress"]) + ":" + assertString(cfg["listenPort"])
	s.ConfigFile = cfgFileName
	return &s, err

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

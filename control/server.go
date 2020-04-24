//Package control defines the control server that  runs household uitlites
package control

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

// Server provides the http handler for the webserver which routes and muxes
type Server struct {
	Port          int    `json:"port,omitempty"`
	ListenAddress string `json:"listenAddress,omitempty"`
	CertFile      string `json:"certFile,omitempty"`
	KeyFile       string `json:"keyFile,omitempty"`
	ConfigFile    string `json:"configFile,omitempty"`
}

//NewServer inits the server from the default config or the passed-in Reader
func NewServer(c io.Reader) (*Server, error) {
	cfgFile := c
	const cfgDir = "/etc/homecontrol"
	s := Server{}
	if cfgFile == nil {
		r := bytes.NewReader([]byte(
			`{
				"port":"44800",
				"listenAddress":"0.0.0.0",
				"certFile":"` + cfgDir + `/cert.dat",
				"keyFile":"` + cfgDir + `/key.dat",
				"configFile":"` + cfgDir + `/homecontrol.conf",
			}`,
		))
		cfgFile = r
	}
	cfg, err := ioutil.ReadAll(cfgFile)
	if err != nil {
		return nil, fmt.Errorf("Server Init: %v", err)
	}
	err = json.Unmarshal(cfg, &s)
	if err != nil {
		return nil, fmt.Errorf("Server Init: %v", err)
	}
	return &s, err

}

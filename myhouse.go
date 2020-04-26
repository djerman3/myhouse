package myhouse

import "github.com/djerman3/myhouse/server"

func NewServer(cfgFileName string) (*server.Server, err) {
	return server.NewServer(cfgFileName)
}

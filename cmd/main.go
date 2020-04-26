// Package main runs the executable
package main

import (
	"log"
	"net/http"

	control "github.com/djerman3/homecontrol"
)

func main() {
	cfgFileName := "/home/djerman/projects/homecontrol/etc/homecontrol.json"

	s, err := control.NewHomeServer(cfgFileName)
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(http.ListenAndServe(s.ListenAddress, s))
}

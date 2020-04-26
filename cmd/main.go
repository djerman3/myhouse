// Package main runs the executable
package main

import (
	"log"

	"github.com/djerman3/homecontrol"
)

func main() {
	cfgFileName := "/home/djerman/projects/homecontrol/etc/homecontrol.json"

	s, err := homecontrol.NewServer(cfgFileName)
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(s.ListenAndServe())
}

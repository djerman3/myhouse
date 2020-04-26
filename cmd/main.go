// Package main runs the executable
package main

import (
	"log"

	control "github.com/djerman3/homecontrol"
)

func main() {
	cfgFileName := "/home/djerman/projects/homecontrol/etc/homecontrol.json"

	s, err := control.NewServer(cfgFileName)
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(s.ListenAndServe())
}

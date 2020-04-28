// Package main runs the executable
package main

import (
	"log"

	"github.com/djerman3/myhouse"
)

func main() {
	cfgFileName := "/home/djerman/projects/myhouse/etc/myhouse.json"

	s, err := myhouse.NewServer(&cfgFileName)
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(s.ListenAndServe())
}

// Package main runs the executable
package main

import (
	"log"
)

func main() {
	cfgFileName := "/home/djerman/projects/myhouse/etc/myhouse.json"

	s, err := myhouse.NewServer()
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(s.ListenAndServe())
}

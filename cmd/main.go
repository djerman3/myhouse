// Package main runs the executable
package main

import (
	"log"
	"os"
	"net/http"
)

func main() {
	cfgFileName := "~djerman/projects/control/etc/homecontrol/homecontrol.cfg"
	cfgFile, err := os.Open(cfgFileName)
	if err != nil {
		log.Printf("Opening Config File: %v")
	}
	s := control.NewServer(cfgFile) 
}

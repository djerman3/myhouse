package main

import (
	"log"

	"github.com/djerman3/homecontrol/client"
)

var cfgFilePath = "/home/djerman/projects/homecontrol/etc/homecontrol.json"

func main() {
	routerClient, err := client.NewClient(cfgFilePath)
	if err != nil {
		log.Fatalf("Bad client:%v\n", err)
	}
	err = routerClient.Auth()
	if err != nil {
		log.Fatalf("Client Auth() failed:%v\n", err)
	}
}

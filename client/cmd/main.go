package main

import "github.com/djerman3/homecontrol/client"

var cfgFilePath := "/home/djerman/projects/homecontrol/etc/homecontrol.json"
func main() {
	router_client,err := client.NewClient(cfgFilePath)
	if err != nil{
		log.Fatalf("Bad client:%v\n",err)
	}
}

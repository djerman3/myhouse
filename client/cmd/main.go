package main

import "github.com/djerman3/homecontrol/client"

func main() {
	router_client,err := client.NewClient("[2606:a000:4c84:82f0::1]")
	if err != nil{
		log.Fatalf("Bad client:%v\n",err)
	}
}

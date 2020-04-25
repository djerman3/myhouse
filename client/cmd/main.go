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
	rules, err := routerClient.GetFirewallRules()
	if err != nil {
		log.Printf("Traffic Rule Get failed:%v\n", err)
	}
	//	log.Printf("%#v]n", rules)
	err = routerClient.EnableFirewallRule(rules["ferdinand"].DotName)
	if err != nil {
		log.Printf("Traffic Rule Enable failed:%v\n", err)
	}
	err = routerClient.EnableFirewallRule(rules["ferdinand"].DotName)
	if err != nil {
		log.Printf("Traffic Rule Enable failed:%v\n", err)
	}
	err = routerClient.DisableFirewallRule(rules["ferdinand"].DotName)
	if err != nil {
		log.Printf("Traffic Rule Enable failed:%v\n", err)
	}

}

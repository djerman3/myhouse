package myhouse

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type myhouseConfig struct {
	AwsClient struct {
		Key    string `json:"key"`
		Secret string `json:"secret"`
	} `json:"awsClient"`
	Client struct {
		Password string `json:"password"`
		RPCURL   string `json:"rpcURL"`
		User     string `json:"user"`
	} `json:"client"`
	Server struct {
		ListenAddress string `json:"listenAddress"`
		Password      string `json:"password"`
		Username      string `json:"username"`
	} `json:"server`
}

func NewServer(cfgFileName string) (*server.Server, err) {
	cfgFile, err := os.Open(cfgFileName)
	if err != nil {
		return nil, fmt.Errorf("Can't open config file %s:%v", cfgFileName, err)
	}

	cfgJSON, err := ioutil.ReadAll(cfgFile)
	if err != nil {
		return nil, fmt.Errorf("Server Init: %v", err)
	}
	result := myhouseConfig{}
	err = json.Unmarshal(cfgJSON, &result)
	if err != nil {
		return nil, fmt.Errorf("Server Init: %v", err)
	}
	return server.NewServer(cfgFileName)
}

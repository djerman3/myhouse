package myhouse

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// secret global
var cfgcache *MyConfig

//MyConfig should not be exported, there's a TODO
type MyConfig struct {
	fileName  string //`json:"-"`
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
	} `json:"server"`
}

// GetConfig Load a config so different test entrypoints can be consistent
func GetConfig(cfgFileName *string) (*MyConfig, error) {
	fileName := "/etc/myhouse.json"
	if cfgFileName != nil && *cfgFileName != fileName {
		fileName = *cfgFileName
		cfgcache = nil // new file forces reload
	}
	if cfgcache == nil {
		cfgFile, err := os.Open(fileName)
		if err != nil {
			return nil, fmt.Errorf("Can't open config file %s:%v", fileName, err)
		}
		defer cfgFile.Close()
		cfgJSON, err := ioutil.ReadAll(cfgFile)
		if err != nil {
			return nil, fmt.Errorf("Server Init: %v", err)
		}
		result := MyConfig{}
		err = json.Unmarshal(cfgJSON, &result)
		if err != nil {
			return &result, err
		}
		cfgcache = &result

	}
	return cfgcache, nil

}

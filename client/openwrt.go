// Package client sends commands to an openwrt router for uci managment of router stuff
package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

//LuciRPCMethodRequest allows serialization of the struct to a valid request body
type LuciRPCMethodRequest struct {
	ID     int      `json:"id,omitempty"`
	Method string   `json:"method"`
	Params []string `json:"params"`
}

//LuciRPCSingleResponse allows deserialization of the result body to a valid struct
// where we expect a simple name:string
type LuciRPCSingleResponse struct {
	ID     int    `json:"id,omitempty"`
	Result string `json:"result"`
	Error  string `json:"error,omitempty"`
}

//LuciRPCBoolResponse allows deserialization of the result body to a valid struct
// where we expect a simple name:string
type LuciRPCBoolResponse struct {
	ID     int    `json:"id,omitempty"`
	Result bool   `json:"result"`
	Error  string `json:"error,omitempty"`
}

//LuciRPCComplexResponse allows deserialization of the result body to a valid struct
// where we expect a simple name:string
type LuciRPCComplexResponse struct {
	ID     int                     `json:"id,omitempty"`
	Result *map[string]interface{} `json:"result"`
	Error  string                  `json:"error,omitempty"`
}

//Client remembers the client connection and auth
type Client struct {
	cfgFileName string
	BaseURL     string
	AuthToken   string
	AuthRequest *LuciRPCMethodRequest // we're going to need this one  alot
}

//NewClient returns a new client
func NewClient(cfn string) (*Client, error) {
	cfgFileName := cfn
	cfgFile, err := os.Open(cfgFileName)
	if err != nil {
		//quietly try a fallback:
		//TODO: signal this with an error type for realerr
		realerr := err
		cfgFileName := "/etc/homecontrol/homecontrol.json"
		cfgFile, err = os.Open(cfgFileName)
		if err != nil {
			return nil, fmt.Errorf("Failed to open client config:%v", realerr)
		}
	}
	defer cfgFile.Close()
	jsonConfig, err := ioutil.ReadAll(cfgFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to read client config:%v", err)
	}
	var result map[string]interface{}
	err = json.Unmarshal([]byte(jsonConfig), &result)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse client config:%v", err)
	}
	result = result["client"].(map[string]interface{})

	c := &Client{
		BaseURL: result["rpcURL"].(string),
		AuthRequest: &LuciRPCMethodRequest{
			Method: "login",
			Params: []string{result["user"].(string), result["password"].(string)},
		},
	}
	return c, nil
}

// doPost is a sink to capture http.Client creation on posts
func doPost(URL string, body []byte) (*http.Response, error) {

	// workaround: accept self-signed TLS
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	netClient := &http.Client{Transport: tr,
		Timeout: time.Second * 10,
	}
	return netClient.Post(URL, "application/json", bytes.NewBuffer(body))
}

//Auth attempts client authentication and remembers the auth token in the client struct for the next call(s)
func (c *Client) Auth() error {
	// build a valid request for authentication
	//c.AuthRequest.Params[1] = "foo" // test auth error
	body, err := json.Marshal(c.AuthRequest)
	if err != nil {
		return fmt.Errorf("Error marshalling auth request:%v", err)
	}
	URL := c.BaseURL + "/auth"
	// get answer
	response, err := doPost(URL, body)
	if err != nil {
		return fmt.Errorf("Failed to Authenticate at Router:%v", err)
	}
	defer response.Body.Close()

	// find token
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("Failed to read router Auth response:%v", err)
	}
	fmt.Printf("body;\n%v", string(responseBody))
	result := LuciRPCSingleResponse{} //map[string]interface{}
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return fmt.Errorf("Failed to parse router Auth response:%v", err)
	}
	if len(result.Result) < 4 || len(result.Error) != 0 {
		return fmt.Errorf("Failed to get auth token")
	}
	fmt.Printf("CLient token obtained\n")
	c.AuthToken = result.Result

	return nil
}

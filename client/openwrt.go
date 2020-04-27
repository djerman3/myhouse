// Package client sends commands to an openwrt router for uci managment of router stuff
package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

//LuciFirewallRuleResult is a cut down model for firewall rule activation
type LuciFirewallRuleResult struct {
	DotName string `json:".name"`
	Name    string `json:"name,omitempty"`
	Enabled string `json:"enabled,omitempty"`
}

//LuciRPCFirewallGetAllResponse allows deserialization of the result body to a valid struct
// where we expect a list of named rules
type LuciRPCFirewallGetAllResponse struct {
	ID     int                       `json:"id,omitempty"`
	Result *LuciFirewallEntryJSONMap `json:"result"`
	Error  string                    `json:"error,omitempty"`
}

// LuciFirewallEntryMap defines the map type as it's used several times, and we've been jerking around the firewall enrry type name
type LuciFirewallEntryMap map[string]LuciRPCFirewallEntry
type LuciFirewallEntryJSONMap map[string]LuciRPCFirewallEntryJSON

//Client remembers the client connection and auth
type Client struct {
	cfgFileName string
	BaseURL     string
	AuthToken   string
	AuthRequest *LuciRPCMethodRequest // we're going to need this one  alot
}

// LuciRPCFirewallEntry my attempt to capture firewall entries: fully populated it's invalid but thoughtfully filled out it generates a variety of correct rule types
type LuciRPCFirewallEntry struct {
	Anonymous   bool
	Index       int64
	DotName     string
	Type        string
	Dest        string
	DestIP      []string
	DestPort    []string
	DropInvalid string
	Enabled     string
	Family      string
	Forward     string
	IcmpType    []string
	Input       string
	Limit       string
	Masq        string
	MtuFix      string
	Name        string
	Network     []string
	Output      string
	Path        string
	Proto       []string
	Src         string
	SrcDport    string
	SrcIP       []string
	SrcMac      []string
	SrcPort     string
	SynFlood    string
	Target      string
}

func assertStringArray(s interface{}) []string {
	// s is a nil interface{}, a string, or a []string
	// convert to []string or nil
	if s != nil {

		if st, ok := s.(string); ok {
			return []string{st}
		}

		if sa, ok := s.([]string); ok {
			return sa
		}
		if sa, ok := s.([]interface{}); ok {
			var ss []string
			for i := range sa {
				if st, ok := (sa[i]).(string); ok {
					ss = append(ss, st)
				}

			}
			return ss
		}
	}
	return nil

}
func downList(sa []string) interface{} {
	if sa == nil {
		return nil
	}
	if len(sa) == 1 {
		return interface{}(sa[0])
	}
	return interface{}(sa)
}

// Marshall does special handling of stringarrays that must be either single strings, nil, or arrays of 2 or mormem entries
func (lf *LuciRPCFirewallEntry) preMarshall() *LuciRPCFirewallEntryJSON {
	lfej := LuciRPCFirewallEntryJSON{}
	lfej.Anonymous = lf.Anonymous
	lfej.Index = lf.Index
	lfej.DotName = lf.DotName
	lfej.Type = lf.Type
	lfej.Dest = lf.Dest
	lfej.DestIP = downList(lf.DestIP)
	lfej.DestPort = downList(lf.DestPort)
	lfej.DropInvalid = lf.DropInvalid
	lfej.Enabled = lf.Enabled
	lfej.Family = lf.Family
	lfej.Forward = lf.Forward
	lfej.IcmpType = downList(lf.IcmpType)
	lfej.Input = lf.Input
	lfej.Limit = lf.Limit
	lfej.Masq = lf.Masq
	lfej.MtuFix = lf.MtuFix
	lfej.Name = lf.Name
	lfej.Network = downList(lf.Network)
	lfej.Output = lf.Output
	lfej.Path = lf.Path
	lfej.Proto = downList(lf.Proto)
	lfej.Src = lf.Src
	lfej.SrcDport = lf.SrcDport
	lfej.SrcIP = downList(lf.SrcIP)
	lfej.SrcMac = downList(lf.SrcMac)
	lfej.SrcPort = lf.SrcPort
	lfej.SynFlood = lf.SynFlood
	lfej.Target = lf.Target

	return &lfej
}

// UnMarshall does special handlig to coerce attributes that can be either single strings or
// string arrays into string arrays of one or more entnries
func (lf *LuciRPCFirewallEntry) postUnMarshall(lfej *LuciRPCFirewallEntryJSON) error {

	if lfej == nil {
		if lf != nil {
			*lf = LuciRPCFirewallEntry{}
		}
		return nil
	}
	lf.Anonymous = lfej.Anonymous
	lf.Index = lfej.Index
	lf.DotName = lfej.DotName
	lf.Type = lfej.Type
	lf.Dest = lfej.Dest
	lf.DestIP = assertStringArray(lfej.DestIP)
	lf.DestPort = assertStringArray(lfej.DestPort)
	lf.DropInvalid = lfej.DropInvalid
	lf.Enabled = lfej.Enabled
	lf.Family = lfej.Family
	lf.Forward = lfej.Forward
	lf.IcmpType = assertStringArray(lfej.IcmpType)
	lf.Input = lfej.Input
	lf.Limit = lfej.Limit
	lf.Masq = lfej.Masq
	lf.MtuFix = lfej.MtuFix
	lf.Name = lfej.Name
	lf.Network = assertStringArray(lfej.Network)
	lf.Output = lfej.Output
	lf.Path = lfej.Path
	lf.Proto = assertStringArray(lfej.Proto)
	lf.Src = lfej.Src
	lf.SrcDport = lfej.SrcDport
	lf.SrcIP = assertStringArray(lfej.SrcIP)
	lf.SrcMac = assertStringArray(lfej.SrcMac)
	lf.SrcPort = lfej.SrcPort
	lf.SynFlood = lfej.SynFlood
	lf.Target = lfej.Target

	return nil
}

type LuciRPCFirewallEntryJSON struct {
	Anonymous   bool        `json:".anonymous,omitempty"`
	Index       int64       `json:".index,omitempty"`
	DotName     string      `json:".name,omitempty"`
	Type        string      `json:".type,omitempty"`
	Dest        string      `json:"dest,omitempty"`
	DestIP      interface{} `json:"dest_ip,omitempty"`
	DestPort    interface{} `json:"dest_port,omitempty"`
	DropInvalid string      `json:"drop_invalid,omitempty"`
	Enabled     string      `json:"enabled,omitempty"`
	Family      string      `json:"family,omitempty"`
	Forward     string      `json:"forward,omitempty"`
	IcmpType    interface{} `json:"icmp_type,omitempty"`
	Input       string      `json:"input,omitempty"`
	Limit       string      `json:"limit,omitempty"`
	Masq        string      `json:"masq,omitempty"`
	MtuFix      string      `json:"mtu_fix,omitempty"`
	Name        string      `json:"name,omitempty"`
	Network     interface{} `json:"network,omitempty"`
	Output      string      `json:"output,omitempty"`
	Path        string      `json:"path,omitempty"`
	Proto       interface{} `json:"proto,omitempty"`
	Src         string      `json:"src,omitempty"`
	SrcDport    string      `json:"src_dport,omitempty"`
	SrcIP       interface{} `json:"src_ip,omitempty"`
	SrcMac      interface{} `json:"src_mac,omitempty"`
	SrcPort     string      `json:"src_port,omitempty"`
	SynFlood    string      `json:"syn_flood,omitempty"`
	Target      string      `json:"target,omitempty"`
}

//NewClient returns a new client
func NewClient(cfn string) (*Client, error) {
	cfgFileName := cfn
	cfgFile, err := os.Open(cfgFileName)
	if err != nil {
		//quietly try a fallback:
		//TODO: signal this with an error type for realerr
		realerr := err
		cfgFileName := "/etc/myhouse.json"
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

// Post is a sink to capture http.Client creation on posts
func (c *Client) Post(URL string, body []byte) (*http.Response, error) {
	r, err := doPost(URL, body) // 1st attempt
	if err != nil {
		return r, err
	}
	if r.StatusCode == 403 {
		// retry once
		c.Auth()
		r, err = doPost(URL, body)
	}
	return r, err
}

// doPost is a sink to capture http.Client creation on posts
// no re-auth so we can use it in Auth()
func doPost(URL string, body []byte) (*http.Response, error) {
	log.Printf("URL:\n%s\nBody\n%s\n", URL, string(body))

	// workaround: accept self-signed TLS
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	netClient := &http.Client{Transport: tr,
		Timeout: time.Second * 10,
	}
	newBody := []byte(body) //copy body in case of retry
	r, err := netClient.Post(URL, "application/json", bytes.NewBuffer(newBody))
	return r, err
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
	// fmt.Printf("body;\n%v", string(responseBody))
	result := LuciRPCSingleResponse{} //map[string]interface{}
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return fmt.Errorf("Failed to parse router Auth response:%v", err)
	}
	if len(result.Result) < 4 || len(result.Error) != 0 {
		a := c.AuthRequest
		a.Params[1] = ""
		log.Printf("Auth failed;%#v got %v", a, string(responseBody))
		return fmt.Errorf("Failed to get auth token")
	}
	log.Println("Client token obtained")
	c.AuthToken = result.Result

	return nil
}
func assertString(i interface{}) string {
	s, ok := i.(string)
	if ok {
		return s
	}
	return ""
}

//GetFirewallRules gets traffic rules from the router firewall
// NOTE in the return map, label names are used as the index but dotnames are the index in the payload
func (c *Client) GetFirewallRules() (LuciFirewallEntryMap, error) {
	request := LuciRPCMethodRequest{
		Method: "get_all",
		Params: []string{"firewall"},
	}
	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling rule request:%v", err)
	}
	URL := c.BaseURL + "/uci?auth=" + c.AuthToken
	// get answer
	response, err := c.Post(URL, body)
	if err != nil {
		return nil, fmt.Errorf("Failed to contact Router:%v", err)
	}
	defer response.Body.Close()

	// find token
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read router rules response:%v", err)
	}

	fwResponse := &LuciRPCFirewallGetAllResponse{}
	err = json.Unmarshal(responseBody, &fwResponse)
	//var rules map[string]LuciFirewallRuleResult
	if err != nil {
		return nil, fmt.Errorf("Failed to parse router rules response:%v", err)
	}
	rules := make(LuciFirewallEntryMap)
	// noname rule fixup - use label names as index
	for name, val := range *fwResponse.Result {
		if len(val.Name) == 0 {
			val.Name = name
		}
		v := LuciRPCFirewallEntry{}
		v.postUnMarshall(&val)
		if val.SrcMac != nil {
			log.Printf("%#v\n%#v\n", val.SrcMac, v.SrcMac)
		}
		rules[val.Name] = v
	}

	return rules, nil

}

// EnableFirewallRule uses tje donamme
func (c *Client) EnableFirewallRule(dotname string) error {
	request := LuciRPCMethodRequest{
		Method: "set",
		Params: []string{"firewall", dotname, "enabled", ""},
	}
	body, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("Error marshalling rule request:%v", err)
	}
	URL := c.BaseURL + "/uci?auth=" + c.AuthToken
	// get answer
	response, err := c.Post(URL, body)
	if err != nil {
		return fmt.Errorf("Failed to contact Router:%v", err)
	}
	defer response.Body.Close()

	// find token
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("Failed to read router rules response:%v", err)
	}
	fmt.Printf("body;\n%v", string(responseBody))
	result := LuciRPCBoolResponse{} //map[string]interface{}
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return fmt.Errorf("Failed to parse router rules response:%v", err)
	}
	if result.Result == false {
		return fmt.Errorf("Set rule failed:%v", result.Error)

	}

	return c.doCommitFirewall()
}

// DisableFirewallRule uses tje donamme
func (c *Client) DisableFirewallRule(dotname string) error {
	request := LuciRPCMethodRequest{
		Method: "set",
		Params: []string{"firewall", dotname, "enabled", "0"},
	}
	body, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("Error marshalling rule request:%v", err)
	}
	URL := c.BaseURL + "/uci?auth=" + c.AuthToken
	// get answer
	response, err := c.Post(URL, body)
	if err != nil {
		return fmt.Errorf("Failed to contact Router:%v", err)
	}
	defer response.Body.Close()

	// find token
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("Failed to read router rules response:%v", err)
	}
	fmt.Printf("body;\n%v", string(responseBody))
	result := LuciRPCBoolResponse{} //map[string]interface{}
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return fmt.Errorf("Failed to parse router rules response:%v", err)
	}
	if result.Result == false {
		return fmt.Errorf("Set rule failed:%v", result.Error)

	}

	return c.doCommitFirewall()
}

func (c *Client) doCommitFirewall() error {
	request := LuciRPCMethodRequest{
		Method: "commit",
		Params: []string{"firewall"},
	}
	body, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("Error marshalling commit request:%v", err)
	}
	URL := c.BaseURL + "/uci?auth=" + c.AuthToken
	// get answer
	response, err := c.Post(URL, body)
	if err != nil {
		return fmt.Errorf("Failed to contact Router:%v", err)
	}
	defer response.Body.Close()

	// find token
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("Failed to read router commmit response:%v", err)
	}
	fmt.Printf("body;\n%v", string(responseBody))
	result := LuciRPCBoolResponse{} //map[string]interface{}
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return fmt.Errorf("Failed to parse router commit response:%v", err)
	}
	if result.Result == false {
		return fmt.Errorf("Commit rule failed:%v", result.Error)
	}
	return nil
}

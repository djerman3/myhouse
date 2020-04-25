// Package client sends commands to an openwrt router for uci managment of router stuff
package client

//LuciRPCMethodRequest allows serialization of the struct to a valid request body
type LuciRPCMethodRequest struct {
	ID     int      `json:"id"`
	Method string   `json:"method"`
	Params []string `json:"params"`
}

//LuciRPCSingleResponse allows deserialization of the result body to a valid struct
// where we expect a simple name:string
type LuciRPCSingleResponse struct {
	ID     int    `json:"id"`
	Result string `json:"result"`
	Error  string `json:"error"`
}

//LuciRPCComplexResponse allows deserialization of the result body to a valid struct
// where we expect a simple name:string
type LuciRPCComplexResponse struct {
	ID     int                     `json:"id"`
	Result *map[string]interface{} `json:"result"`
	Error  string                  `json:"error"`
}

//Client remembers the client connection and auth
type Client struct {
	cfgFileName   string
	ServerAddress string
	BaseURL       string
	AuthToken     string
	AuthRequest   *LuciRPCMethodRequest // we're going to need this one  alot
}

//NewClient returns a new client
func NewClient(cfn string) (*Client, error) {
	cfgFileName = cfn

	c := &Client{address, "/cgi-bin/luci/RPC", "", nil}
	return c, nil
}

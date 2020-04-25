// Package client sends commands to an openwrt router for uci managment of router stuff
package client

type LuciRpcAuthRequest struct {
	ID     int      `json:"id"`
	Method string   `json:"method"`
	Params []string `json:"params"`
}
type LuciRpcAuthResponse struct {
	ID     int    `json:"id"`
	Result string `json:"result"`
	Error  string `json:"error"`
}

type Client struct {
	ServerAddress string
	BaseURL       string
	AuthToken     string              `json:"authtoken"`
	AuthRequest   *LuciRpcAuthRequest `json:ref`
}

//NewClient returns a new client
func NewClient(address string) (*Client, error) {
	c := &Client{address, "/cgi-bin/luci/rpc", "", nil}
	return c, nil
}

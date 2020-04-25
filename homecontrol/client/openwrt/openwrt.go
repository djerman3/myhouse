// Package openwrt/client sends commands to an openwrt router for uci managment of router stuff
package client

type LuciRpcMethodRequest struct {
	ID     int      `json:"id"`
	Method string   `json:"method"`
	Params []string `json:"params"`
}
type LuciRpcAuthResponse struct {
	ID     int    `json:"id"`
	Result string `json:"result"`
	Error  string `json:"error"`
}
type LuceFirewallRuleResult{
	DotName string `json:".name"`
	Name string `json:"name"`
	Enabled string `json:"enabled"`
}

package coin

////////////////////////////////////////////////////////////////////////////////
/*

TODO:

	Other coins might have a websocket based RPC protocol (bitcoin, dash(?)),
	and it would be appropriate for the RPC access to be abstracted behind an
	interface which selects the type of transport based on the coin's config
	parameters as it is being registered.

*/
////////////////////////////////////////////////////////////////////////////////

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
)

////////////////////////////////////////////////////////////////////////////////

var (
	rpcId int64 // Atomic counter for JSON RPC unique ID
)

////////////////////////////////////////////////////////////////////////////////

// DoJSONRPCCommand accepts a `method` and a list of values in `params` which
// will be sent over JSON RPC to the corresponding coin's daemon.
func (c *Coin) DoJSONRPCCommand(method string, params []interface{}) (*http.Response, error) {
	url := fmt.Sprintf("http://%s:%d", c.GetConfigValue("rpcallowip"), c.GetRPCPort())

	atomic.AddInt64(&rpcId, 1)
	dto := &struct {
		Method string        `json:"method"`
		ID     int64         `json:"id"`
		Params []interface{} `json:"params"`
	}{
		Method: method,
		ID:     rpcId,
		Params: params,
	}
	bs, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.GetConfigValue("rpcuser"), c.GetConfigValue("rpcpassword"))
	return http.DefaultClient.Do(req)
}

////////////////////////////////////////////////////////////////////////////////

func init() {
	atomic.StoreInt64(&rpcId, 0)
}

////////////////////////////////////////////////////////////////////////////////

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
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync/atomic"
)

////////////////////////////////////////////////////////////////////////////////

var (
	rpcId int64 // Atomic counter for JSON RPC unique ID
)

////////////////////////////////////////////////////////////////////////////////

var (
	ErrCouldNotConnectToServer = errors.New("couldn't connect to server")
	ErrAuthorizationFailed     = errors.New("authorization failure, bad rpcuser/rpcpass")
	ErrNoResponse              = errors.New("no response from server")
)

////////////////////////////////////////////////////////////////////////////////

type JSONRPCError struct {
	Code    int64  `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type JSONRPCResponse struct {
	ID     int64                  `json:"id,omitempty"`
	Result map[string]interface{} `json:"result,omitempty"`
	Error  JSONRPCError           `json:"error,omitempty"`
}

////////////////////////////////////////////////////////////////////////////////

// DoJSONRPCCommand accepts a `method` and a list of values in `params` which
// will be sent over JSON RPC to the corresponding coin's daemon.
func (c *Coin) DoJSONRPCCommand(method string, params []interface{}) (*JSONRPCResponse, error) {
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
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, ErrCouldNotConnectToServer
	}
	if rsp.StatusCode == http.StatusUnauthorized {
		return nil, ErrAuthorizationFailed
	}

	// Validate the response and return it to the caller.  NOTE: it is no the job
	// of this function to verify any RPC errors, this is just the transport for the
	// packet.
	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil || len(data) == 0 {
		return nil, ErrNoResponse
	}

	jrrsp := &JSONRPCResponse{}
	if err := json.Unmarshal(data, jrrsp); err != nil {
		return nil, ErrNoResponse
	}
	return jrrsp, nil
}

////////////////////////////////////////////////////////////////////////////////

func init() {
	atomic.StoreInt64(&rpcId, 0)
}

////////////////////////////////////////////////////////////////////////////////

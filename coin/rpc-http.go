package coin

////////////////////////////////////////////////////////////////////////////////

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
)

var (
	rpcId int64
)

////////////////////////////////////////////////////////////////////////////////

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

package chaingo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	host       string
	port       int32
	httpclient http.Client
}

func new(host string, port int32) (*Client, error) {
	httpcli := http.DefaultClient
	return &Client{
		host:       host,
		port:       port,
		httpclient: *httpcli,
	}, nil
}

func NewClient(host string, port int32) (*Client, error) {
	return new(host, port)
}

type httpClientRequest struct {
	Method string         `json:"method"`
	Params [1]interface{} `json:"params"`
	ID     uint64         `json:"id"`
}

type httpClientResponse struct {
	ID     uint64           `json:"id"`
	Result *json.RawMessage `json:"result"`
	Error  interface{}      `json:"error"`
}

func (client *Client) Call(method string, params interface{}) (interface{}, error) {
	req := &httpClientRequest{}
	req.Method = method
	req.Params[0] = params
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("http://%s:%d", client.host, client.port)
	postres, err := client.httpclient.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer postres.Body.Close()
	res := &httpClientResponse{}
	result, err := ioutil.ReadAll(postres.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(result, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

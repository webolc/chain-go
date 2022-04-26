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
	result, err := ioutil.ReadAll(postres.Body)
	return result, nil
}

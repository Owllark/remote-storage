package http_client

import (
	"bytes"
	"net/http"
)

/*type Client interface {
	Get()
	Put()
	Post()
	Delete()
}*/

type HttpClient struct {
	client     http.Client
	serviceUrl string
}

func NewHttpClient(url string) *HttpClient {
	client := http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}
	res := HttpClient{client: client, serviceUrl: url}
	return &res
}

func (c *HttpClient) Get(api string) *http.Response {
	resp, _ := http.Get(c.serviceUrl + api)
	return resp
}

func (c *HttpClient) Post(api string, contentType string, body []byte) *http.Response {
	resp, _ := http.Post(c.serviceUrl+api, contentType, bytes.NewReader(body))
	return resp
}

func (c *HttpClient) DoRequest(request *http.Request) *http.Response {
	resp, _ := c.client.Do(request)
	return resp
}

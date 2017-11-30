// Bitstamp API client
package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

type Response struct {
	Err       []string `json:"-"`
	High      string   `json:"high"`
	Last      string   `json:"last"`
	Timestamp string   `json:"timestamp"`
	Bid       string   `json:"bid"`
	Low       string   `json:"low"`
	Ask       string   `json:"ask"`
	Open      string   `json:"open"`
}

type Client struct {
	url    string
	client *http.Client
}

func NewClient() *Client {
	c := new(Client)
	// hardcoded for now
	c.url = "https://www.bitstamp.net/api/v2/"
	c.client = &http.Client{
		Timeout: time.Second * 5,
	}
	return c
}

func (c *Client) GetTicker(crypto, fiat string) *Response {
	crypto = strings.ToLower(crypto)
	fiat = strings.ToLower(fiat)
	url := c.url + "ticker/" + crypto + fiat

	resp := new(Response)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		resp.Err = make([]string, 1, 1)
		resp.Err[0] = err.Error()
		return resp
	}

	res, err := c.client.Do(req)
	if err != nil {
		resp.Err = make([]string, 1, 1)
		resp.Err[0] = err.Error()
		return resp
	}
	if res.StatusCode != 200 {
		resp.Err = make([]string, 1, 1)
		resp.Err[0] = "Unexpected API response"
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		resp.Err = make([]string, 1, 1)
		resp.Err[0] = err.Error()
		return resp
	}

	log.WithFields(log.Fields{
		"body": string(body),
	}).Debug("Raw response for Ticker")

	err = json.Unmarshal(body, &resp)
	if err != nil {
		resp.Err = make([]string, 1, 1)
		resp.Err[0] = err.Error()
		return resp
	}

	return resp
}

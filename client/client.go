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
	resp := new(Response)

	body, err := c.sendRequest(c.url + "ticker/" + crypto + fiat)
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
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error unmarshaling response data.")
		resp.Err = make([]string, 1, 1)
		resp.Err[0] = "Error unmarshaling response data."
		return resp
	}

	return resp
}

func (c *Client) sendRequest(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

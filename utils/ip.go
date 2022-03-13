package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type IPRequestData struct {
	IP          string `json:"ip"`
	CountryCode string `json:"country_code"`
	Tor         bool   `json:"tor"`
	Reverse     string `json:"reverse"`
}

type IPChecker struct {
	HTTPClient *http.Client
}

func (c *IPChecker) request(req *http.Request) (*http.Response, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *IPChecker) GetPublicIPv4Address() (address string, err error) {
	return c.getIPAddress(4)
}

func (c *IPChecker) GetPublicIPv6Address() (address string, err error) {
	return c.getIPAddress(6)
}

func (c *IPChecker) getIPAddress(version int) (address string, err error) {
	url := fmt.Sprintf("https://ipv%v.ipleak.net/json/", version)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	logrus.Info("Requesting IP address from ", url)
	resp, err := c.request(req)
	if err != nil {
		return
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	data := IPRequestData{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return
	}

	address = data.IP
	return
}

func NewIPChecker() (ipChecker *IPChecker, err error) {
	ipChecker = &IPChecker{
		HTTPClient: &http.Client{
			Timeout: 3 * time.Second},
	}
	return
}

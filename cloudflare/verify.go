package cloudflare

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

func (c *CloudFlareClient) VerifyToken() (response CloudFlareTokenVerify, err error) {
	url := fmt.Sprintf("%suser/tokens/verify", c.BaseUrl)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

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

	var data CloudFlareTokenVerify
	err = json.Unmarshal(body, &data)
	if err != nil {
		return
	}

	logrus.Infof("CloudFlare token verification: %+v", data)
	return
}

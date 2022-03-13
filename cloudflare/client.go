package cloudflare

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

type CloudFlareClient struct {
	BaseUrl    string
	apiToken   string
	limiter    *rate.Limiter
	HTTPClient *http.Client
}

func (c *CloudFlareClient) request(req *http.Request) (*http.Response, error) {
	ctx := context.Background()
	err := c.limiter.Wait(ctx)
	req.Header.Add("Authorization", "Bearer "+c.apiToken)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	logrus.Info("Requesting ", req.URL.String())
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *CloudFlareClient) GetZone(name string) (zone Zone, err error) {
	req, err := http.NewRequest("GET", c.BaseUrl+"zones", nil)
	if err != nil {
		return
	}

	query := req.URL.Query()
	query.Add("name", name)
	req.URL.RawQuery = query.Encode()

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
	parsedResponse := CloudFlareResponse{}
	err = json.Unmarshal(body, &parsedResponse)
	if err != nil {
		return
	}

	for _, item := range body {
		zone = Zone{}
		err = json.Unmarshal(body, &item)
		if err != nil {
			return
		}
		if zone.Name == name {
			return zone, nil
		}
	}
	return
}

func NewCloudFlareClient(apiToken string) *CloudFlareClient {
	BaseUrl := "https://api.cloudflare.com/client/v4/"
	limiter := rate.NewLimiter(rate.Every(5*time.Minute), 1200)
	return &CloudFlareClient{
		BaseUrl:  BaseUrl,
		apiToken: apiToken,
		limiter:  limiter,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second},
	}
}

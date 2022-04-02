package cloudflare

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

	logrus.Infof("Calling %s %s", req.Method, req.URL.String())
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *CloudFlareClient) GetZone(name string) (zone CloudFlareZone, err error) {
	url := c.BaseUrl + "zones"
	req, err := http.NewRequest("GET", url, nil)
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

	parsedResponse := CloudFlareZonesResponse{}
	err = json.Unmarshal(body, &parsedResponse)
	if err != nil {
		return
	}

	for _, zone := range parsedResponse.Result {
		if zone.Name == name {
			return zone, nil
		}
	}
	return
}

func (c *CloudFlareClient) VerifyToken() (response CloudFlareApiTokenVerifyResponse, err error) {
	url := c.BaseUrl + "user/tokens/verify"
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

	var data CloudFlareApiTokenVerifyResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return
	}
	return data, nil
}

func (c *CloudFlareClient) CheckDNSRecordAlreadyExists(zoneID string, recordType string, recordName string) (exists bool, err error) {
	record, err := c.GetDNSRecord(zoneID, recordType, recordName)
	if err != nil {
		return false, nil
	} else if record.Name == "" {
		return false, nil
	}
	return true, nil
}

func (c *CloudFlareClient) GetDNSRecord(zoneID string, recordType string, recordName string) (record CloudFlareDNSRecord, err error) {
	url := c.BaseUrl + fmt.Sprintf("zones/%s/dns_records", zoneID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	query := req.URL.Query()
	query.Add("name", recordName)
	query.Add("type", recordType)
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

	parsedResponse := CloudFlareDNSRecordsResponse{}
	err = json.Unmarshal(body, &parsedResponse)
	if err != nil {
		return
	}

	for _, record := range parsedResponse.Result {
		if record.Name == recordName && record.Type == recordType {
			logrus.Debugf("Found DNS Record for %s - %s", recordType, recordName)
			return record, nil
		}
	}
	logrus.Debugf("Did not found DNS Record for %s - %s", recordType, recordName)
	return
}

func (c *CloudFlareClient) CreateDNSRecord(zoneID string, record *CloudFlareDNSRecordBody) (err error) {
	method := "POST"
	url := c.BaseUrl + fmt.Sprintf("zones/%s/dns_records", zoneID)
	_, err = c.editDNSRecord(method, url, zoneID, record)
	return
}

func (c *CloudFlareClient) UpdateDNSRecord(zoneID string, recordID string, record *CloudFlareDNSRecordBody) (err error) {
	method := "PUT"
	url := c.BaseUrl + fmt.Sprintf("zones/%s/dns_records/%s", zoneID, recordID)
	_, err = c.editDNSRecord(method, url, zoneID, record)
	return
}

func (c *CloudFlareClient) editDNSRecord(method string, url string, zoneID string, record *CloudFlareDNSRecordBody) (respRecord CloudFlareDNSRecord, err error) {
	jsonString, err := json.Marshal(record)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonString))
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

	var data CloudFlareEditDNSRecordResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return
	}

	if !data.Success {
		logrus.Errorf("%v", data.Errors)
		return data.Result, errors.New("could not create dns entry")
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

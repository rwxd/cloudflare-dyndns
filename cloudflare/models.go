package cloudflare

import (
	"encoding/json"
	"time"
)

type CloudFlareResponse struct {
	Success  bool                     `json:"success"`
	Errors   []CloudFlareResponseInfo `json:"errors"`
	Messages []CloudFlareResponseInfo `json:"messages"`
}

type CloudFlareResponseInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type CloudFlareResultInfo struct {
	Page       int                         `json:"page"`
	PerPage    int                         `json:"per_page"`
	TotalPages int                         `json:"total_pages"`
	Count      int                         `json:"count"`
	Total      int                         `json:"total_count"`
	Cursor     string                      `json:"cursor"`
	Cursors    CloudFlareResultInfoCursors `json:"cursors"`
}

type CloudFlareResultInfoCursors struct {
	Before string `json:"before"`
	After  string `json:"after"`
}

type CloudFlareRawResponse struct {
	CloudFlareResponse
	Result json.RawMessage `json:"result"`
}

type CloudFlareApiTokenVerifyResponse struct {
	CloudFlareResponse
	Result struct {
		ID        string    `json:"id"`
		Status    string    `json:"status"`
		NotBefore time.Time `json:"not_before"`
		ExpiresOn time.Time `json:"expires_on"`
	}
}

type CloudFlareZonesResponse struct {
	CloudFlareResponse
	Result               []CloudFlareZone `json:"result"`
	CloudFlareResultInfo `json:"result_info"`
}

type CloudFlareZone struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type CloudFlareDNSRecordsResponse struct {
	CloudFlareResponse
	Result               []CloudFlareDNSRecord `json:"result"`
	CloudFlareResultInfo `json:"result_info"`
}

type CloudFlareDNSRecord struct {
	CreatedOn  time.Time   `json:"created_on,omitempty"`
	ModifiedOn time.Time   `json:"modified_on,omitempty"`
	Type       string      `json:"type,omitempty"`
	Name       string      `json:"name,omitempty"`
	Content    string      `json:"content,omitempty"`
	Meta       interface{} `json:"meta,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	ID         string      `json:"id,omitempty"`
	ZoneID     string      `json:"zone_id,omitempty"`
	ZoneName   string      `json:"zone_name,omitempty"`
	Priority   *uint16     `json:"priority,omitempty"`
	TTL        int         `json:"ttl,omitempty"`
	Proxied    *bool       `json:"proxied,omitempty"`
	Proxiable  bool        `json:"proxiable,omitempty"`
	Locked     bool        `json:"locked,omitempty"`
}

type CloudFlareDNSRecordBody struct {
	Type    string
	Name    string
	Content string
	Ttl     int
}

type CloudFlareEditDNSRecordResponse struct {
	CloudFlareResponse
	Result CloudFlareDNSRecord `json:"result"`
}

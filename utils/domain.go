package utils

import "strings"

func RemoveZoneFromDomainName(domain string, zone string) string {
	if domain == zone {
		return domain
	}
	domain = strings.Replace(domain, zone, "", 1)
	domain = strings.TrimSuffix(domain, ".")
	return domain
}

func CombineRecordAndZone(record string, zone string) string {
	if strings.Contains(record, zone) {
		return record
	}
	zone = strings.Trim(zone, ".")
	record = strings.Trim(record, ".")
	return record + "." + zone
}

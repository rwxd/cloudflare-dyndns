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

package utils

import "testing"

func TestRemoveZoneFromDomainName(t *testing.T) {
	tables := []struct {
		domain   string
		zone     string
		expected string
	}{
		{"test.example.com", "example.com", "test"},
		{"example", "example.com", "example"},
		{"example.com", "example.com", "example.com"},
	}

	for _, table := range tables {
		actual := RemoveZoneFromDomainName(table.domain, table.zone)
		if actual != table.expected {
			t.Errorf("RemoveZoneFromDomainName(%s, %s) == %s, expected %s", table.domain, table.zone, actual, table.expected)
		}
	}
}

func TestCombineRecordAndZone(t *testing.T) {
	tables := []struct {
		record   string
		zone     string
		expected string
	}{
		{"test", "example.com", "test.example.com"},
		{"test", ".example.com", "test.example.com"},
		{".test", "example.com", "test.example.com"},
		{"test.example.com", "example.com", "test.example.com"},
	}

	for _, table := range tables {
		actual := CombineRecordAndZone(table.record, table.zone)
		if actual != table.expected {
			t.Errorf("RemoveZoneFromDomainName(%s, %s) == %s, expected %s", table.record, table.zone, actual, table.expected)
		}
	}
}

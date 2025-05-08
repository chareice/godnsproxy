package dns

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/chareice/godnsproxy/domain"
)

func TestGetDomainFromDnsQuery(t *testing.T) {
	// Example DNS query for "example.com"
	query := []byte{
		0x00, 0x00, // Transaction ID
		0x01, 0x00, // Flags
		0x00, 0x01, // Questions
		0x00, 0x00, // Answer RRs
		0x00, 0x00, // Authority RRs
		0x00, 0x00, // Additional RRs
		// Query section
		0x07, 'e', 'x', 'a', 'm', 'p', 'l', 'e',
		0x03, 'c', 'o', 'm',
		0x00,       // Null terminator
		0x00, 0x01, // QTYPE (A record)
		0x00, 0x01, // QCLASS (IN)
	}

	domain := getDomainFromDnsQuery(query)
	if domain != "example.com" {
		t.Errorf("Expected example.com, got %s", domain)
	}
}

func TestServerRouting(t *testing.T) {
	mgr := domain.NewManager()
	// Create temp domains file for testing
	tmpFile := filepath.Join(t.TempDir(), "domains.txt")
	os.WriteFile(tmpFile, []byte("test.cn"), 0644)
	mgr.LoadDomains(tmpFile)

	server := &Server{
		config: &Config{
			ChinaServer:    "223.5.5.5",
			TrustDNSServer: "https://1.1.1.1/dns-query",
		},
		domainMgr: mgr,
	}

	testCases := []struct {
		domain      string
		expectChina bool
	}{
		{"test.cn", true},
		{"google.com", false},
	}

	for _, tc := range testCases {
		t.Run(tc.domain, func(t *testing.T) {
			isChina := server.domainMgr.IsChinaDomain(tc.domain)
			if isChina != tc.expectChina {
				t.Errorf("For domain %s, expected China=%v but got %v",
					tc.domain, tc.expectChina, isChina)
			}
		})
	}
}

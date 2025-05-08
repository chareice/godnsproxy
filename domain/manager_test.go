package domain

import (
	"path/filepath"
	"testing"
)

func TestLoadDomains(t *testing.T) {
	mgr := NewManager()
	testFile := filepath.Join("..", "tests", "fixtures", "domains.txt")

	err := mgr.LoadDomains(testFile)
	if err != nil {
		t.Fatalf("Failed to load domains: %v", err)
	}

	testCases := []struct {
		domain   string
		expected bool
	}{
		{"example.com", true},
		{"test.cn", true},
		{"example.org", true},
		{"google.com", false},
		{"sub.test.cn", true},
		{"nonexistent.com", false},
	}

	for _, tc := range testCases {
		t.Run(tc.domain, func(t *testing.T) {
			result := mgr.IsChinaDomain(tc.domain)
			if result != tc.expected {
				t.Errorf("For domain %s, expected %v but got %v",
					tc.domain, tc.expected, result)
			}
		})
	}
}

func TestEmptyManager(t *testing.T) {
	mgr := NewManager()

	if mgr.IsChinaDomain("any.com") {
		t.Error("Empty manager should return false for any domain")
	}
}

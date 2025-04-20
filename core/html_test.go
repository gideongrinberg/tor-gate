package core

import (
	"strings"
	"testing"
)

func TestRewriteLinks(t *testing.T) {
	tests := []struct {
		name        string
		inputHTML   string
		proxyDomain string
		expected    string
	}{
		{
			name:        "Rewrite single .onion link",
			inputHTML:   `<html><body><a href="http://example.onion/page">Link</a></body></html>`,
			proxyDomain: "proxy.example.com",
			expected:    `<html><body><a href="http://example.proxy.example.com/page">Link</a></body></html>`,
		},
		{
			name:        "No .onion link present",
			inputHTML:   `<html><body><a href="http://example.com/page">Link</a></body></html>`,
			proxyDomain: "proxy.example.com",
			expected:    `<html><body><a href="http://example.com/page">Link</a></body></html>`,
		},
		{
			name:        "Multiple .onion links",
			inputHTML:   `<html><body><a href="http://first.onion/page1">First</a><a href="http://second.onion/page2">Second</a></body></html>`,
			proxyDomain: "proxy.example.com",
			expected:    `<html><body><a href="http://first.proxy.example.com/page1">First</a><a href="http://second.proxy.example.com/page2">Second</a></body></html>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modifiedHTML := RewriteLinks([]byte(tt.inputHTML), tt.proxyDomain)
			output := strings.TrimSpace(string(modifiedHTML))
			expected := strings.TrimSpace(tt.expected)
			if output != expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", expected, output)
			}
		})
	}
}

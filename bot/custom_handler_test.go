package bot

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMatchForUrl(t *testing.T) {
	testCase := []struct {
		message string
		url     string
	}{
		{"http://google.com", "http://google.com"},
		{"https://www.google.com", "https://www.google.com"},
		{"Some https://arbitrary.link.org to be checked", "https://arbitrary.link.org"},
		{"Special chars https://tools.ietf.org/html/rfc1459#section-2.3", "https://tools.ietf.org/html/rfc1459#section-2.3"},
	}

	for _, tc := range testCase {
		assert.Equal(t, tc.url, matchForUrl(tc.message))

	}
}

package bot

import (
	"github.com/stretchr/testify/assert"
	"strings"
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
		{"No link at all", ""},
	}

	for _, tc := range testCase {
		assert.Equal(t, tc.url, matchForUrl(tc.message))

	}
}

func TestGetTitle(t *testing.T) {
	testCase := []struct {
		content string
		title   string
	}{
		{`"<html><head><title>This is the title</title></head><body>This is the body</body></html>"`, "This is the title"},
		{`"<html><head><title>This is <b>the</b> title</title></head><body>This is the body</body></html>"`, "This is &lt;b&gt;the&lt;/b&gt; title"},
		{`"<html><head><title></title></head><body>This is the body</body></html>"`, ""},
		{`"<html><head></head><body>This is the body</body></html>"`, ""},
		{`"<html>incomplete"`, ""},
	}

	for _, tc := range testCase {
		title, err := getTitle(strings.NewReader(tc.content))
		assert.Nil(t, err)
		assert.Equal(t, tc.title, title)
	}
}

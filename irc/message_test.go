package irc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractNormalReply(t *testing.T) {
	assert := assert.New(t)

	in := ":adams.freenode.net NOTICE * :*** No Ident response"
	out, ok := extractNormalReply(in)

	if !ok {
		t.Fatal("Failed to extract values")
	}

	assert.Equal("adams.freenode.net", out.Prefix)
	assert.Equal("NOTICE", out.Command)
	assert.Equal("* :*** No Ident response", out.Parameter)
}

func TestExtractStatusReply(t *testing.T) {
	assert := assert.New(t)

	in := "PING :adams.freenode.net"
	out, ok := extractStatusReply(in)

	if !ok {
		t.Fatal("Failed to extract values")
	}

	assert.Equal("", out.Prefix)
	assert.Equal("PING", out.Command)
	assert.Equal("adams.freenode.net", out.Parameter)
}

func TestExtractResponseNormalReply(t *testing.T) {
	assert := assert.New(t)

	in := ":adams.freenode.net NOTICE * :*** No Ident response"
	out := ExtractResponse(in)

	assert.Equal("adams.freenode.net", out.Prefix)
	assert.Equal("NOTICE", out.Command)
	assert.Equal("* :*** No Ident response", out.Parameter)
}

func TestExtractResponseStatusReply(t *testing.T) {
	assert := assert.New(t)

	in := "PING :adams.freenode.net\n"
	out := ExtractResponse(in)

	assert.Equal("", out.Prefix)
	assert.Equal("PING", out.Command)
	assert.Equal("adams.freenode.net", out.Parameter)
}

func TestExtractResponseUnconventional(t *testing.T) {
	assert := assert.New(t)

	in := "This is not an IRC server"
	out := ExtractResponse(in)

	assert.Equal("", out.Prefix)
	assert.Equal("", out.Command)
	assert.Equal("This is not an IRC server", out.Parameter)
}

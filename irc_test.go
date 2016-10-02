package main

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

	assert.Equal("adams.freenode.net", out.initiator)
	assert.Equal("NOTICE", out.command)
	assert.Equal("* :*** No Ident response", out.message)
}

func TestExtractStatusReply(t *testing.T) {
	assert := assert.New(t)

	in := "PING :adams.freenode.net"
	out, ok := extractStatusReply(in)

	if !ok {
		t.Fatal("Failed to extract values")
	}

	assert.Equal("", out.initiator)
	assert.Equal("PING", out.command)
	assert.Equal("adams.freenode.net", out.message)
}

func TestExtractResponseNormalReply(t *testing.T) {
	assert := assert.New(t)

	in := ":adams.freenode.net NOTICE * :*** No Ident response"
	out := extractResponse(in)

	assert.Equal("adams.freenode.net", out.initiator)
	assert.Equal("NOTICE", out.command)
	assert.Equal("* :*** No Ident response", out.message)
}

func TestExtractResponseStatusReply(t *testing.T) {
	assert := assert.New(t)

	in := "PING :adams.freenode.net\n"
	out := extractResponse(in)

	assert.Equal("", out.initiator)
	assert.Equal("PING", out.command)
	assert.Equal("adams.freenode.net", out.message)
}

func TestExtractResponseUnconventional(t *testing.T) {
	assert := assert.New(t)

	in := "This is not an IRC server"
	out := extractResponse(in)

	assert.Equal("", out.initiator)
	assert.Equal("", out.command)
	assert.Equal("This is not an IRC server", out.message)
}

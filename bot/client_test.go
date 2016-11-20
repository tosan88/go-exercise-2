package bot

import (
	"testing"

	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tosan88/go-exercise-2/irc"
	"log"
)

type mockConn struct {
	respCh chan string
}

type testCase struct {
	name      string
	msg       irc.Message
	response  string
	client    Client
	logOutput string
}

func TestHandleMessage(t *testing.T) {
	assert := assert.New(t)

	listenCh := make(chan string)
	clientConn := &mockConn{listenCh}
	defer clientConn.Close()

	go clientConn.ReadContinuously(listenCh)

	ircChannel := "test-your-client"
	testCases := []testCase{
		{"logged successful join to server, command 001",
			irc.Message{Command: "001"},
			"",
			Client{conn: clientConn, config: &Conf{Server: "test"}, handlers: getHandlers()},
			"Successfully joined to server test"},
		{"send request to join channel, command 376",
			irc.Message{Command: "376"},
			fmt.Sprintf("JOIN #%v\n", ircChannel),
			Client{conn: clientConn, config: &Conf{Channel: ircChannel}, handlers: getHandlers()},
			""},
		{"send pong response",
			irc.Message{Command: "PING", Parameter: "1234"},
			"PONG :1234\n",
			Client{conn: clientConn, handlers: getHandlers()},
			"Sending PONG response"},
		{"log successful join to channel",
			irc.Message{Command: "JOIN", Prefix: "test-bot!home@home"},
			"",
			Client{conn: clientConn, registeredBotName: "test-bot", config: &Conf{Channel: ircChannel}, handlers: getHandlers()},
			fmt.Sprintf("Successfully joined to channel #%v as %v\n", ircChannel, "test-bot")},
		{"send new nick upon collision, command ERR_NICKCOLLISION",
			irc.Message{Command: "ERR_NICKCOLLISION"},
			"NICK test-bot",
			Client{conn: clientConn, registeredBotName: "test-bot", config: &Conf{BotName: "test-bot"}, handlers: getHandlers()},
			"Bot name could not be used. Adding suffix"},
		{"send new nick upon nick name used, command ERR_NICKNAMEINUSE",
			irc.Message{Command: "ERR_NICKNAMEINUSE"},
			"NICK test-bot",
			Client{conn: clientConn, registeredBotName: "test-bot", config: &Conf{BotName: "test-bot"}, handlers: getHandlers()},
			"Bot name could not be used. Adding suffix"},
		{"send new nick upon nick name used, command 433",
			irc.Message{Command: "433"},
			"NICK test-bot",
			Client{conn: clientConn, registeredBotName: "test-bot", config: &Conf{BotName: "test-bot"}, handlers: getHandlers()},
			"Bot name could not be used. Adding suffix"},
	}

	for _, tc := range testCases {
		t.Logf("Testing %v\n", tc.name)
		runTestCase(&tc, assert, listenCh)
	}

}

func runTestCase(tc *testCase, assert *assert.Assertions, listenCh chan string) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	tc.client.handleCommand(&tc.msg)

	if tc.response != "" {
		response := <-listenCh
		assert.Contains(response, tc.response)
	}

	if tc.logOutput != "" {
		assert.Contains(buf.String(), tc.logOutput)
	}

}

func (conn *mockConn) ReadContinuously(respCh chan string) {
	for {
		serverResponse, open := <-conn.respCh
		if !open {
			return
		}
		respCh <- serverResponse
	}
}

func (conn *mockConn) Send(msg string) {
	conn.respCh <- msg
}

func (conn *mockConn) Close() error {
	close(conn.respCh)
	return nil
}

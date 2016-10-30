package main

import (
	"testing"

	"bufio"
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"net"
	"sync"
	"time"
	"io"
)

type testCase struct {
	name      string
	msg       ircMessage
	response  string
	client    botClient
	logOutput string
}

func TestHandleMessage(t *testing.T) {
	assert := assert.New(t)

	listenCh := make(chan string)
	var wg sync.WaitGroup
	wg.Add(1)
	go startServer(&wg, assert, listenCh)

	wg.Wait()
	time.Sleep(time.Second) //do we need to make sure the server is accepting connections?

	clientConn, err := net.Dial("tcp", ":3000")
	assert.Nil(err)
	defer func() {
		if clientConn != nil {
			clientConn.Close()
		}
		fmt.Println("Client connection stopped")
	}()

	ircChannel := "test-your-client"
	testCases := []testCase{
		{"logged successful join to server, command 001",
			ircMessage{command: "001"},
			"",
			botClient{conn: clientConn, config: &conf{server: "test"}},
			"Successfully joined to server test"},
		{"logged successful join to server, command RPL_WELCOME",
			ircMessage{command: "RPL_WELCOME"},
			"",
			botClient{conn: clientConn, config: &conf{server: "test"}},
			"Successfully joined to server test"},
		{"send request to join channel, command 376",
			ircMessage{command: "376"},
			fmt.Sprintf("JOIN #%v\n", ircChannel),
			botClient{conn: clientConn, config: &conf{channel: ircChannel}},
			""},
		{"send request to join channel, command RPL_ENDOFMOTD",
			ircMessage{command: "RPL_ENDOFMOTD"},
			fmt.Sprintf("JOIN #%v\n", ircChannel),
			botClient{conn: clientConn, config: &conf{channel: ircChannel}},
			""},
		{"send pong response",
			ircMessage{command: "PING", message: "1234"},
			"PONG :1234\n",
			botClient{conn: clientConn},
			"Sending PONG response"},
		{"send greeting to newcomer",
			ircMessage{command: "JOIN", initiator: "newcomer!home@home"},
			fmt.Sprintf("PRIVMSG #%v :Welcome in this channel newcomer\n", ircChannel),
			botClient{conn: clientConn, registeredBotName: "test-bot", config: &conf{channel: ircChannel}, names: make(map[string]*user)},
			""},
		{"log successful join to channel",
			ircMessage{command: "JOIN", initiator: "test-bot!home@home"},
			"",
			botClient{conn: clientConn, registeredBotName: "test-bot", config: &conf{channel: ircChannel}},
			fmt.Sprintf("Successfully joined to channel #%v as %v\n", ircChannel, "test-bot")},
		{"send new nick upon collision, command ERR_NICKCOLLISION",
			ircMessage{command: "ERR_NICKCOLLISION"},
			"NICK test-bot",
			botClient{conn: clientConn, registeredBotName: "test-bot", config: &conf{botName: "test-bot"}},
			"Bot name could not be used. Adding suffix"},
		{"send new nick upon nick name used, command ERR_NICKNAMEINUSE",
			ircMessage{command: "ERR_NICKNAMEINUSE"},
			"NICK test-bot",
			botClient{conn: clientConn, registeredBotName: "test-bot", config: &conf{botName: "test-bot"}},
			"Bot name could not be used. Adding suffix"},
		{"send new nick upon nick name used, command 433",
			ircMessage{command: "433"},
			"NICK test-bot",
			botClient{conn: clientConn, registeredBotName: "test-bot", config: &conf{botName: "test-bot"}},
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
	tc.client.handleMessageCommand(&tc.msg)

	if tc.response != "" {
		response := <-listenCh
		assert.Contains(response, tc.response)
	}

	if tc.logOutput != "" {
		assert.Contains(buf.String(), tc.logOutput)
	}

}

func startServer(wg *sync.WaitGroup, assert *assert.Assertions, listenCh chan string) {
	mockConn, err := net.Listen("tcp", ":3000")
	assert.Nil(err)
	defer mockConn.Close()

	wg.Done()
	conn, err := mockConn.Accept()
	if err != nil {
		return
	}
	defer func() {
		if conn != nil {
			conn.Close()
		}
		fmt.Println("Server connection stopped")
	}()
	for {
		buf := bufio.NewReader(conn)
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		}
		assert.Nil(err)

		listenCh <- line
	}
}

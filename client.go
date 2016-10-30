package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
)

type handler func(*botClient, *ircMessage)

type user struct {
	name      string
	available bool
	lastSeen  time.Time
}

type botClient struct {
	config            *conf
	registeredBotName string
	conn              net.Conn
	response          chan string
	shouldStop        chan bool
	namesDB           *sql.DB
	handlers          map[string]handler
}

func newClient(config *conf, namesDB *sql.DB) *botClient {
	return &botClient{
		config:            config,
		registeredBotName: config.botName,
		response:          make(chan string),
		shouldStop:        make(chan bool),
		namesDB:           namesDB,
		handlers:          getHandlers(),
	}
}

func (c *botClient) Run() {
	rand.Seed(time.Now().UTC().UnixNano())
	var err error
	c.conn, err = net.Dial("tcp", c.config.server)
	if err != nil {
		log.Fatalf("ERROR - %v\n", err)
	}
	defer c.conn.Close()
	go func() {
		reader := bufio.NewReader(c.conn)
		for {
			serverResponse, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalf("Error while reading from connection - %v\n", err)
			}
			c.response <- serverResponse
		}
	}()
	c.registerUser()
	c.communicate()
}

func (c *botClient) communicate() {
	for {
		select {
		case <-c.shouldStop:
			fmt.Fprintf(c.conn, "QUIT :%v\n", "Bye!")
			log.Println("DEBUG - Sent QUIT command with message")
			return
		case serverResponse := <-c.response:
			message := extractResponse(strings.TrimRight(serverResponse, "\r\n"))
			log.Printf("DEBUG - %+v", message)

			c.handleCommand(message)
		}
	}
}

func (c *botClient) handleCommand(message *ircMessage) {
	if handle, found := c.handlers[message.command]; found {
		handle(c, message)
	}
}

func (c *botClient) Stop() {
	c.shouldStop <- true
}

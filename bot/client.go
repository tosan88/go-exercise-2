package bot

import (
	"fmt"
	"github.com/tosan88/go-exercise-2/conn"
	"github.com/tosan88/go-exercise-2/irc"
	"github.com/tosan88/go-exercise-2/storage"
	"log"
	"math/rand"
	"strings"
	"time"
)

type Conf struct {
	Server  string
	Channel string
	BotName string
}

type handler func(*Client, *irc.Message)

type Client struct {
	config            *Conf
	registeredBotName string
	conn              conn.Conn
	response          chan string
	shouldStop        chan bool
	db                storage.Persister
	handlers          map[string]handler
}

func New(config *Conf, db storage.DB) *Client {
	return &Client{
		config:            config,
		registeredBotName: config.BotName,
		response:          make(chan string),
		shouldStop:        make(chan bool),
		db:                db,
		handlers:          getHandlers(),
	}
}

func (c *Client) Run() {
	rand.Seed(time.Now().UTC().UnixNano())
	c.conn = conn.New(c.config.Server)
	defer c.conn.Close()
	go c.conn.ReadContinuously(c.response)
	c.registerUser()
	c.communicate()
}

func (c *Client) communicate() {
	for {
		select {
		case <-c.shouldStop:
			c.conn.Send(fmt.Sprint("QUIT :Bye!\n"))
			log.Println("DEBUG - Sent QUIT command with message")
			return
		case serverResponse := <-c.response:
			message := irc.ExtractResponse(strings.TrimRight(serverResponse, "\r\n"))
			log.Printf("DEBUG - %+v", message)

			c.handleCommand(message)
		}
	}
}

func (c *Client) handleCommand(message *irc.Message) {
	if handle, found := c.handlers[message.Command]; found {
		handle(c, message)
	}
}

func (c *Client) Stop() {
	c.shouldStop <- true
}

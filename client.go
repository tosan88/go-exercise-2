package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"strings"
)

type botClient struct {
	config            conf
	registeredBotName string
	conn              net.Conn
	response          chan string
	shouldStop        chan bool
}

func (c *botClient) Run() {
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
			if err == io.EOF {
				break
			}
			if err != nil {
				//could be this really fatal?
				log.Fatalf("%v\n", err)
			}
			c.response <- serverResponse
		}
	}()
	c.registerUser()
	c.communicate()
}

//https://tools.ietf.org/html/rfc2812#section-3.1
func (c *botClient) registerUser() {
	fmt.Fprintf(c.conn, "NICK %v\n", c.config.botName)
	fmt.Fprintf(c.conn, "USER %v 8 * :Greeting Bot Written in GoLang\n", c.config.botName)

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

			c.handleMessageCommand(message)
		default:
			//do nothing
		}
	}
}

//https://www.alien.net.au/irc/irc2numerics.html
func (c *botClient) handleMessageCommand(message *irc) {
	switch message.command {
	case "001":
		fallthrough
	case "RPL_WELCOME":
		log.Printf("Successfully joined to server %v\n", c.config.server)
	case "376":
		fallthrough
	case "RPL_ENDOFMOTD":
		fmt.Fprintf(c.conn, "JOIN #%v\n", c.config.channel)
	case "PING":
		log.Println("Sending PONG response")
		fmt.Fprintf(c.conn, "PONG :%s\n", message.message)
	case "JOIN":
		if strings.HasPrefix(message.initiator, c.registeredBotName+"!") {
			log.Printf("Successfully joined to channel #%v as %v\n", c.config.channel, c.registeredBotName)
			break
		}
		fmt.Fprintf(c.conn, "PRIVMSG #%v :Welcome in this channel %v\n", c.config.channel, strings.Split(message.initiator, "!")[0])
	case "432":
		fallthrough
	case "433":
		fallthrough
	case "436":
		fallthrough
	case "ERR_NICKCOLLISION":
		fallthrough
	case "ERR_ERRONEUSNICKNAME":
		fallthrough
	case "ERR_NICKNAMEINUSE":
		suffix := rand.Intn(1000)
		c.registeredBotName = fmt.Sprintf("%v%v", c.config.botName, suffix)
		log.Printf("Bot name could not be used. Adding suffix '%v' and retrying as %v\n", suffix, c.registeredBotName)
		fmt.Fprintf(c.conn, "NICK %v\n", c.registeredBotName)
	default:
		//do nothing
	}
}

func (c *botClient) Stop() {
	c.shouldStop <- true
}

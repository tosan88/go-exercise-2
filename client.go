package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
)

type user struct {
	available bool
	lastSeen time.Time
}

type botClient struct {
	config            *conf
	registeredBotName string
	conn              net.Conn
	response          chan string
	shouldStop        chan bool
	names             map[string]*user
}

func NewClient(config *conf) *botClient {
	return &botClient{
		config:            config,
		registeredBotName: config.botName,
		response:          make(chan string),
		shouldStop:        make(chan bool),
		names:             make(map[string]*user),
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
			if err == io.EOF {
				log.Fatalf("EOF - %v\n", err) //TODO might never appear
				break
			}
			if err != nil {
				log.Fatalf("NOT EOF - %v\n", err) //TODO handle differently?
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
	fmt.Fprintf(c.conn, "USER %v 8 * :Multifunctional Bot Written in GoLang\n", c.config.botName)

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
		}
	}
}

//https://www.alien.net.au/irc/irc2numerics.html
func (c *botClient) handleMessageCommand(message *ircMessage) {
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
		initiator := strings.Split(message.initiator, "!")[0]
		c.names[initiator] = &user{available:true}
		fmt.Fprintf(c.conn, "PRIVMSG #%v :Welcome in this channel %v\n", c.config.channel, initiator)
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
	case "PRIVMSG":
		//TODO polish things
		initiator := strings.Split(message.initiator, "!")[0]
		if initiator != c.registeredBotName{
			if strings.HasPrefix(message.message, fmt.Sprintf("#%v :!seen ", c.config.channel)) {
				nick := strings.Split(message.message, " ")[2]
				usr, found := c.names[nick]
				if !found {
					fmt.Fprintf(c.conn, "PRIVMSG #%v :Hello %v, unfortunately there are no records about %v\n", c.config.channel, initiator, nick)
					break
				}
				if usr.available {
					fmt.Fprintf(c.conn, "PRIVMSG #%v :Hello %v, %v is still present in this channel\n", c.config.channel, initiator, nick)
				} else {
					fmt.Fprintf(c.conn, "PRIVMSG #%v :Hello %v, %v was last seen on %v at %v\n", c.config.channel, initiator, nick, c.config.channel, usr.lastSeen)
				}
				break
			}
		}
		if strings.HasPrefix(message.message, c.registeredBotName+" ") && rand.Intn(100)%3 == 0 {
			fmt.Fprintf(c.conn, "PRIVMSG %v :Hello %v, I'm afraid I can't understand you, I'm just a bot...\n", initiator, initiator)
			break
		}
		if strings.Contains(message.message, c.registeredBotName) && rand.Intn(100)%3 == 0 {
			fmt.Fprintf(c.conn, "PRIVMSG #%v :Hello %v, would you like to tell me some cat fats?\n", c.config.channel, initiator)
			break
		}
		if rand.Intn(100)%2 == 0  && initiator != c.registeredBotName{
			var randomName string
			idx := rand.Intn(len(c.names))
			counter := 0
			for name := range c.names {
				if counter == idx {
					randomName = name
					break
				}
				counter++
			}
			fmt.Fprintf(c.conn, "PRIVMSG #%v :Check this out, %v - %v\n", c.config.channel, randomName, randomText[rand.Intn(len(randomText))])
			break
		}
	case "353": //RPL_NAMREPLY
		split := strings.Split(message.message, ":")
		if len(split) == 2 {
			names := strings.Split(split[1], " ")
			for _, name := range names {
				c.names[strings.TrimPrefix(name,"@")] = &user{available:true}
			}
			log.Printf("DEBUG - Names: %v\n", c.names)
		}
	case "KICK":
		initiator := strings.Split(message.initiator, "!")[0]
		fmt.Fprintf(c.conn, "PRIVMSG %v :That was rude!\n", initiator)
	case "PART":
		fallthrough
	case "QUIT":
		initiator := strings.Split(message.initiator, "!")[0]
		usr, found := c.names[initiator]
		if !found {
			log.Printf("WARN - Logical error. %v quit who was not in the channel :-/\n", initiator)
			break
		}
		usr.available = false
		usr.lastSeen = time.Now()
	default:
		//do nothing
	}
}

func (c *botClient) Stop() {
	c.shouldStop <- true
}

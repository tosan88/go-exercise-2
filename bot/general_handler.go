package bot

import (
	"database/sql"
	"fmt"
	. "github.com/tosan88/go-exercise-2/irc"
	"github.com/tosan88/go-exercise-2/storage"
	"log"
	"math/rand"
	"strings"
	"time"
)

func getHandlers() map[string]handler {
	return map[string]handler{
		RPL_WELCOME:              (*Client).handleWelcomeReply,
		RPL_ENDOFMOTD:            (*Client).handleEndMessageOfTheDayCommand,
		PING:                     (*Client).handlePingCommand,
		JOIN:                     (*Client).handleJoinCommand,
		RPL_NAMREPLY:             (*Client).handleNamesReply,
		MESSAGE_CMD:              (*Client).handleUserMessageCommand,
		KICK:                     (*Client).handleKickCommand,
		QUIT:                     (*Client).handleUserDeparture,
		PART:                     (*Client).handleUserDeparture,
		ERR_NICKNAMEINUSE:        (*Client).handleNickNotGoodCommand,
		ERR_NICKNAMEINUSE_NUM:    (*Client).handleNickNotGoodCommand,
		ERR_ERRONEUSNICKNAME:     (*Client).handleNickNotGoodCommand,
		ERR_ERRONEUSNICKNAME_NUM: (*Client).handleNickNotGoodCommand,
		ERR_NICKCOLLISION:        (*Client).handleNickNotGoodCommand,
		ERR_NICKCOLLISION_NUM:    (*Client).handleNickNotGoodCommand,
	}
}

//https://tools.ietf.org/html/rfc2812#section-3.1
func (c *Client) registerUser() {
	c.conn.Send(fmt.Sprintf("NICK %v\n", c.config.BotName))
	c.conn.Send(fmt.Sprintf("USER %v 8 * :Multifunctional Bot Written in GoLang\n", c.config.BotName))
}

func (c *Client) handleWelcomeReply(message *Message) {
	log.Printf("Successfully joined to server %v\n", c.config.Server)
}

func (c *Client) handlePingCommand(message *Message) {
	log.Println("Sending PONG response")
	c.conn.Send(fmt.Sprintf("PONG :%s\n", message.Parameter))
}

func (c *Client) handleEndMessageOfTheDayCommand(message *Message) {
	c.conn.Send(fmt.Sprintf("JOIN #%v\n", c.config.Channel))
}

func (c *Client) handleJoinCommand(message *Message) {
	if strings.HasPrefix(message.Prefix, c.registeredBotName+"!") {
		log.Printf("Successfully joined to channel #%v as %v\n", c.config.Channel, c.registeredBotName)
		return
	}
	initiator := strings.Split(message.Prefix, "!")[0]

	_, err := c.db.GetUser(initiator)
	if err == nil || err == sql.ErrNoRows {
		err = c.db.UpdateUser(&storage.User{Name: initiator, Available: true, LastSeen: time.Now()})
	} else {
		log.Printf("WARN - get user: %v\n", err)
		err = c.db.AddUser(&storage.User{Available: true, LastSeen: time.Now(), Name: initiator})
	}
	if err != nil {
		log.Printf("ERROR - insert into names: %v\n", err)
		return
	}
	c.conn.Send(fmt.Sprintf("PRIVMSG #%v :Welcome in this channel %v\n", c.config.Channel, initiator))
}

func (c *Client) handleNickNotGoodCommand(message *Message) {
	suffix := rand.Intn(1000)
	c.registeredBotName = fmt.Sprintf("%v%v", c.config.BotName, suffix)
	log.Printf("Bot name could not be used. Adding suffix '%v' and retrying as %v\n", suffix, c.registeredBotName)
	c.conn.Send(fmt.Sprintf("NICK %v\n", c.registeredBotName))
}

func (c *Client) handleNamesReply(message *Message) {
	split := strings.Split(message.Parameter, ":")
	if len(split) == 2 {
		names := strings.Split(split[1], " ")
		for _, name := range names {
			usrName := strings.TrimPrefix(name, "@")
			usr := &storage.User{Available: true, Name: usrName, LastSeen: time.Now()}

			storedUsr, err := c.db.GetUser(usrName)
			if err != nil && err != sql.ErrNoRows {
				log.Printf("ERROR - getting user %v: %v\n", usrName, err)
				continue
			}
			if storedUsr.Name == "" {
				c.db.AddUser(usr)
			}
		}
	}
}

func (c *Client) handleKickCommand(message *Message) {
	initiator := strings.Split(message.Prefix, "!")[0]
	kicked := strings.Split(message.Parameter, " ")[1]
	_, err := c.db.GetUser(kicked)
	if err != nil {
		log.Printf("ERROR - No stored user %v: %v\n", kicked, err)
		return
	}
	c.db.UpdateUser(&storage.User{Name: kicked, Available: false, LastSeen: time.Now()})
	c.conn.Send(fmt.Sprintf("PRIVMSG %v :That was rude!\n", initiator))
}

func (c *Client) handleUserDeparture(message *Message) {
	initiator := strings.Split(message.Prefix, "!")[0]
	_, err := c.db.GetUser(initiator)
	if err != nil {
		log.Printf("ERROR - No stored user %v: %v\n", initiator, err)
		return
	}
	c.db.UpdateUser(&storage.User{Name: initiator, Available: false, LastSeen: time.Now()})
}

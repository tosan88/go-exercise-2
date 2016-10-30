package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
)

func getHandlers() map[string]handler {
	return map[string]handler{
		RPL_WELCOME:              (*botClient).handleWelcomeReply,
		RPL_ENDOFMOTD:            (*botClient).handleEndMessageOfTheDayCommand,
		PING:                     (*botClient).handlePingCommand,
		JOIN:                     (*botClient).handleJoinCommand,
		RPL_NAMREPLY:             (*botClient).handleNamesReply,
		MESSAGE_CMD:              (*botClient).handleUserMessageCommand,
		KICK:                     (*botClient).handleKickCommand,
		QUIT:                     (*botClient).handleUserDeparture,
		PART:                     (*botClient).handleUserDeparture,
		ERR_NICKNAMEINUSE:        (*botClient).handleNickNotGoodCommand,
		ERR_NICKNAMEINUSE_NUM:    (*botClient).handleNickNotGoodCommand,
		ERR_ERRONEUSNICKNAME:     (*botClient).handleNickNotGoodCommand,
		ERR_ERRONEUSNICKNAME_NUM: (*botClient).handleNickNotGoodCommand,
		ERR_NICKCOLLISION:        (*botClient).handleNickNotGoodCommand,
		ERR_NICKCOLLISION_NUM:    (*botClient).handleNickNotGoodCommand,
	}
}

//https://tools.ietf.org/html/rfc2812#section-3.1
func (c *botClient) registerUser() {
	fmt.Fprintf(c.conn, "NICK %v\n", c.config.botName)
	fmt.Fprintf(c.conn, "USER %v 8 * :Multifunctional Bot Written in GoLang\n", c.config.botName)
}

func (c *botClient) handleWelcomeReply(message *ircMessage) {
	log.Printf("Successfully joined to server %v\n", c.config.server)
}

func (c *botClient) handlePingCommand(message *ircMessage) {
	log.Println("Sending PONG response")
	fmt.Fprintf(c.conn, "PONG :%s\n", message.message)
}

func (c *botClient) handleEndMessageOfTheDayCommand(message *ircMessage) {
	fmt.Fprintf(c.conn, "JOIN #%v\n", c.config.channel)
}

func (c *botClient) handleJoinCommand(message *ircMessage) {
	if strings.HasPrefix(message.initiator, c.registeredBotName+"!") {
		log.Printf("Successfully joined to channel #%v as %v\n", c.config.channel, c.registeredBotName)
		return
	}
	initiator := strings.Split(message.initiator, "!")[0]
	c.names[initiator] = &user{available: true}
	fmt.Fprintf(c.conn, "PRIVMSG #%v :Welcome in this channel %v\n", c.config.channel, initiator)
}

func (c *botClient) handleNickNotGoodCommand(message *ircMessage) {
	suffix := rand.Intn(1000)
	c.registeredBotName = fmt.Sprintf("%v%v", c.config.botName, suffix)
	log.Printf("Bot name could not be used. Adding suffix '%v' and retrying as %v\n", suffix, c.registeredBotName)
	fmt.Fprintf(c.conn, "NICK %v\n", c.registeredBotName)
}

func (c *botClient) handleUserMessageCommand(message *ircMessage) {
	initiator := strings.Split(message.initiator, "!")[0]
	if initiator == c.registeredBotName {
		return
	}
	if c.isUserCommand(message, "seen") {
		c.handleNotSeenUserCommand(initiator, message)
	}
	if c.isUserCommand(message, "cat fact") {
		c.handleCatFactUserCommand(message)
	}
	if strings.HasPrefix(message.message, c.registeredBotName+" ") {
		fmt.Fprintf(c.conn, "PRIVMSG %v :Hello %v, I'm afraid I can't understand you, I'm just a bot...\n", initiator, initiator)
		return
	}
}

func (c *botClient) isUserCommand(message *ircMessage, cmd string) bool {
	return strings.HasPrefix(message.message, fmt.Sprintf("#%v :!%v", c.config.channel, cmd)) ||
		strings.HasPrefix(message.message, fmt.Sprintf("%v :!%v", c.registeredBotName, cmd))
}

func (c *botClient) handleNotSeenUserCommand(initiator string, message *ircMessage) {
	splitCommand := strings.Split(message.message, " ")
	location := splitCommand[0]
	if location != fmt.Sprintf("#%v", c.config.channel) {
		location = initiator
	}
	if len(splitCommand) < 3 {
		fmt.Fprintf(c.conn, "PRIVMSG %v :Hello %v, not enough arguments for !seen command. "+
			"Please use in the form \"!seen <nickname>\"\n", location, initiator)
		return
	}

	nick := splitCommand[2]
	usr, found := c.names[nick]
	if !found {
		fmt.Fprintf(c.conn, "PRIVMSG %v :Hello %v, unfortunately there are no records about %v\n", location, initiator, nick)
		return
	}

	if usr.available {
		fmt.Fprintf(c.conn, "PRIVMSG %v :Hello %v, %v is still present in this channel\n", location, initiator, nick)
	} else {
		fmt.Fprintf(c.conn, "PRIVMSG %v :Hello %v, %v was last seen on %v at %v\n", location, initiator, nick, c.config.channel, usr.lastSeen)
	}
}

func (c *botClient) handleCatFactUserCommand(message *ircMessage) {
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

}
func (c *botClient) handleNamesReply(message *ircMessage) {
	split := strings.Split(message.message, ":")
	if len(split) == 2 {
		names := strings.Split(split[1], " ")
		for _, name := range names {
			c.names[strings.TrimPrefix(name, "@")] = &user{available: true}
		}
	}
}

func (c *botClient) handleKickCommand(message *ircMessage) {
	initiator := strings.Split(message.initiator, "!")[0]
	fmt.Fprintf(c.conn, "PRIVMSG %v :That was rude!\n", initiator)
}

func (c *botClient) handleUserDeparture(message *ircMessage) {
	initiator := strings.Split(message.initiator, "!")[0]
	usr, found := c.names[initiator]
	if !found {
		log.Printf("WARN - Logical error. %v quit but was not in the channel :-/\n", initiator)
		return
	}
	usr.available = false
	usr.lastSeen = time.Now()
}

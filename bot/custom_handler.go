package bot

import (
	"database/sql"
	"fmt"
	"github.com/tosan88/go-exercise-2/irc"
	"log"
	"regexp"
	"strings"
)

var urlRegex = regexp.MustCompile("^.*(https?:\\/\\/[^ ]+\\.[^ ]+).*$")

func (c *Client) handleUserMessageCommand(message *irc.Message) {
	initiator := strings.Split(message.Prefix, "!")[0]
	if initiator == c.registeredBotName {
		return
	}
	if c.isUserCommand(message, "seen") {
		c.handleNotSeenUserCommand(initiator, message)
		return
	}
	if c.isUserCommand(message, "cat fact") {
		c.handleCatFactUserCommand(initiator, message)
		return
	}

	if matchForUrl(message.Parameter) != "" {
		//TODO get page title
	}

	if strings.HasPrefix(message.Parameter, c.registeredBotName+" ") {
		c.conn.Send(fmt.Sprintf("PRIVMSG %v :Hello %v, I'm afraid I can't understand you, I'm just a bot...\n", initiator, initiator))
		return
	}
}

func matchForUrl(message string) string {
	matches := urlRegex.FindStringSubmatch(message)
	if len(matches) == 2 {
		return matches[1]
	}
	return ""
}

func (c *Client) isUserCommand(message *irc.Message, cmd string) bool {
	return strings.HasPrefix(message.Parameter, fmt.Sprintf("#%v :!%v", c.config.Channel, cmd)) ||
		strings.HasPrefix(message.Parameter, fmt.Sprintf("%v :!%v", c.registeredBotName, cmd))
}

func (c *Client) handleNotSeenUserCommand(initiator string, message *irc.Message) {
	splitCommand := strings.Split(message.Parameter, " ")
	location := splitCommand[0]
	if location != fmt.Sprintf("#%v", c.config.Channel) {
		location = initiator
	}
	if len(splitCommand) < 3 {
		c.conn.Send(fmt.Sprintf("PRIVMSG %v :Hello %v, not enough arguments for !seen command. "+
			"Please use in the form \"!seen <nickname>\"\n", location, initiator))
		return
	}

	nick := splitCommand[2]
	usr, err := c.db.GetUser(nick)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("ERROR - getting user %v: %v\n", nick, err)
	}

	if usr.Name == "" {
		c.conn.Send(fmt.Sprintf("PRIVMSG %v :Hello %v, unfortunately there are no records about %v\n", location, initiator, nick))
		return
	}

	if usr.Available {
		c.conn.Send(fmt.Sprintf("PRIVMSG %v :Hello %v, %v is still present in %v channel\n", location, initiator, nick, c.config.Channel))
	} else {
		c.conn.Send(fmt.Sprintf("PRIVMSG %v :Hello %v, %v was last seen on %v channel at %v\n", location, initiator, nick, c.config.Channel, usr.LastSeen))
	}
}

func (c *Client) handleCatFactUserCommand(initiator string, message *irc.Message) {
	splitCommand := strings.Split(message.Parameter, " ")
	location := splitCommand[0]
	if location != fmt.Sprintf("#%v", c.config.Channel) {
		location = initiator
	}
	c.conn.Send(fmt.Sprintf("PRIVMSG %v :%v\n", location, GetCatFact()))

}

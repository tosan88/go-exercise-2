package bot

import (
	"database/sql"
	"fmt"
	"github.com/tosan88/go-exercise-2/irc"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
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

	if url := matchForUrl(message.Parameter); url != "" {
		c.handleUrl(url, initiator, message)
		return
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

func (c *Client) handleUrl(url string, initiator string, message *irc.Message) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("ERROR - Cannot create url %v: %v\n", url, err)
		return
	}
	req.Header.Add("Accept", "text/html")
	resp, err := c.hc.Do(req)
	if err != nil {
		log.Printf("ERROR - Requesting %v: %v\n", url, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("ERROR - Request to %v returned status: %v\n", url, resp.StatusCode)
		return
	}

	title, err := getTitle(resp.Body)
	if err != nil {
		log.Printf("WARN - Cannot get title for %v: %v\n", url, err)
		return
	}

	splitCommand := strings.Split(message.Parameter, " ")
	location := splitCommand[0]
	if location != fmt.Sprintf("#%v", c.config.Channel) {
		location = initiator
	}

	if title != "" {
		c.conn.Send(fmt.Sprintf("PRIVMSG %v :The title for %v is: %v\n", location, url, title))
	} else {
		log.Printf("WARN - No title for %v\n", url)
	}

}

func getTitle(r io.Reader) (string, error) {
	tokenizer := html.NewTokenizer(r)
	for {
		tokenType := tokenizer.Next()
		token := tokenizer.Token()

		switch tokenType {
		case html.ErrorToken:
			err := tokenizer.Err()
			if err == io.EOF {
				return "", nil
			}
			return "", err
		case html.StartTagToken:
			if token.Data == "title" {
				if tokenizer.Next() == html.TextToken {
					return tokenizer.Token().String(), nil
				}
			}
		}
	}
}

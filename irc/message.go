package irc

import (
	"regexp"
)

//https://www.alien.net.au/irc/irc2numerics.html
const (
	RPL_WELCOME              = "001"
	RPL_ENDOFMOTD            = "376"
	PING                     = "PING"
	JOIN                     = "JOIN"
	MESSAGE_CMD              = "PRIVMSG"
	RPL_NAMREPLY             = "353"
	KICK                     = "KICK"
	PART                     = "PART"
	QUIT                     = "QUIT"
	ERR_ERRONEUSNICKNAME     = "ERR_ERRONEUSNICKNAME"
	ERR_ERRONEUSNICKNAME_NUM = "432"
	ERR_NICKNAMEINUSE        = "ERR_NICKNAMEINUSE"
	ERR_NICKNAMEINUSE_NUM    = "433"
	ERR_NICKCOLLISION        = "ERR_NICKCOLLISION"
	ERR_NICKCOLLISION_NUM    = "436"
)

type Message struct {
	Prefix    string
	Command   string
	Parameter string
}

var normalReply = regexp.MustCompile(`^:([^ ]+) ([^ ]+) (.*)`)
var statusReply = regexp.MustCompile(`^([^:][^ ]+) :(.*)`)

func ExtractResponse(response string) (msg *Message) {
	msg, extracted := extractNormalReply(response)
	if extracted {
		return
	}
	msg, extracted = extractStatusReply(response)
	if extracted {
		return
	}

	//we couldn't parse it, use whole response as the message
	msg = &Message{
		Parameter: response,
	}
	return
}

func extractNormalReply(response string) (msg *Message, extracted bool) {
	matches := normalReply.FindStringSubmatch(response)
	if len(matches) == 4 {
		msg = &Message{
			Prefix:    matches[1],
			Command:   matches[2],
			Parameter: matches[3],
		}
		extracted = true
	}
	return
}
func extractStatusReply(response string) (msg *Message, extracted bool) {
	matches := statusReply.FindStringSubmatch(response)
	if len(matches) == 3 {
		msg = &Message{
			Command:   matches[1],
			Parameter: matches[2],
		}
		extracted = true
	}
	return
}

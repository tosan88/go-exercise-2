package main

import (
	"regexp"
)

type ircMessage struct {
	initiator string
	command   string
	message   string
}

var normalReply = regexp.MustCompile(`^:([^ ]+) ([^ ]+) (.*)`)
var statusReply = regexp.MustCompile(`^([^:][^ ]+) :(.*)`)

func extractResponse(response string) (msg *ircMessage) {
	msg, extracted := extractNormalReply(response)
	if extracted {
		return
	}
	msg, extracted = extractStatusReply(response)
	if extracted {
		return
	}

	//we couldn't parse it, use whole response as the message
	msg = &ircMessage{
		message: response,
	}
	return
}

func extractNormalReply(response string) (msg *ircMessage, extracted bool) {
	matches := normalReply.FindStringSubmatch(response)
	if len(matches) == 4 {
		msg = &ircMessage{
			initiator: matches[1],
			command:   matches[2],
			message:   matches[3],
		}
		extracted = true
	}
	return
}
func extractStatusReply(response string) (msg *ircMessage, extracted bool) {
	matches := statusReply.FindStringSubmatch(response)
	if len(matches) == 3 {
		msg = &ircMessage{
			command: matches[1],
			message: matches[2],
		}
		extracted = true
	}
	return
}

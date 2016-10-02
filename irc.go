package main

import (
	"regexp"
)

type irc struct {
	initiator string
	command   string
	message   string
}

var normalReply = regexp.MustCompile(`^:([^ ]+) ([^ ]+) (.*)`)
var statusReply = regexp.MustCompile(`^([^:][^ ]+) :(.*)`)

func extractResponse(response string) (message *irc) {
	message, extracted := extractNormalReply(response)
	if extracted {
		return
	}
	message, extracted = extractStatusReply(response)
	if extracted {
		return
	}

	//we couldn't parse it, use whole response as the message
	message = &irc{
		message: response,
	}
	return
}

func extractNormalReply(response string) (message *irc, extracted bool) {
	matches := normalReply.FindStringSubmatch(response)
	if len(matches) == 4 {
		message = &irc{
			initiator: matches[1],
			command:   matches[2],
			message:   matches[3],
		}
		extracted = true
	}
	return
}
func extractStatusReply(response string) (message *irc, extracted bool) {
	matches := statusReply.FindStringSubmatch(response)
	if len(matches) == 3 {
		message = &irc{
			command: matches[1],
			message: matches[2],
		}
		extracted = true
	}
	return
}

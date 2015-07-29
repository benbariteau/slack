package rtm

import (
	"regexp"
	"strings"

	"github.com/firba1/slack"

	"github.com/gorilla/websocket"
)

type Conn struct {
	conn           *websocket.Conn
	messageCounter int
	// maps user ID string to users.
	Users map[string]slack.User
}

func (c *Conn) Close() error {
	return c.conn.Close()
}

type Event struct {
	Type string `json:"type"`
}

type Message struct {
	Event
	Subtype   string `json:"subtype"`
	Channel   string `json:"channel"`
	User      string `json:"user"`
	Text      string `json:"text"`
	Timestamp string `json:"ts"`
}

/*
ReadChannel returns a channel that reads and sends the Messages coming from the websocket
*/
func (c *Conn) MessageChan() <-chan Message {
	ch := make(chan Message, 10)
	go func() {
		for {
			message := Message{}
			c.conn.ReadJSON(&message)
			if message.Type != "message" {
				continue
			}
			ch <- message
		}
	}()
	return ch
}

type sendMessage struct {
	ID      int    `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func (c *Conn) SendMessage(text, channel string) error {
	c.messageCounter++
	msg := sendMessage{
		ID:      c.messageCounter,
		Type:    "message",
		Channel: channel,
		Text:    text,
	}
	return c.conn.WriteJSON(msg)
}

var escapeRegex = regexp.MustCompile("<(.*?)>")

var escapeTypePostprocessors = map[int]func(string) string{
	userEscape:    func(s string) string { return "@" + s },
	channelEscape: func(s string) string { return "#" + s },
}

/*
UnescapeMessage takes in the escape string text of a message and returns a new string that appears as it would to a user.

UnescapeMessage does so by parsing escape sequences according to <https://api.slack.com/docs/formatting> and substituting the appropriate user-facing junk (e.g. <@UABC123> would become @firba1, assuming there's a user named firba1 with the user ID UABC123).
*/
func (c Conn) UnescapeMessage(message string) string {
	message = escapeRegex.ReplaceAllStringFunc(message, func(match string) string {
		unescapedMatch, escapeType := replaceEscapeHelper(c, match)
		postprocess := escapeTypePostprocessors[escapeType]
		if postprocess != nil {
			unescapedMatch = postprocess(unescapedMatch)
		}
		return unescapedMatch
	})

	// finally replace all html entity escapes
	message = strings.Replace(message, "&amp;", "&", -1)
	message = strings.Replace(message, "&lt;", "<", -1)
	message = strings.Replace(message, "&gt;", ">", -1)
	return message
}

func replaceEscapeHelper(c Conn, match string) (unescape string, escapeType int) {
	// remove < and > from each end
	fullEscape := match[1 : len(match)-1]

	// check for display string
	escapeParts := strings.Split(fullEscape, "|")

	// this is a case we don't recognize, just return the original match and treat it as a link (the default)
	if len(escapeParts) > 2 || len(escapeParts) <= 0 {
		unescape = match
		escapeType = linkEscape
		return
	}

	escape := escapeParts[0]
	escapeType = parseEscapeType(escape)

	// if we have an alias, just return that
	if len(escapeParts) == 2 {
		unescape = escapeParts[1]
		return
	}

	// since there's no alias, now it's time for idenitifier lookup
	escapeType = parseEscapeType(escape)

	switch escapeType {
	case userEscape:
		// user link
		user := c.Users[escape[1:]]
		// if user is zero value, this will just be empty string, which we handle later
		unescape = user.Name
	}

	// if we couldn't unescape properly, just return the original match text, make it a linkEscape type to prevent post processing
	if unescape == "" {
		escapeType = linkEscape
		unescape = match
	}
	return
}

const (
	linkEscape = iota
	userEscape
	channelEscape
	commandEscape
)

/*
parseEscapeType is a convience function for getting an easily comparable type from an escape sequence (e.g. "@U123A56BC" for users "#C789D10EF" for channels, etc)
*/
func parseEscapeType(escapeString string) int {
	switch {
	case escapeString[0:2] == "@U":
		return userEscape
	case escapeString[0:2] == "#C":
		return channelEscape
	case escapeString[0] == '!':
		return commandEscape
	default:
		// as per the docs, anything we can't recognize like this is a link
		return linkEscape
	}
}

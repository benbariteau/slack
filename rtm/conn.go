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

/*
UnescapeMessage takes in the escape string text of a message and returns a new string that appears as it would to a user.

UnescapeMessage does so by parsing escape sequences according to <https://api.slack.com/docs/formatting> and substituting the appropriate user-facing junk (e.g. <@UABC123> would become @firba1, assuming there's a user named firba1 with the user ID UABC123).
*/
func (c Conn) UnescapeMessage(message string) string {
	message = escapeRegex.ReplaceAllStringFunc(message, func(match string) string {
		// remove < and > from each end
		escape := match[1 : len(match)-1]
		escapeParts := strings.Split(escape, "|")

		// this is a case we don't recognize, just return the original match
		if len(escapeParts) > 2 {
			return match
		}

		if len(escapeParts) == 2 {
			return escapeParts[1]
		}

		// since there's no alias, now it's time for idenitifier lookup
		escape = escapeParts[0]

		switch {
		case escape[0:2] == "@U":
			// user link
			user, ok := c.Users[escape[1:]]
			if ok {
				// mentions always have an "@" preceding the username
				return "@" + user.Name
			}
		}
		// if we were unable to unescape this, just return the original match
		return match
	})

	// finally replace all html entity escapes
	message = strings.Replace(message, "&amp;", "&", -1)
	message = strings.Replace(message, "&lt;", "<", -1)
	message = strings.Replace(message, "&gt;", ">", -1)
	return message
}

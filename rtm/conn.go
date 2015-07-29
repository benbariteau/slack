package rtm

import (
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

/*
UnescapeMessage takes in the escape string text of a message and returns a new string that appears as it would to a user.

UnescapeMessage does so by parsing escape sequences according to <https://api.slack.com/docs/formatting> and substituting the appropriate user-facing junk (e.g. <@UABC123> would become @firba1, assuming there's a user named firba1 with the user ID UABC123).
*/
func (c Conn) UnescapeMessage(message string) string {
	// first replace all html entity escapes
	message = strings.Replace(message, "&amp;", "&", -1)
	message = strings.Replace(message, "&lt;", "<", -1)
	message = strings.Replace(message, "&gt;", ">", -1)
	return message
}

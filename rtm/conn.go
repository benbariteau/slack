package rtm

import (
	"github.com/gorilla/websocket"
)

type Conn struct {
	conn *websocket.Conn
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

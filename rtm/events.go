package rtm

import (
	"fmt"
)

type Event interface {
	Type() string
}

type BasicEvent struct {
	TypeStr  string
	RawEvent map[string]interface{}
}

func (e BasicEvent) Type() string {
	return e.TypeStr
}

type Hello struct{}

func (h Hello) Type() string { return "hello" }

type Message interface {
	Event
	Subtype() string
	Text() string
	User() string
	Channel() string
	Timestamp() string
}

type BasicMessage struct {
	channel   string
	user      string
	text      string
	timestamp string
	subtype   string
}

func (m BasicMessage) Type() string    { return "message" }
func (m BasicMessage) Subtype() string { return m.subtype }
func (m BasicMessage) Text() string {
	return m.text
}
func (m BasicMessage) User() string {
	return m.user
}
func (m BasicMessage) Channel() string {
	return m.channel
}
func (m BasicMessage) Timestamp() string {
	return m.timestamp
}

type ChannelCreated struct {
	ID      string
	Name    string
	Created int
	Creator string
}

func (c ChannelCreated) Type() string { return "channel_created" }

type ChannelRename struct {
	ID      string
	Name    string
	Created int
}

func (c ChannelRename) Type() string { return "channel_rename" }

func toEvent(rawEvent map[string]interface{}) Event {
	fmt.Println(rawEvent)
	rawEventType, ok := rawEvent["type"]
	if !ok {
		return BasicEvent{"invalid", rawEvent}
	}
	switch eventType := rawEventType.(string); eventType {
	case "hello":
		return Hello{}
	case "message":
		return BasicMessage{
			channel:   getStringField(rawEvent, "channel"),
			user:      getStringField(rawEvent, "user"),
			text:      getStringField(rawEvent, "text"),
			timestamp: getStringField(rawEvent, "ts"),
			subtype:   getStringField(rawEvent, "subtype"),
		}
	case "channel_created":
		channelInfo := rawEvent["channel"].(map[string]interface{})
		return ChannelCreated{
			ID:      channelInfo["id"].(string),
			Name:    channelInfo["name"].(string),
			Created: int(channelInfo["created"].(float64)),
			Creator: channelInfo["creator"].(string),
		}
	case "channel_rename":
		channelInfo := rawEvent["channel"].(map[string]interface{})
		return ChannelRename{
			ID:      channelInfo["id"].(string),
			Name:    channelInfo["name"].(string),
			Created: int(channelInfo["created"].(float64)),
		}
	default:
		return BasicEvent{eventType, rawEvent}
	}
}

func getStringField(m map[string]interface{}, key string) string {
	val := m[key]
	switch val.(type) {
	case string:
		return val.(string)
	default:
		return ""
	}
}

package rtm

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
	switch eventType := rawEvent["type"].(string); eventType {
	case "hello":
		return Hello{}
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

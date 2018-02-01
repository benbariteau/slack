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
	userChanges    chan<- slack.User
	infoRequests   chan<- userInfoRequest
	cancel         chan struct{}
}

func (c *Conn) Close() error {
	close(c.cancel)
	return c.conn.Close()
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
NextEvent blocks until the next Event is sent and then returns it to the caller.
*/
func (c *Conn) NextEvent() Event {
	rawEvent := make(map[string]interface{})
	c.conn.ReadJSON(&rawEvent)
	return toEvent(rawEvent)
}

var escapeRegex = regexp.MustCompile("<(.*?)>")

var escapeTypePostprocessors = map[EscapeType]func(string) string{
	userEscape:    func(s string) string { return "@" + s },
	channelEscape: func(s string) string { return "#" + s },
}

/*
UnescapeMessage takes in the escape string text of a message and returns a new string that appears as it would to a user.

UnescapeMessage does so by parsing escape sequences according to <https://api.slack.com/docs/formatting> and substituting the appropriate user-facing junk (e.g. <@UABC123> would become @firba1, assuming there's a user named firba1 with the user ID UABC123).
*/
func (c Conn) UnescapeMessage(message string) string {
	return c.UnescapeMessagePostprocess(message, func(s string, i EscapeType) string { return s })
}

func (c Conn) UnescapeMessagePostprocess(
	message string,
	postprocessor func(userString string, escapeType EscapeType) string,
) string {
	message = escapeRegex.ReplaceAllStringFunc(message, func(match string) string {
		unescapedMatch, escapeType := replaceEscapeHelper(c, match)
		postprocess := escapeTypePostprocessors[escapeType]
		if postprocess != nil {
			unescapedMatch = postprocess(unescapedMatch)
		}
		return postprocessor(unescapedMatch, escapeType)
	})

	// finally replace all html entity escapes
	message = strings.Replace(message, "&amp;", "&", -1)
	message = strings.Replace(message, "&lt;", "<", -1)
	message = strings.Replace(message, "&gt;", ">", -1)
	return message
}

func replaceEscapeHelper(c Conn, match string) (unescape string, escapeType EscapeType) {
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
		user := c.UserInfo(escape[1:])
		// if user is zero value, this will just be empty string, which we handle later
		unescape = user.Profile.DisplayName
	}

	// if we couldn't unescape properly, just return the original match text, make it a linkEscape type to prevent post processing
	if unescape == "" {
		escapeType = linkEscape
		unescape = match
	}
	return
}

type EscapeType int

const (
	linkEscape EscapeType = iota
	userEscape
	channelEscape
	commandEscape
)

/*
parseEscapeType is a convience function for getting an easily comparable type from an escape sequence (e.g. "@U123A56BC" for users "#C789D10EF" for channels, etc)
*/
func parseEscapeType(escapeString string) EscapeType {
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

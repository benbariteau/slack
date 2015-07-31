package rtm

import (
	"testing"

	"github.com/firba1/slack"

	"github.com/firba1/util/assert"
	"github.com/gorilla/websocket"
)

func TestUnescapeMessage(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		// HTML entity replacements
		{"Simon &amp; Garfunkel", "Simon & Garfunkel"},
		{"Independence Day &gt; Stargate", "Independence Day > Stargate"},
		{"):&lt;", "):<"},
		// existing user
		{"<@U123A56BC>++", "@fart++"},
		// user with an alias
		{"<@U123A56BC|butt>++", "@butt++"},
		// unrecognized user
		{"<@UXYZ987WV>++", "<@UXYZ987WV>++"},
	}

	conn, err := dialer{
		rtmStartFunc: func(token string) (*websocket.Conn, slack.RTMStartInfo, error) {
			rtmStartInfo := slack.RTMStartInfo{
				Users: []slack.User{
					slack.User{
						ID:   "U123A56BC",
						Name: "fart",
					},
				},
			}
			return nil, rtmStartInfo, nil
		},
	}.Dial("")
	assert.NilError(t, err)

	for _, test := range tests {
		assert.Equal(t, conn.UnescapeMessage(test.in), test.out)
	}
}

func TestUserInfoWithUpdates(t *testing.T) {
	conn, err := dialer{
		rtmStartFunc: func(token string) (*websocket.Conn, slack.RTMStartInfo, error) {
			return nil, slack.RTMStartInfo{}, nil
		},
	}.Dial("")
	assert.NilError(t, err)

	// user should not exist (zero value)
	assert.Equal(t, conn.UserInfo("U123"), slack.User{})

	// send a user update
	user := slack.User{ID: "U123", Name: "butt fart"}
	conn.userChanges <- user

	// user should exist now
	assert.Equal(t, conn.UserInfo("U123"), user)
}

func TestParseEscapeType(t *testing.T) {
	tests := []struct {
		in  string
		out int
	}{
		{"@U123A56BC", userEscape},
		{"#C789D01EF", channelEscape},
		{"!everyone", commandEscape},
		{"http://www.facebook.com/prenis", linkEscape},
	}

	for _, test := range tests {
		assert.Equal(t, parseEscapeType(test.in), test.out)
	}
}

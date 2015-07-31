package rtm

import (
	"testing"

	"github.com/firba1/slack"

	"github.com/firba1/util/assert"
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

	conn := Conn{cancel: make(chan struct{})}
	users := []slack.User{
		slack.User{
			ID:   "U123A56BC",
			Name: "fart",
		},
	}
	conn.userChanges, conn.infoRequests = serveUserInfo(users, conn.cancel)

	for _, test := range tests {
		assert.Equal(t, conn.UnescapeMessage(test.in), test.out)
	}
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

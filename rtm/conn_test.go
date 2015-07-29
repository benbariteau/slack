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
		{"Simon &amp; Garfunkel", "Simon & Garfunkel"},
		{"Independence Day &gt; Stargate", "Independence Day > Stargate"},
		{"):&lt;", "):<"},
		{"<@U123A56BC>++", "@fart++"},
	}

	conn := Conn{
		Users: map[string]slack.User{
			"U123A56BC": slack.User{
				ID:   "U123A56BC",
				Name: "fart",
			},
		},
	}

	for _, test := range tests {
		assert.Equal(t, conn.UnescapeMessage(test.in), test.out)
	}
}

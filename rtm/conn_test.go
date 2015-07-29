package rtm

import (
	"testing"

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
	}

	conn := Conn{}

	for _, test := range tests {
		assert.Equal(t, conn.UnescapeMessage(test.in), test.out)
	}
}

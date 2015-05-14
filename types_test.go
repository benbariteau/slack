package slack

import (
	"testing"

	"github.com/firba1/util/assert"
)

func TestBasicResponseOk(t *testing.T) {
	testcases := []struct {
		in  BasicResponse
		out bool
	}{
		{BasicResponse{}, false},
		{BasicResponse{OK: false}, false},
		{BasicResponse{OK: true}, true},
	}

	for _, test := range testcases {
		assert.Equal(t, test.in.Ok(), test.out)
	}
}

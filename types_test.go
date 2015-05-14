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

func TestBasicResponseError(t *testing.T) {
	testcases := []struct {
		in  BasicResponse
		out error
	}{
		{BasicResponse{OK: true}, nil},
		{BasicResponse{OK: false, ErrorString: "too_many_farts"}, SlackError{"too_many_farts"}},
		{BasicResponse{OK: false}, SlackError{}},
	}

	for _, test := range testcases {
		assert.Equal(t, test.in.Error(), test.out)
	}
}

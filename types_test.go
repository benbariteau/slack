package slack

import (
	"encoding/json"
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

func TestUnmarshalUser(t *testing.T) {
	src := `{
		"id": "U023BECGF",
		"name": "bobby",
		"deleted": false,
		"color": "9f69e7",
		"profile": {
			"first_name": "Bobby",
			"last_name": "Tables",
			"real_name": "Bobby Tables",
			"email": "bobby@slack.com",
			"skype": "my-skype-name",
			"phone": "+1 (123) 456 7890",
			"image_24": "https:\/\/www.example.org/image/24.jpg",
			"image_32": "https:\/\/www.example.org/image/32.jpg",
			"image_48": "https:\/\/www.example.org/image/48.jpg",
			"image_72": "https:\/\/www.example.org/image/72.jpg",
			"image_192": "https:\/\/www.example.org/image/192.jpg"
		},
		"is_admin": true,
		"is_owner": true,
		"is_primary_owner": true,
		"is_restricted": false,
		"is_ultra_restricted": false,
		"has_2fa": false,
		"has_files": true
	}`
	expected := User{
		ID:      "U023BECGF",
		Name:    "bobby",
		Deleted: false,
		Color:   "9f69e7",
		Profile: UserProfile{
			FirstName: "Bobby",
			LastName:  "Tables",
			RealName:  "Bobby Tables",
			Email:     "bobby@slack.com",
			Skype:     "my-skype-name",
			Phone:     "+1 (123) 456 7890",
			Image24:   "https://www.example.org/image/24.jpg",
			Image32:   "https://www.example.org/image/32.jpg",
			Image48:   "https://www.example.org/image/48.jpg",
			Image72:   "https://www.example.org/image/72.jpg",
			Image192:  "https://www.example.org/image/192.jpg",
		},
		IsAdmin:           true,
		IsOwner:           true,
		IsPrimaryOwner:    true,
		IsRestricted:      false,
		IsUltraRestricted: false,
		Has2FA:            false,
		HasFiles:          true,
	}

	u := User{}
	err := json.Unmarshal([]byte(src), &u)
	assert.NilError(t, err)
	assert.EqualFields(t, u, expected)
}

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

func TestUnmarshalChannel(t *testing.T) {
	src := `
{
    "id": "C024BE91L",
    "name": "fun",
    "is_channel": "true",
    "created": 1360782804,
    "creator": "U024BE7LH",
    "is_archived": false,
    "is_general": false,
    "members": [
        "U024BE7LH",
		"U024BE7LV"
    ],
    "topic": {
        "value": "Fun times",
        "creator": "U024BE7LV",
        "last_set": 1369677212
    },
    "purpose": {
        "value": "This channel is for fun",
        "creator": "U024BE7LH",
        "last_set": 1360782804
    },
    "is_member": true,
    "last_read": "1401383885.000061",
    "unread_count": 0,
    "unread_count_display": 0,
	"latest": {
		"text": "butt"
	}
}`
	expected := Channel{
		ID:         "C024BE91L",
		Name:       "fun",
		IsChannel:  "true",
		Created:    1360782804,
		Creator:    "U024BE7LH",
		IsArchived: false,
		IsGeneral:  false,
		Members:    []string{"U024BE7LH", "U024BE7LV"},
		Topic: Topic{
			Value:   "Fun times",
			Creator: "U024BE7LV",
			LastSet: 1369677212,
		},
		Purpose: Topic{
			Value:   "This channel is for fun",
			Creator: "U024BE7LH",
			LastSet: 1360782804,
		},
		IsMember:           true,
		LastRead:           "1401383885.000061",
		UnreadCount:        0,
		UnreadCountDisplay: 0,
		Latest: map[string]interface{}{
			"text": "butt",
		},
	}

	c := Channel{}
	err := json.Unmarshal([]byte(src), &c)
	assert.NilError(t, err)
	assert.EqualFields(t, c, expected)
}

func TestUnmarshalGroup(t *testing.T) {
	src := `
{
    "id": "G024BE91L",
    "name": "secretplans",
    "is_group": "true",
    "created": 1360782804,
    "creator": "U024BE7LH",
    "is_archived": false,
    "members": [
        "U024BE7LH"
    ],
    "topic": {
        "value": "Secret plans on hold",
        "creator": "U024BE7LV",
        "last_set": 1369677212
    },
    "purpose": {
        "value": "Discuss secret plans that no-one else should know",
        "creator": "U024BE7LH",
        "last_set": 1360782804
    },
    "last_read": "1401383885.000061",
    "latest": { "text": "butt" },
    "unread_count": 0,
    "unread_count_display": 0
}
	`
	expected := Channel{
		ID:         "G024BE91L",
		Name:       "secretplans",
		IsGroup:    "true",
		Created:    1360782804,
		Creator:    "U024BE7LH",
		IsArchived: false,
		Members:    []string{"U024BE7LH"},
		Topic: Topic{
			Value:   "Secret plans on hold",
			Creator: "U024BE7LV",
			LastSet: 1369677212,
		},
		Purpose: Topic{
			Value:   "Discuss secret plans that no-one else should know",
			Creator: "U024BE7LH",
			LastSet: 1360782804,
		},
		LastRead: "1401383885.000061",
		Latest: map[string]interface{}{
			"text": "butt",
		},
		UnreadCount:        0,
		UnreadCountDisplay: 0,
	}

	c := Channel{}
	err := json.Unmarshal([]byte(src), &c)
	assert.NilError(t, err)
	assert.EqualFields(t, c, expected)
}

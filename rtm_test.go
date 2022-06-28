package slack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRTMStart(t *testing.T) {
	rtmInfoSent := RTMStartInfo{}
	rtmInfoSent.OK = true // ok = true so we don't return an error

	rtmInfoSentJson, err := json.Marshal(rtmInfoSent)
	assert.NoError(t, err)
	server, client := jsonTestServer(200, string(rtmInfoSentJson))

	api := API{"deadbeef", server.URL + "/", client}

	rtmInfoRecieved, err := api.RTMStart()

	assert.NoError(t, err)
	assert.Equal(t, rtmInfoRecieved, rtmInfoSent)
}

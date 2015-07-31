package rtm

import (
	"net/http"

	"github.com/firba1/slack"

	"github.com/gorilla/websocket"
)

const (
	startURL = "https://slack.com/api/rtm.start"
	origin   = "http://slack.com/"

	paramToken   = "token"
	headerOrigin = "origin"
)

func Dial(token string) (*Conn, error) {
	slackAPI := slack.NewAPI(token)
	rtmConn := Conn{cancel: make(chan struct{})}

	rtmStartInfo, err := slackAPI.RTMStart()
	if err != nil {
		return &rtmConn, err
	}

	conn, err := connectWebsocket(rtmStartInfo)
	if err != nil {
		return &rtmConn, err
	}
	rtmConn.conn = conn

	// start userinfo "server"
	rtmConn.userChanges, rtmConn.infoRequests = serveUserInfo(rtmStartInfo.Users, rtmConn.cancel)

	return &rtmConn, nil
}

func connectWebsocket(rtmInfo slack.RTMStartInfo) (*websocket.Conn, error) {
	header := http.Header{}
	header.Set(headerOrigin, origin)
	conn, _, err := websocket.DefaultDialer.Dial(rtmInfo.URL, header)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

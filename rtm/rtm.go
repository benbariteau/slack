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

func Dial(token string) (conn *Conn, err error) {
	conn := Conn{cancel: make(chan struct{})}

	rtmStartInfo, err := slackAPI.RTMStart()

	conn.conn, err = rtmStart(token)
	if err != nil {
		return &conn, err
	}

	// start userinfo "server"
	conn.userChanges, conn.infoRequests = serveUserInfo(rtmStartInfo.Users, conn.cancel)

	return &conn, nil
}

func rtmStart(token string) (conn *websocket.Conn, err error) {
	rtmStartInfo, err := slack.NewAPI(token).RTMStart()
	if err != nil {
		return
	}

	conn, err := connectWebsocket(rtmStartInfo.URL)
	return
}

func connectWebsocket(url string) (*websocket.Conn, error) {
	header := http.Header{}
	header.Set(headerOrigin, origin)
	conn, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

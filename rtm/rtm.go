package rtm

import (
	"net/http"

	"github.com/benbariteau/slack"

	"github.com/gorilla/websocket"
)

const (
	startURL = "https://slack.com/api/rtm.start"
	origin   = "http://slack.com/"

	paramToken   = "token"
	headerOrigin = "origin"
)

func Dial(token string) (conn *Conn, err error) {
	conn = &Conn{cancel: make(chan struct{})}

	conn.conn, err = rtmConnect(token)
	if err != nil {
		return
	}

	return
}

func rtmStart(token string) (conn *websocket.Conn, err error) {
	rtmConnectInfo := slack.RTMConnectInfo{}
	rtmConnectInfo, err = slack.NewAPI(token).RTMConnect()
	if err != nil {
		return
	}

	conn, err = connectWebsocket(rtmConnectInfo.URL)
	return
}

func rtmStart(token string) (conn *websocket.Conn, rtmStartInfo slack.RTMStartInfo, err error) {
	rtmStartInfo, err = slack.NewAPI(token).RTMStart()
	if err != nil {
		return
	}

	conn, err = connectWebsocket(rtmStartInfo.URL)
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

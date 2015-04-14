package rtm

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

const (
	startURL = "https://slack.com/api/rtm.start"
	origin   = "http://slack.com/"

	paramToken   = "token"
	headerOrigin = "origin"
)

func Dial(token string) (*Conn, error) {
	slackConn := Conn{}

	params := url.Values{}
	params.Set(paramToken, token)
	conn, err := connectWebsocket(parseStart(http.PostForm(startURL, params)))
	if err != nil {
		return &slackConn, err
	}
	slackConn.Conn = conn
	return &slackConn, nil
}

type basicResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

type startResponse struct {
	basicResponse
	Url string `json:"url"`
}

type message struct {
	Type string `json:"type"`
}

func parseStart(resp *http.Response, err error) (startResponse, error) {
	if err != nil {
		return startResponse{}, err
	}
	return parseStartJson(ioutil.ReadAll(resp.Body))
}

func parseStartJson(bytes []byte, err error) (startResponse, error) {
	if err != nil {
		return startResponse{}, err
	}
	resp := startResponse{}
	err = json.Unmarshal(bytes, &resp)

	if err != nil {
		return startResponse{}, err
	} else if !resp.Ok {
		return startResponse{}, errors.New(resp.Error)
	}

	return resp, err
}

func connectWebsocket(startResp startResponse, err error) (*websocket.Conn, error) {
	if err != nil {
		return nil, err
	}
	header := http.Header{}
	header.Set(headerOrigin, origin)
	conn, _, err := websocket.DefaultDialer.Dial(startResp.Url, header)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

type Conn struct {
	Conn *websocket.Conn
}

func (c *Conn) Close() error {
	return c.conn.Close()
}

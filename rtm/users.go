package rtm

import (
	"github.com/firba1/slack"
)

type userInfoRequest struct {
	id     string
	respCh chan slack.User
}

/*
serverUserInfo "owns" the maps of users (user ID to user) owned by the connections that hold its channels. Any access or mutation of the users map is handled by this function alone, so as to prevent data races.

The userChanges channel is meant to be written to by user info updates from the slack server.
The infoRequests channel is used to read from the users map. The user must create a channel to recieve the response on, but this function will close it.
*/
func serveUserInfo(userList []slack.User, cancel chan struct{}) (userChanges chan<- slack.User, infoRequests chan<- userInfoRequest) {
	//build users map
	users := make(map[string]slack.User)
	for _, user := range userList {
		users[user.ID] = user
	}

	userChangesCh := make(chan slack.User)
	infoRequestsCh := make(chan userInfoRequest)

	go func() {
		for {
			select {
			case user := <-userChangesCh:
				users[user.ID] = user
			case request := <-infoRequestsCh:
				request.respCh <- users[request.id]
				close(request.respCh)
			case <-cancel:
				close(userChangesCh)
				close(infoRequestsCh)
				return
			}
		}
	}()
	return userChangesCh, infoRequestsCh
}

func (c Conn) UserInfo(id string) slack.User {
	responseChannel := make(chan slack.User)
	c.infoRequests <- userInfoRequest{
		id:     id,
		respCh: responseChannel,
	}
	return <-responseChannel
}

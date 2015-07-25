package slack

type RTMStartInfo struct {
	BasicResponse
	URL      string    `json:"url"`
	Users    []User    `json:"users"`
	//Channels []Channel `json:"channel"`
	//Groups   []Channel `json:"group"`
	//IMs      []IM      `json:"ims"`
	//TODO self, team, bots
}

func (a API) RTMStart() (r RTMStartInfo, err error) {
	err = a.request("rtm.start", make(map[string]string), &r)
	return
}

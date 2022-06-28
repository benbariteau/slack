package slack

type RTMConnectInfo struct {
	BasicResponse
	URL string `json:"url"`
}

func (a API) RTMConnect() (r RTMConnectInfo, err error) {
	err = a.request("rtm.connect", make(map[string]string), &r)
	return
}

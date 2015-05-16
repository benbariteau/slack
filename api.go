package slack

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	apiBaseURL = "https://slack.com/api/"
)

type API struct {
	Token      string
	baseURL    string
	httpClient *http.Client
}

func NewAPI(token string) API {
	return API{
		Token:      token,
		baseURL:    apiBaseURL,
		httpClient: http.DefaultClient,
	}
}

func (a API) request(method string, params map[string]string, r Response) (err error) {
	form := url.Values{}
	form.Set("token", a.Token)
	for k, v := range params {
		form.Set(k, v)
	}

	resp, err := a.httpClient.PostForm(a.baseURL+method, form)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(respData, r)
	if err != nil {
		return
	}
	return r.Error()
}

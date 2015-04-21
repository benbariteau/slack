package slack

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

type usersInfoResp struct {
	Ok   bool     `json:"ok"`
	Info UserInfo `json:"user"`
	Err  string   `json:"error"`
}

type UserInfo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Delete  bool   `json:"delete"`
	Color   string `json:"color"`
	Profile struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		RealName  string `json:"real_name"`
		Email     string `json:"email"`
		Skype     string `json:"skype"`
		Phone     string `json:"phone"`
		Image24   string `json:"image_24"`
		Image32   string `json:"image_32"`
		Image48   string `json:"image_48"`
		Image72   string `json:"image_72"`
		Image192  string `json:"image_192"`
	} `json:"profile"`
	IsAdmin  bool
	IsOwner  bool
	Has2FA   bool
	HasFiles bool
}

func (a API) UsersInfo(user string) (u UserInfo, err error) {
	form := url.Values{}
	form.Set("token", a.Token)
	form.Set("user", user)
	resp, err := http.PostForm(apiBaseURL+"users.info", form)
	if err != nil {
		return
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	uiResp := usersInfoResp{}
	err = json.Unmarshal(respData, &uiResp)
	if uiResp.Ok {
		return uiResp.Info, nil
	} else {
		return UserInfo{}, errors.New(uiResp.Err)
	}
}

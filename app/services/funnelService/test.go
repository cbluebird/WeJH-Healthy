package funnelService

import (
	"encoding/json"
	"healthy/app/apiException"
	"healthy/app/utils/fetch"
	"healthy/config/api"
	"healthy/config/user"
	"net/url"
)

func OauthTest() error {
	form := genTermForm("OAUTH")
	_, err := FetchHandleOfPost(form, api.ZFExam)
	return err
}

func ZFTest() error {
	form := genTermForm("ZF")
	_, err := FetchHandleOfPost(form, api.ZFExam)
	return err
}

func genTermForm(loginType string) url.Values {
	var password string

	if loginType == "OAUTH" {
		password = user.UserInfo.Oauth
	} else {
		password = user.UserInfo.ZF
	}

	form := url.Values{}
	form.Add("username", user.UserInfo.StudentID)
	form.Add("password", password)
	form.Add("type", loginType)
	form.Add("year", user.UserInfo.Year)
	form.Add("term", user.UserInfo.Term)
	return form
}

type FunnelResponse struct {
	Code int         `json:"code" binding:"required"`
	Msg  string      `json:"message" binding:"required"`
	Data interface{} `json:"data"`
}

func FetchHandleOfPost(form url.Values, url api.FunnelApi) (interface{}, error) {
	f := fetch.Fetch{}
	f.Init()
	res, err := f.PostForm(api.FunnelHost+string(url), form)
	if err != nil {
		return nil, apiException.RequestError
	}
	rc := FunnelResponse{}
	err = json.Unmarshal(res, &rc)
	if err != nil {
		return nil, apiException.RequestError
	}
	i := 0
	for rc.Code == 413 && i < 5 {
		i++
		res, err = f.PostForm(api.FunnelHost+string(url), form)
		if err != nil {
			return nil, apiException.RequestError
		}
		rc = FunnelResponse{}
		err = json.Unmarshal(res, &rc)
		if err != nil {
			return nil, apiException.RequestError
		}
	}

	if rc.Code == 413 {
		return rc.Data, apiException.ServerError
	}
	if rc.Code == 412 {
		return rc.Data, apiException.NoThatPasswordOrWrong
	}
	return rc.Data, nil
}

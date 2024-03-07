package captchaService

import (
	"encoding/json"
	"healthy/app/apiException"
	"healthy/app/utils/fetch"
	"healthy/config/api"
	"healthy/config/user"
	"strings"
)

type captchaServerResponse struct {
	Status int    `json:"status"`
	Data   string `json:"msg"`
}

func CaptchaTest() error {
	f := fetch.Fetch{}
	f.Init()
	_, err := f.Get(api.ZFLogin)
	if err != nil {
		return err
	}
	if len(f.Cookie) < 1 {
		return apiException.UnknownLoginError
	}
	captcha, err := f.Get(api.CaptchaHost + "?session=" + f.Cookie[0].Value + "&route=" + "/jwglxt")
	if err != nil {
		return err
	}
	captchaRes := &captchaServerResponse{}
	_ = json.Unmarshal(captcha, captchaRes)
	if captchaRes.Status != 0 {
		return apiException.WrongCaptcha
	}

	loginData := genLoginData(user.UserInfo.StudentID, user.UserInfo.ZF, f)

	s, err := f.PostForm(api.ZFLogin, loginData)
	if err != nil {
		return err
	} else if strings.Contains(string(s), "请先滑动图片进行验证") {
		return apiException.WrongCaptcha
	}
	return nil
}

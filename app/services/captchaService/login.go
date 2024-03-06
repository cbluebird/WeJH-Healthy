package captchaService

import (
	"healthy/app/utils/fetch"
	"healthy/app/utils/security"
	"healthy/config/api"
	"net/url"
)

func genLoginData(username, password string, f fetch.Fetch) url.Values {
	s, _ := f.Get(api.ZfLoginGetPublickey)
	encodePassword, _ := security.GetEncodePassword(s, []byte(password))
	return url.Values{
		"yhm": {username},
		"mm":  {encodePassword}}
}

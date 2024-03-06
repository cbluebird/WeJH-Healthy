package user

import "healthy/config/config"

type User struct {
	Oauth     string
	ZF        string
	StudentID string
	Year      string
	Term      string
}

var UserInfo User

func init() {
	UserInfo = User{
		Oauth:     "",
		ZF:        "",
		StudentID: "",
		Year:      "",
		Term:      "",
	}
	if config.Config.IsSet("user.oauth") {
		UserInfo.Oauth = config.Config.GetString("user.oauth")
	}
	if config.Config.IsSet("user.zf") {
		UserInfo.ZF = config.Config.GetString("user.zf")
	}
	if config.Config.IsSet("user.student_id") {
		UserInfo.StudentID = config.Config.GetString("user.student_id")
	}
	if config.Config.IsSet("user.year") {
		UserInfo.Year = config.Config.GetString("user.year")
	}
	if config.Config.IsSet("user.term") {
		UserInfo.Term = config.Config.GetString("user.term")
	}
}

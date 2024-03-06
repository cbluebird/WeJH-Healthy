package api

import (
	"healthy/config/config"
	"strconv"
	"time"
)

type FunnelApi string

const (
	ZFExam FunnelApi = "student/zf/exam"
)

var CaptchaHost string = config.Config.GetString("captcha")

var ZFLogin string = config.Config.GetString("zf") + "xtgl/login_slogin.html?time=" + strconv.FormatInt(time.Now().Unix()*1000, 10)

var ZfLoginGetPublickey string = config.Config.GetString("zf") + "xtgl/login_getPublicKey.html?time=" + strconv.FormatInt(time.Now().Unix()*1000, 10)

var FunnelHost = config.Config.GetString("funnel")

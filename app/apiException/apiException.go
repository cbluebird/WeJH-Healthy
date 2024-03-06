package apiException

import "net/http"

type Error struct {
	StatusCode int    `json:"-"`
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
}

var (
	ServerError           = NewError(http.StatusInternalServerError, 200500, "系统异常，请稍后重试!")
	ParamError            = NewError(http.StatusInternalServerError, 200501, "参数错误")
	NotLogin              = NewError(http.StatusBadRequest, 200503, "未登录")
	NoThatPasswordOrWrong = NewError(http.StatusBadRequest, 200504, "密码错误")
	HttpTimeout           = NewError(http.StatusInternalServerError, 200505, "网络连接超时，请稍后重试!")
	RequestError          = NewError(http.StatusInternalServerError, 200506, "服务响应错误，请稍后重试!")
	NotInit               = NewError(http.StatusNotFound, 200404, http.StatusText(http.StatusNotFound))
	NotFound              = NewError(http.StatusNotFound, 200404, http.StatusText(http.StatusNotFound))
	Unknown               = NewError(http.StatusInternalServerError, 300500, "系统未知异常，请稍后重试!")
	UnknownLoginError     = NewError(http.StatusBadRequest, 200405, "unknown login error")
	WrongCaptcha          = NewError(http.StatusBadRequest, 200406, "验证码错误")
	UnexpectedTaskType    = NewError(http.StatusInternalServerError, 200507, "未知task类型")
)

func OtherError(message string) *Error {
	return NewError(http.StatusForbidden, 100403, message)
}

func (e *Error) Error() string {
	return e.Msg
}

func NewError(statusCode, Code int, msg string) *Error {
	return &Error{
		StatusCode: statusCode,
		Code:       Code,
		Msg:        msg,
	}
}

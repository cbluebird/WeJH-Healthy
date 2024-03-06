package taskService

type TaskPayload struct {
	Cnt  int
	Type TaskType
}

type TaskType int

const (
	ZF           TaskType = 1
	ZFRetry      TaskType = 2
	Oauth        TaskType = 3
	OauthRetry   TaskType = 4
	Captcha      TaskType = 5
	CaptchaRetry TaskType = 6
)

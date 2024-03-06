package taskService

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"healthy/app/apiException"
	"healthy/app/services/captchaService"
	"healthy/app/services/emailService"
	"healthy/app/services/funnelService"
	"healthy/config/redis"
	"log"
	"os/exec"
)

func Handler(ctx context.Context, t *asynq.Task) error {
	switch t.Type() {
	case "healthy-check":
		var p TaskPayload
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return err
		}
		switch p.Type {
		case ZF:
			err := funnelService.ZFTest()
			if err != nil {
				PutRetryTask(ZFRetry, 0)
				log.Println(err)
				return err
			}
			PutTestTask(ZF)

		case ZFRetry:
			err := funnelService.ZFTest()
			p.Cnt++
			if err != nil {
				if p.Cnt == 5 {
					cmd := exec.Command("sudo systemctl restart funnel.service")
					if err = cmd.Run(); err != nil {
						log.Println(err)
					}
				} else if p.Cnt == 10 {
					redis.RedisClient.Set(ctx, "ZF_KEY", true, 0)
					PutEmailTask(ZFRetry)
				}
				PutRetryTask(ZFRetry, p.Cnt)
				return err
			}
			redis.RedisClient.Del(ctx, "ZF_KEY")
			PutTestTask(ZF)

		case Oauth:
			err := funnelService.OauthTest()
			if err != nil {
				PutRetryTask(OauthRetry, 0)
				log.Println(err)
				return err
			}
			PutTestTask(Oauth)

		case OauthRetry:
			err := funnelService.OauthTest()
			p.Cnt++
			if err != nil {
				if p.Cnt == 5 {
					cmd := exec.Command("sudo systemctl restart funnel.service")
					if err = cmd.Run(); err != nil {
						log.Println(err)
					}
				} else if p.Cnt == 10 {
					redis.RedisClient.Set(ctx, "Oauth_KEY", true, 0)
					PutEmailTask(p.Type)
				}
				PutRetryTask(p.Type, p.Cnt)
				return err
			}
			redis.RedisClient.Del(ctx, "Oauth_KEY")
			PutTestTask(Oauth)

		case Captcha:
			err := captchaService.CaptchaTest()
			if err != nil {
				PutRetryTask(CaptchaRetry, 0)
				log.Println(err)
				return err
			}
			PutTestTask(Captcha)
		case CaptchaRetry:
			err := captchaService.CaptchaTest()
			p.Cnt++
			if err != nil {
				if p.Cnt == 5 {
					cmd := exec.Command("sudo systemctl restart zf.service")
					if err = cmd.Run(); err != nil {
						log.Println(err)
					}
				} else if p.Cnt == 10 {
					PutEmailTask(CaptchaRetry)
				}
				PutRetryTask(CaptchaRetry, p.Cnt)
				return err
			}
			PutTestTask(Captcha)
		}

	case "email":
		var p TaskPayload
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return err
		}
		service := ""
		switch p.Type {
		case ZFRetry:
			service = "正方教务"
		case OauthRetry:
			service = "统一验证"
		case CaptchaRetry:
			service = "正方验证码"
		}
		emailService.SendEmail(service)

	default:
		return apiException.UnexpectedTaskType
	}
	return nil
}

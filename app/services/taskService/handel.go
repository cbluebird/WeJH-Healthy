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
	"healthy/config/task"
	"log"
	"os/exec"
	"time"
)

var ctx context.Context

func clearAll(info asynq.RedisClientOpt) error {
	inspector := asynq.NewInspector(info)

	if err := inspector.DeleteQueue("healthy-check", true); err != nil {
		return err
	}
	if err := inspector.DeleteQueue("healthy-retry", true); err != nil {
		return err
	}
	return nil
}

func Init() (*asynq.ServeMux, error) {
	info := task.Init()
	task.AsynqClient = asynq.NewClient(info)
	task.AsynqServer = asynq.NewServer(
		info,
		asynq.Config{
			Concurrency:    10, //Concurrency表示最大并发处理任务数。
			RetryDelayFunc: retryFunc,
			Queues: map[string]int{
				"email":         2,
				"healthy-check": 4,
				"healthy-retry": 4,
			},
		},
	)
	// 遍历每个队列，清空队列中的消息
	if err := clearAll(info); err != nil {
		log.Println(err)
	}

	if err := CheckInit(); err != nil {
		log.Fatal(err)
	}

	redis.RedisClient.Del(ctx, "email")
	mux := asynq.NewServeMux()
	mux.HandleFunc("email", emailHandel)
	mux.HandleFunc("healthy-check", healthyHandel)
	mux.HandleFunc("healthy-retry", retryHandel)
	return mux, nil
}

func emailHandel(ctx context.Context, t *asynq.Task) error {
	var p TaskPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}
	service := ""
	switch p.Type {
	case ZF:
		service = "正方教务"
	case Oauth:
		service = "统一验证"
	case Captcha:
		service = "正方验证码"
	default:
		return apiException.UnexpectedTaskType
	}
	return emailService.SendEmail(service)
}

func healthyHandel(ctx context.Context, t *asynq.Task) error {
	var p TaskPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}
	switch p.Type {
	case ZF:
		err := funnelService.ZFTest()
		if err != nil {
			return err
		}
		PutTestTask(ZF)

	case Captcha:
		err := captchaService.CaptchaTest()
		if err != nil {
			return err
		}
		PutTestTask(Captcha)

	case Oauth:
		err := funnelService.OauthTest()
		if err != nil {
			return err
		}
		PutTestTask(Oauth)
	default:
		return apiException.UnexpectedTaskType
	}

	return nil
}

func retryHandel(ctx context.Context, t *asynq.Task) error {
	var p TaskPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}
	switch p.Type {

	case ZFRetry:
		err := funnelService.ZFTest()
		if err == nil {
			redis.RedisClient.Del(ctx, "ZF_KEY")
			PutTestTask(ZF)
		}
		return err

	case OauthRetry:
		err := funnelService.OauthTest()
		if err == nil {
			redis.RedisClient.Del(ctx, "Oauth_KEY")
			PutTestTask(Oauth)
		}
		return err

	case CaptchaRetry:
		err := captchaService.CaptchaTest()
		if err == nil {
			PutTestTask(Captcha)
		}
		return err
	default:
		return apiException.UnexpectedTaskType
	}
}

func retryFunc(n int, e error, t *asynq.Task) time.Duration {
	switch t.Type() {
	case "healthy-check":
		if n == 5 {
			var p TaskPayload
			if err := json.Unmarshal(t.Payload(), &p); err != nil {
			}
			var cmd *exec.Cmd
			switch p.Type {
			case ZF:
				cmd = exec.Command("sudo systemctl restart funnel.service")
			case Oauth:
				cmd = exec.Command("sudo systemctl restart funnel.service")
			case Captcha:
				cmd = exec.Command("sudo systemctl restart zf.service")
			}
			if err := cmd.Run(); err != nil {
				log.Println(err)
			}
		} else if n == 9 {
			var p TaskPayload
			if err := json.Unmarshal(t.Payload(), &p); err != nil {
			}
			switch p.Type {
			case ZF:
				redis.RedisClient.Set(ctx, "ZF_KEY", true, 0)
				PutRetryTask(ZFRetry)
			case Oauth:
				redis.RedisClient.Set(ctx, "Oauth_KEY", true, 0)
				PutRetryTask(OauthRetry)
			case Captcha:
				PutRetryTask(CaptchaRetry)
			}
			PutEmailTask(p.Type)
		}

	case "healthy-retry":
		var p TaskPayload
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
		}
		if n == 4 {
			PutRetryTask(p.Type)
		}
	}

	return 5 * time.Minute
}

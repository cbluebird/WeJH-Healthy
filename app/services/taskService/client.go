package taskService

import (
	"encoding/json"
	"github.com/hibiken/asynq"
	"healthy/config/task"
	"log"
	"time"
)

func PutTestTask(t TaskType) {
	payload, err := json.Marshal(TaskPayload{
		Type: t,
	})
	if err != nil {
		log.Println(err)
		return
	}
	t1 := asynq.NewTask("healthy-check", payload)
	_, err = task.AsynqClient.Enqueue(t1, asynq.ProcessIn(6*time.Hour), asynq.MaxRetry(10), asynq.Queue("healthy-check"))
	if err != nil {
		log.Println(err)
		return
	}
}

func PutRetryTask(t TaskType) {
	payload, err := json.Marshal(TaskPayload{
		Type: t,
	})
	if err != nil {
		log.Println(err)
		return
	}
	t1 := asynq.NewTask("healthy-retry", payload)
	_, err = task.AsynqClient.Enqueue(t1, asynq.ProcessIn(30*time.Minute), asynq.MaxRetry(5), asynq.Queue("healthy-retry"))
	if err != nil {
		log.Println(err)
		return
	}
}

func PutEmailTask(t TaskType) {
	payload, err := json.Marshal(TaskPayload{
		Type: t,
	})
	if err != nil {
		log.Println(err)
		return
	}
	t1 := asynq.NewTask("email", payload)
	_, err = task.AsynqClient.Enqueue(t1, asynq.MaxRetry(3), asynq.Queue("email"))
	if err != nil {
		log.Println(err)
		return
	}
}

func CheckInit() error {
	payload1, err := json.Marshal(TaskPayload{
		Type: ZF,
	})
	payload2, err := json.Marshal(TaskPayload{
		Type: Oauth,
	})
	payload3, err := json.Marshal(TaskPayload{
		Type: Captcha,
	})
	if err != nil {
		log.Println(err)
		return err
	}

	t1 := asynq.NewTask("healthy-check", payload1)
	t2 := asynq.NewTask("healthy-check", payload2)
	t3 := asynq.NewTask("healthy-check", payload3)

	if _, err = task.AsynqClient.Enqueue(t1, asynq.MaxRetry(10), asynq.Queue("healthy-check")); err != nil {
		return err
	}
	if _, err = task.AsynqClient.Enqueue(t2, asynq.MaxRetry(10), asynq.Queue("healthy-check")); err != nil {
		return err
	}
	if _, err = task.AsynqClient.Enqueue(t3, asynq.MaxRetry(10), asynq.Queue("healthy-check")); err != nil {
		return err
	}
	return nil
}

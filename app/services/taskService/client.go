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
		Cnt:  0,
		Type: t,
	})
	if err != nil {
		log.Println(err)
		return
	}
	t1 := asynq.NewTask("healthy-check", payload)
	_, err = task.AsynqClient.Enqueue(t1, asynq.ProcessIn(6*time.Hour))
	if err != nil {
		log.Println(err)
		return
	}
}

func PutRetryTask(t TaskType, cnt int) {
	payload, err := json.Marshal(TaskPayload{
		Cnt:  cnt,
		Type: t,
	})
	if err != nil {
		log.Println(err)
		return
	}
	t1 := asynq.NewTask("healthy-check", payload)
	_, err = task.AsynqClient.Enqueue(t1, asynq.ProcessIn(5*time.Minute))
	if err != nil {
		log.Println(err)
		return
	}
}

func PutEmailTask(t TaskType) {
	payload, err := json.Marshal(TaskPayload{
		Cnt:  0,
		Type: t,
	})
	if err != nil {
		log.Println(err)
		return
	}
	t1 := asynq.NewTask("email", payload)
	_, err = task.AsynqClient.Enqueue(t1)
	if err != nil {
		log.Println(err)
		return
	}
}

func Init() {

}

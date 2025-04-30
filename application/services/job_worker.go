package services

import (
	"encoder/domain"
	"encoder/framework/utils"
	"encoding/json"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

type JobWorkerResult struct {
	Job     *domain.Job
	Message *amqp.Delivery
	Error   error
}

func JobWorker(messageChannel chan amqp.Delivery, returnChan chan JobWorkerResult, jobService JobService, job *domain.Job, workerID int) {
	for message := range messageChannel {
		processMessage(&message, returnChan, jobService, job, workerID)
	}
}

func processMessage(message *amqp.Delivery, returnChan chan JobWorkerResult, jobService JobService, job *domain.Job, workerID int) {
	var err error

	if err = utils.IsJson(string(message.Body)); err != nil {
		returnJobResult(&domain.Job{}, message, err)
		return
	}

	if err = json.Unmarshal(message.Body, &jobService.VideoService.Video); err != nil {
		returnJobResult(&domain.Job{}, message, err)
		return
	}

	jobService.VideoService.Video.ID = uuid.NewString()

	if err = jobService.VideoService.Video.Validate(); err != nil {
		returnJobResult(&domain.Job{}, message, err)
		return
	}

	if err = jobService.VideoService.InsertVideo(); err != nil {
		returnJobResult(&domain.Job{}, message, err)
		return
	}

	job.Video = jobService.VideoService.Video
	job.OutputBucketPath = os.Getenv(OUTPUT_BUCKET_NAME)
	job.ID = uuid.NewString()
	job.CreatedAt = time.Now()
	job.Status, err = domain.JobStatusStarting.Description()

	if err != nil {
		returnJobResult(&domain.Job{}, message, err)
		return
	}

	if _, err = jobService.JobRepository.Insert(job); err != nil {
		returnJobResult(&domain.Job{}, message, err)
		return
	}

	jobService.Job = job

	if err = jobService.Start(); err != nil {
		returnJobResult(&domain.Job{}, message, err)
		return
	}

	returnChan <- returnJobResult(job, message, nil)

}

func returnJobResult(job *domain.Job, message *amqp.Delivery, err error) JobWorkerResult {
	result := JobWorkerResult{
		Job:     job,
		Message: message,
		Error:   err,
	}

	return result

}

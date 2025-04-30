package services

import (
	"encoder/application/repositories"
	"encoder/domain"
	"errors"
	"os"
	"strconv"
)

type JobService struct {
	Job           *domain.Job
	JobRepository repositories.JobRepository
	VideoService  VideoService
}

func (j *JobService) Start() error {
	err := j.changeJobStatus(domain.JobStatusDownloading)
	if err != nil {
		return j.failJob(err)
	}

	err = j.VideoService.Download(os.Getenv("INPUT_BUCKET_NAME"))
	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus(domain.JobStatusFragmenting)
	if err != nil {
		return j.failJob(err)
	}

	err = j.VideoService.Fragment()
	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus(domain.JobStatusEncoding)
	if err != nil {
		return j.failJob(err)
	}

	err = j.VideoService.Encode()
	if err != nil {
		return j.failJob(err)
	}

	err = j.performUpload()
	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus(domain.JobStatusFinishing)
	if err != nil {
		return j.failJob(err)
	}

	err = j.VideoService.Finish()
	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus(domain.JobStatusCompleted)
	if err != nil {
		return j.failJob(err)
	}

	return nil
}

func (j *JobService) performUpload() error {
	err := j.changeJobStatus(domain.JobStatusUploading)
	if err != nil {
		return j.failJob(err)
	}

	videoUpload := NewVideoUpload()
	videoUpload.OutputBucket = os.Getenv(OUTPUT_BUCKET_NAME)
	videoUpload.VideoPath = os.Getenv(LOCAL_STORAGE_PATH) + "/" + j.VideoService.Video.ID

	concurrency, _ := strconv.Atoi(os.Getenv("CONCURRENCY_UPLOAD"))
	doneUpload := make(chan string)

	go videoUpload.ProcessUpload(concurrency, doneUpload)
	uploadResult := <-doneUpload

	if uploadResult != "upload completed" {
		return j.failJob(errors.New(uploadResult))
	}

	return nil
}

func (j *JobService) changeJobStatus(status domain.JobStatus) error {
	var err error

	j.Job.Status, err = status.Description()
	if err != nil {
		return j.failJob(err)
	}

	j.Job, err = j.JobRepository.Update(j.Job)

	if err != nil {
		return j.failJob(err)
	}

	return nil
}

func (j *JobService) failJob(err error) error {
	j.Job.Error = err.Error()
	j.Job.Status, err = domain.JobStatusFailed.Description()

	if err != nil {
		return err
	}

	_, err = j.JobRepository.Update(j.Job)
	if err != nil {
		return err
	}

	return nil
}

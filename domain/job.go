package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

type Job struct {
	ID               string    `json:"job_id" valid:"uuid" gorm:"type:uuid;primary_key"`
	OutputBucketPath string    `json:"output_bucket_path" valid:"notnull"`
	Status           string    `json:"status" valid:"notnull"`
	Video            *Video    `json:"video" valid:"-"`
	VideoID          string    `json:"-" valid:"-" gorm:"column:video_id;type:uuid;notnull"`
	Error            string    `valid:"-"`
	CreatedAt        time.Time `json:"created_at" valid:"-"`
	UpdatedAt        time.Time `json:"updated_at" valid:"-"`
}

type JobStatus int

const (
	JobStatusStarting JobStatus = iota
	JobStatusDownloading
	JobStatusUploading
	JobStatusFragmenting
	JobStatusEncoding
	JobStatusFinishing
	JobStatusCompleted
	JobStatusFailed
)

func (j JobStatus) Description() (string, error) {
	switch j {
	case JobStatusStarting:
		return "STARTING", nil
	case JobStatusDownloading:
		return "DOWNLOADING", nil
	case JobStatusUploading:
		return "UPLOADING", nil
	case JobStatusFragmenting:
		return "FRAGMENTING", nil
	case JobStatusEncoding:
		return "ENCODING", nil
	case JobStatusFinishing:
		return "FINISHING", nil
	case JobStatusCompleted:
		return "COMPLETED", nil
	case JobStatusFailed:
		return "FAILED", nil
	}

	return "", errors.New(fmt.Sprintf("invalid job status %v", j))

}

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

func NewJob(output string, status string, video *Video) (*Job, error) {
	job := Job{
		OutputBucketPath: output,
		Status:           status,
		Video:            video,
	}

	job.prepare()

	if err := job.Validate(); err != nil {
		return nil, err
	}

	return &job, nil
}

func (job *Job) prepare() {
	job.ID = uuid.NewString()
	job.CreatedAt = time.Now()
	job.UpdatedAt = time.Now()
}

func (job *Job) Validate() error {
	_, err := govalidator.ValidateStruct(job)
	if err != nil {
		return err
	}

	return nil

}

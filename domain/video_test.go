package domain_test

import (
	"encoder/domain"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestValidateIfVideoIsEmpty(t *testing.T) {
	video := domain.NewVideo()
	err := video.Validate()

	require.Error(t, err)
}

func TestVideoIdIsNotAUuid(t *testing.T) {
	video := domain.NewVideo()
	video.ID = "abc"
	video.ResourceID = "a"
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	err := video.Validate()
	require.Error(t, err)
}

func TestVideoValidation(t *testing.T) {
	video := domain.NewVideo()
	video.ID = uuid.New().String()
	video.ResourceID = "a"
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	err := video.Validate()
	require.Nil(t, err)
}

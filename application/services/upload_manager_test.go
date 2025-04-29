package services_test

import (
	"encoder/application/services"
	"os"
	"testing"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func TestVideoServiceUpload(t *testing.T) {
	video, repo := prepare()

	videoService := services.NewVideoService()
	videoService.Video = video
	videoService.VideoRepository = repo

	err := videoService.Download("fc-go-encoder")
	require.Nil(t, err)

	err = videoService.Fragment()
	require.Nil(t, err)

	err = videoService.Encode()
	require.Nil(t, err)

	videoUpload := services.NewVideoUpload()
	videoUpload.OutputBucket = "fc-go-encoder"
	videoUpload.VideoPath = os.Getenv(services.LOCAL_STORAGE_PATH) + "/" + video.ID

	doneUpload := make(chan string)
	go videoUpload.ProcessUpload(25, doneUpload)

	result := <-doneUpload
	require.Equal(t, result, "upload completed")

}

package services

import (
	"context"
	"encoder/application/repositories"
	"encoder/domain"
	"io"
	"os"
	"os/exec"

	"cloud.google.com/go/storage"
	log "github.com/sirupsen/logrus"
)

type VideoService struct {
	Video           *domain.Video
	VideoRepository repositories.VideoRepository
}

func NewVideoService() VideoService {
	return VideoService{}
}

func (v *VideoService) Download(bucketName string) error {

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	bucket := client.Bucket(bucketName)
	object := bucket.Object(v.Video.FilePath)

	reader, err := object.NewReader(ctx)
	if err != nil {
		return err
	}
	defer reader.Close()

	body, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	file, err := os.Create(v.getVideoPath() + ".mp4")
	if err != nil {
		return err
	}

	_, err = file.Write(body)
	if err != nil {
		return err
	}
	defer file.Close()

	log.Infof("video %v has been stored", v.Video.ID)
	return nil
}

func (v *VideoService) Fragment() error {
	err := os.Mkdir(v.getVideoPath(), os.ModePerm)
	if err != nil {
		return err
	}

	source := v.getVideoPath() + ".mp4"
	target := v.getVideoPath() + ".frag"

	cmd := exec.Command("mp4fragment", source, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

func (v *VideoService) Encode() error {
	cmdArgs := []string{}
	cmdArgs = append(cmdArgs, v.getVideoPath()+".frag")
	cmdArgs = append(cmdArgs, "--use-segment-timeline")
	cmdArgs = append(cmdArgs, "-o")
	cmdArgs = append(cmdArgs, v.getVideoPath())
	cmdArgs = append(cmdArgs, "-f")
	cmdArgs = append(cmdArgs, "--exec-dir")
	cmdArgs = append(cmdArgs, "/opt/bento4/bin/")
	cmd := exec.Command("mp4dash", cmdArgs...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

func (v *VideoService) Finish() error {
	err := os.Remove(v.getVideoPath() + ".mp4")
	if err != nil {
		log.Errorln("error removing mp4 ", v.Video.ID, ".mp4")
		return err
	}

	err = os.Remove(v.getVideoPath() + ".frag")
	if err != nil {
		log.Errorln("error removing frag ", v.Video.ID, ".frag")
		return err
	}

	err = os.RemoveAll(v.getVideoPath())
	if err != nil {
		log.Errorln("error removing files ", v.Video.ID)
		return err
	}

	log.Info("files have been removed")
	return nil

}

func (v *VideoService) getVideoPath() string {
	return os.Getenv(LOCAL_STORAGE_PATH) + "/" + v.Video.ID
}

func printOutput(out []byte) {
	if len(out) > 0 {
		log.Infof("====> Output %s\n", string(out))
	}
}

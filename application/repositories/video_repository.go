package repositories

import (
	"encoder/domain"
	"fmt"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type VideoRepository interface {
	Insert(video *domain.Video) (*domain.Video, error)
	Find(id string) (*domain.Video, error)
}

type VideoRepositoryDb struct {
	Db *gorm.DB
}

func NewVideoRepository(db *gorm.DB) *VideoRepositoryDb {
	return &VideoRepositoryDb{Db: db}
}

func (repo VideoRepositoryDb) Insert(video *domain.Video) (*domain.Video, error) {
	if video.ID == "" {
		video.ID = uuid.NewString()
	}

	err := repo.Db.Create(video).Error
	if err != nil {
		return nil, err
	}

	return video, nil
}

func (repo VideoRepositoryDb) Find(id string) (*domain.Video, error) {
	var video domain.Video

	err := repo.Db.Preload("Jobs").First(&video, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	if video.ID == "" {
		return nil, fmt.Errorf("Video does not exist: %s", id)
	}

	return &video, nil

}

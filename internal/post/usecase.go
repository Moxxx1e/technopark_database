package post

import (
	"github.com/technopark_database/internal/helpers/errors"
	"github.com/technopark_database/internal/models"
)

type PostUseCase interface {
	CreateMany(posts []*models.Post) *errors.Error
	ChangeByID(id uint64, message string) (*models.Post, *errors.Error)
}

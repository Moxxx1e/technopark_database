package post

import (
	"github.com/technopark_database/internal/helpers/errors"
	"github.com/technopark_database/internal/models"
)

type PostUseCase interface {
	CreateMany(slugOrID string, posts []*models.Post) ([]*models.Post, *errors.Error)
	ChangeByID(id uint64, message string) (*models.Post, *errors.Error)
	GetPosts(slugOrID string, sort string, since uint64,
		pagination *models.Pagination) ([]*models.Post, *errors.Error)
	GetPostInfo(id uint64, related *models.Related) (*models.PostDetails, *errors.Error)
}

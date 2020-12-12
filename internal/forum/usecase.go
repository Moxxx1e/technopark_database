package forum

import (
	"github.com/technopark_database/internal/helpers/errors"
	"github.com/technopark_database/internal/models"
)

type ForumUseCase interface {
	Create(forum *models.Forum) (*models.Forum, *errors.Error)
	GetDetails(slug string) (*models.Forum, *errors.Error)
	GetUsers(slug string, pagination *models.Pagination) ([]*models.User, *errors.Error)
	GetThreads(slug string, pagination *models.Pagination) ([]*models.Thread, *errors.Error)
}

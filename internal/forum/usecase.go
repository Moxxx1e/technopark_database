package forum

import (
	"github.com/technopark_database/internal/helpers/errors"
	"github.com/technopark_database/internal/models"
)

type ForumUseCase interface {
	Create(forum *models.Forum) *errors.Error
	GetDetails(slug string) (*models.Forum, *errors.Error)
	GetUsers(slug string) ([]*models.User, *errors.Error)
	GetThreads(slug string) ([]*models.Thread, *errors.Error)
}

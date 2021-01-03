package forum

import (
	"github.com/technopark_database/internal/helpers/errors"
	"github.com/technopark_database/internal/models"
)

type ForumUseCase interface {
	Create(forum *models.Forum) (*models.Forum, *errors.Error)
	AddForumUser(nickname string, slug string) *errors.Error
	GetDetails(slug string) (*models.Forum, *errors.Error)
	GetFullDetails(slug string) (*models.Forum, *errors.Error)
	GetUsers(slug string, since string, pagination *models.Pagination) ([]*models.User, *errors.Error)
	GetThreads(slug string, since string, pagination *models.Pagination) ([]*models.Thread, *errors.Error)
}

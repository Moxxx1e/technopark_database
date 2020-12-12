package thread

import (
	"github.com/technopark_database/internal/helpers/errors"
	"github.com/technopark_database/internal/models"
)

type ThreadUsecase interface {
	Create(thread *models.Thread) *errors.Error
	CreateVoteByID(id uint64, nickname string, vote int) (*models.Thread, *errors.Error)
	CreateVoteBySlug(slug string, nickname string, vote int) (*models.Thread, *errors.Error)
	ChangeByID(id uint64, title, message string) (*models.Thread, *errors.Error)
	ChangeBySlug(slug string, title, message string) (*models.Thread, *errors.Error)
	GetByID(id uint64) (*models.Thread, *errors.Error)
	GetBySlug(slug string) (*models.Thread, *errors.Error)
}

package vote

import (
	"github.com/technopark_database/internal/helpers/errors"
	"github.com/technopark_database/internal/models"
)

type VoteUseCase interface {
	Create(vote *models.Vote) *errors.Error
	Get(threadID uint64, userID uint64) (*models.Vote, *errors.Error)
	Update(vote *models.Vote) *errors.Error
	Delete(vote *models.Vote) *errors.Error
}

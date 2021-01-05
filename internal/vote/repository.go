package vote

import "github.com/technopark_database/internal/models"

type VoteRepository interface {
	Insert(vote *models.Vote) error
	Update(vote *models.Vote) error
	Delete(vote *models.Vote) error
	SelectByThreadIDUserID(threadID uint64, userID uint64) (*models.Vote, error)
}

package usecases

import (
	"database/sql"
	"github.com/technopark_database/internal/consts"
	"github.com/technopark_database/internal/helpers/errors"
	"github.com/technopark_database/internal/models"
	"github.com/technopark_database/internal/vote"
)

type VoteUseCase struct {
	rep vote.VoteRepository
}

func (uc *VoteUseCase) Create(vote *models.Vote) (bool, *errors.Error) {
	dbVote, customErr := uc.Get(vote.ThreadID, vote.UserID)
	if customErr == errors.Get(consts.CodeVoteDoesNotExist) {
		err := uc.rep.Insert(vote)
		if err != nil {
			return false, errors.New(consts.CodeInternalServerError, err)
		}
	} else if customErr == nil {
		if dbVote.Likes == vote.Likes {
			return false, nil
		}

		err := uc.rep.Update(vote)
		if err != nil {
			return false, errors.New(consts.CodeInternalServerError, err)
		}
	} else {
		return false, customErr
	}
	return true, nil
}

func (uc *VoteUseCase) Get(threadID uint64, userID uint64) (*models.Vote, *errors.Error) {
	vote, err := uc.rep.SelectByThreadIDUserID(threadID, userID)
	if err == sql.ErrNoRows {
		return nil, errors.Get(consts.CodeVoteDoesNotExist)
	} else if err != nil {
		return nil, errors.New(consts.CodeInternalServerError, err)
	}
	return vote, nil
}

func (uc *VoteUseCase) Update(vote *models.Vote) *errors.Error {
	dbVote, customErr := uc.Get(vote.ThreadID, vote.UserID)
	if customErr != nil {
		return customErr
	}

	if dbVote.Likes == vote.Likes {
		return errors.Get(consts.CodeVoteAlreadyExist)
	}

	err := uc.rep.Update(vote)
	if err != nil {
		return errors.New(consts.CodeInternalServerError, err)
	}
	return nil
}

func (uc *VoteUseCase) Delete(vote *models.Vote) *errors.Error {
	_, customErr := uc.Get(vote.ThreadID, vote.UserID)
	if customErr != nil {
		return customErr
	}

	err := uc.rep.Delete(vote)
	if err != nil {
		return errors.New(consts.CodeInternalServerError, err)
	}
	return nil
}

func NewVoteUseCase(rep vote.VoteRepository) vote.VoteUseCase {
	return &VoteUseCase{rep: rep}
}

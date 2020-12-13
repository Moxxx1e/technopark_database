package repository

import (
	"context"
	"database/sql"
	"github.com/sirupsen/logrus"
	"github.com/technopark_database/internal/models"
	"github.com/technopark_database/internal/vote"
)

type VotePgRepository struct {
	db *sql.DB
}

func NewVoteRepository(db *sql.DB) vote.VoteRepository {
	return &VotePgRepository{db: db}
}

func (rep *VotePgRepository) Insert(vote *models.Vote) error {
	tx, err := rep.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO votes(thread_id, user_id, likes)
		VALUES ($1, $2, $3)`,
		vote.ThreadID, vote.UserID, vote.Likes)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logrus.Info(rollbackErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (rep *VotePgRepository) Update(vote *models.Vote) error {
	tx, err := rep.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE votes
		SET likes=$1
		WHERE user_id=$2 AND thread_id=$3`,
		vote.Likes, vote.UserID, vote.ThreadID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logrus.Info(rollbackErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (rep *VotePgRepository) Delete(vote *models.Vote) error {
	tx, err := rep.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE
		FROM votes
		WHERE thread_id=$1 AND user_id=$2`,
		vote.ThreadID, vote.UserID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logrus.Info(rollbackErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
func (rep *VotePgRepository) SelectByThreadIDUserID(threadID uint64,
	userID uint64) (*models.Vote, error) {
	vote := &models.Vote{}
	err := rep.db.QueryRow(`
		SELECT thread_id, user_id, likes
		FROM votes
		WHERE thread_id=$1 AND user_id=$2`, threadID, userID).
		Scan(&vote.ThreadID, &vote.UserID, &vote.Likes)
	if err != nil {
		return nil, err
	}
	return vote, nil
}

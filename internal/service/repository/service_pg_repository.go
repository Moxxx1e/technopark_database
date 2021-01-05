package repository

import (
	"context"
	"database/sql"
	"github.com/technopark_database/internal/models"
	"github.com/technopark_database/internal/service"
)

type ServicePgRepository struct {
	db *sql.DB
}

func (rep ServicePgRepository) GetStatus() (*models.ServiceStatus, error) {
	serviceStatus := &models.ServiceStatus{}
	row := rep.db.QueryRow(`
		SELECT
		(SELECT count(*) from forums) AS forum,
		(SELECT count(*) from posts) AS post,
		(SELECT count(*) from threads) AS thread,
		(SELECT count(*) from users) AS user`)

	err := row.Scan(&serviceStatus.Forum, &serviceStatus.Post,
		&serviceStatus.Thread, &serviceStatus.User)
	if err != nil {
		return nil, err
	}

	return serviceStatus, nil
}

func NewServicePgRepository(db *sql.DB) service.ServiceRepository {
	return &ServicePgRepository{db: db}
}

func (rep ServicePgRepository) Delete() error {
	tx, err := rep.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		TRUNCATE users, forums, posts, threads, votes RESTART IDENTITY CASCADE`)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

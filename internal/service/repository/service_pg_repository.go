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
	rows, err := rep.db.Query(`
		SELECT reltuples::bigint AS estimate
		FROM pg_class
		WHERE relname='users' OR relname='forums' 
			OR relname='threads' OR relname='posts'
		ORDER BY relname`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&serviceStatus.Forum)
		if err != nil {
			return nil, err
		}

		err = rows.Scan(&serviceStatus.Post)
		if err != nil {
			return nil, err
		}

		err = rows.Scan(&serviceStatus.Thread)
		if err != nil {
			return nil, err
		}

		err = rows.Scan(&serviceStatus.User)
		if err != nil {
			return nil, err
		}
	}

	if rows.Err() != nil {
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

	_, err = tx.Exec(`TRUNCATE TABLE users CASCADE;
				TRUNCATE TABLE forums CASCADE;
				TRUNCATE TABLE posts CASCADE;
				TRUNCATE TABLE threads CASCADE`)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

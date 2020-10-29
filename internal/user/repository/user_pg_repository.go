package repository

import (
	"context"
	"database/sql"
	"github.com/technopark_database/internal/models"
)

type UserPgRepository struct {
	db *sql.DB
}

func NewUserPgRepository(db *sql.DB) *UserPgRepository {
	return &UserPgRepository{db: db}
}

func (ur *UserPgRepository) Insert(user *models.User) error {
	tx, err := ur.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT INTO profile(nickname, fullname, about, email)
		VALUES ($1, $2, $3, $4)`,
		user.Nickname, user.Fullname, user.About, user.Email)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (ur *UserPgRepository) Select(nickname string) (*models.User, error) {
	dbUser := &models.User{}
	err := ur.db.QueryRow(`
		SELECT nickname, fullname, about, email
		FROM profile
		WHERE nickname=$1`, nickname).Scan(
		&dbUser.Nickname, &dbUser.Fullname,
		&dbUser.About, &dbUser.Email)
	if err != nil {
		return nil, err
	}
	return dbUser, nil
}

func (ur *UserPgRepository) Update(nickname string, newUser *models.User) error {
	tx, err := ur.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE profile 
		SET nickname = $1,
		    fullname = $2,
            about = $3,
		    email = $4
		WHERE nickname=$5`, newUser.Nickname, newUser.Fullname, newUser.About, newUser.Email, nickname)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}

func (ur *UserPgRepository) SelectByEmail(email string) (*models.User, error) {
	dbUser := &models.User{}
	err := ur.db.QueryRow(`
		SELECT nickname, fullname, about, email
		FROM profile WHERE email=$1`, email).Scan(
		&dbUser.Nickname, &dbUser.Fullname,
		&dbUser.About, &dbUser.Email)
	if err != nil {
		return nil, err
	}
	return dbUser, nil
}


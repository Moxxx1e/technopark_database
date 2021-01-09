package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/technopark_database/internal/models"
	"github.com/technopark_database/internal/user"
	"strings"
)

type UserPgRepository struct {
	db *sql.DB
}

func buildValuesQuery(valuesCount int) string {
	var values []string
	for i := 1; i <= valuesCount; i++ {
		values = append(values, fmt.Sprintf("$%d", i))
	}
	valuesQuery := fmt.Sprintf("(%s)", strings.Join(values, ", "))
	return valuesQuery
}

func (ur *UserPgRepository) SelectCountNicknames(nicknames []string) (int, error) {
	selectQuery := "SELECT COUNT(nickname) FROM users WHERE nickname IN"
	valuesQuery := buildValuesQuery(len(nicknames))
	query := strings.Join([]string{
		selectQuery,
		valuesQuery,
	}, " ")

	var values []interface{}
	for _, nickname := range nicknames {
		values = append(values, nickname)
	}

	var usersCount int
	err := ur.db.QueryRow(query, values...).Scan(&usersCount)
	if err != nil {
		return 0, err
	}
	return usersCount, nil
}

func NewUserPgRepository(db *sql.DB) user.UserRepository {
	return &UserPgRepository{db: db}
}

func (ur *UserPgRepository) Insert(user *models.User) error {
	tx, err := ur.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	err = tx.QueryRow(`
		INSERT INTO users(nickname, fullname, about, email)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		user.Nickname, user.Fullname, user.About, user.Email).Scan(&user.ID)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			logrus.Error(err)
		}
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
		SELECT id, nickname, fullname, about, email
		FROM users
		WHERE lower(nickname)=$1`, strings.ToLower(nickname)).Scan(
		&dbUser.ID, &dbUser.Nickname, &dbUser.Fullname,
		&dbUser.About, &dbUser.Email)
	if err != nil {
		return nil, err
	}
	return dbUser, nil
}

func (ur *UserPgRepository) Update(newUser *models.User) error {
	tx, err := ur.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE users
		SET email = $2, about = $3, fullname = $4
		WHERE nickname = $1`,
		newUser.Nickname, newUser.Email, newUser.About, newUser.Fullname)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			logrus.Error(err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (ur *UserPgRepository) SelectByEmail(email string) (*models.User, error) {
	dbUser := &models.User{}
	err := ur.db.QueryRow(`
		SELECT id, nickname, fullname, about, email
		FROM users
		WHERE lower(email) = $1`, strings.ToLower(email)).
		Scan(&dbUser.ID, &dbUser.Nickname, &dbUser.Fullname,
			&dbUser.About, &dbUser.Email)
	if err != nil {
		return nil, err
	}
	return dbUser, nil
}

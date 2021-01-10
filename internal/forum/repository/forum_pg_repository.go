package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/technopark_database/internal/forum"
	"github.com/technopark_database/internal/models"
	"strings"
	//"time"
)

type ForumPgRepository struct {
	db *sql.DB
}

func NewForumPgRepository(db *sql.DB) forum.ForumRepository {
	return &ForumPgRepository{db: db}
}

func (rep *ForumPgRepository) Insert(forum *models.Forum) error {
	tx, err := rep.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO forums(title, profile, slug) 
		VALUES($1, $2, $3)`, forum.Title, forum.User, forum.Slug)
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

func (rep *ForumPgRepository) InsertUserForum(nickname string, slug string) error {
	tx, err := rep.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO user_forum(nickname, slug) 
		VALUES($1, $2)`, nickname, slug)
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

func (rep *ForumPgRepository) UpdatePostsCount(slug string, count int) error {
	tx, err := rep.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE forums
		SET posts = posts + $1
		WHERE slug=$2`, count, slug)
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

func (rep *ForumPgRepository) SelectUserForum(nickname string, slug string) (string, string, error) {
	var dbNickname, dbSlug string
	err := rep.db.QueryRow(`
		SELECT nickname, slug
		FROM user_forum
		WHERE slug=$1 AND nickname=$2`,
		slug, nickname).Scan(
		&dbNickname, &dbSlug)
	if err != nil {
		return "", "", err
	}
	return dbNickname, dbSlug, nil
}

func (rep *ForumPgRepository) Select(slug string) (*models.Forum, error) {
	forum := &models.Forum{}
	err := rep.db.QueryRow(`
		SELECT title, profile, slug
		FROM forums
		WHERE slug = $1`, slug).Scan(
		&forum.Title, &forum.User, &forum.Slug)
	if err != nil {
		return nil, err
	}
	return forum, nil
}

func (rep *ForumPgRepository) SelectFull(slug string) (*models.Forum, error) {
	forum := &models.Forum{}
	err := rep.db.QueryRow(`
		SELECT title, profile, slug, posts, threads
		FROM forums
		WHERE slug = $1`, slug).Scan(
		&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
	if err != nil {
		return nil, err
	}
	return forum, nil
}

func (rep *ForumPgRepository) SelectThreads(forumSlug string,
	limit int, since string, desc bool) ([]*models.Thread, error) {
	query := `
		SELECT t.id, t.title, t.author, t.forum, t.message, t.votes, t.slug, t.created
		FROM forums f
		JOIN threads t on t.forum=f.slug
		WHERE forum = $1`

	var values []interface{}
	values = append(values, forumSlug)

	i := 2

	if since != "" {
		if desc {
			query = strings.Join([]string{query,
				"AND created<=$2",
			}, " ")
		} else {
			query = strings.Join([]string{query,
				"AND created>=$2",
			}, " ")
		}
		i++
		values = append(values, since)
	}

	query = strings.Join([]string{query, "ORDER BY created"}, " ")
	if desc {
		query = strings.Join([]string{query,
			"DESC",
		}, " ")
	}
	limitStr := fmt.Sprintf("LIMIT $%d", i)
	query = strings.Join([]string{query, limitStr}, " ")
	values = append(values, limit)

	rows, err := rep.db.Query(query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threads []*models.Thread
	for rows.Next() {
		thread := &models.Thread{}
		err := rows.Scan(&thread.ID, &thread.Title, &thread.Author, &thread.Forum,
			&thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
		if err != nil {
			return nil, err
		}
		threads = append(threads, thread)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return threads, nil
}

func (rep *ForumPgRepository) SelectUsers(slug string, limit int,
	since string, desc bool) ([]*models.User, error) {
	query := `
		SELECT u.id, u.nickname, u.fullname, u.about,
		       u.email
		FROM user_forum
		JOIN users u on u.nickname = user_forum.nickname
		WHERE slug=$1`
	var values []interface{}
	values = append(values, slug)

	i := 2
	if since != "" {
		if desc {
			query = strings.Join([]string{query, "AND u.nickname < $2"}, " ")
		} else {
			query = strings.Join([]string{query, "AND u.nickname > $2"}, " ")
		}
		values = append(values, since)
		i++
	}

	if desc {
		query = strings.Join([]string{query, "ORDER BY u.nickname DESC"}, " ")
	} else {
		query = strings.Join([]string{query, "ORDER BY u.nickname"}, " ")
	}
	if limit != 0 {
		query = strings.Join([]string{query,
			fmt.Sprintf("LIMIT $%d", i)}, " ")
		values = append(values, limit)
	}

	rows, err := rep.db.Query(query, values...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*models.User{}
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.ID, &user.Nickname,
			&user.Fullname, &user.About, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

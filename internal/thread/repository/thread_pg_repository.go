package repository

import (
	"context"
	"database/sql"
	"github.com/technopark_database/internal/models"
	"github.com/technopark_database/internal/thread"
)

type ThreadPgRepository struct {
	db *sql.DB
}

func NewThreadPgRepository(db *sql.DB) thread.ThreadRepository {
	return &ThreadPgRepository{db: db}
}

func (rep *ThreadPgRepository) Insert(thread *models.Thread) error {
	tx, err := rep.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	err = tx.QueryRow(`
		INSERT INTO threads(title, author, forum, message, votes, slug, created) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`,
		thread.Title, thread.Author, thread.Forum, thread.Message,
		thread.Votes, thread.Slug, thread.Created).Scan(&thread.ID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (rep *ThreadPgRepository) UpdateByID(thread *models.Thread) error {
	tx, err := rep.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = rep.db.Exec(`
		UPDATE threads
		SET title=$1,
		message=$2,
		slug=$3,
		votes=$4
		WHERE id=$5`, thread.Title, thread.Message, thread.Slug,
		thread.Votes, thread.ID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (rep *ThreadPgRepository) UpdateBySlug(thread *models.Thread) error {
	tx, err := rep.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = rep.db.Exec(`
		UPDATE threads
		SET title=$1,
		message=$2,
		votes=$3
		WHERE slug=$4`, thread.Title, thread.Message, thread.Votes, thread.Slug)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (rep *ThreadPgRepository) SelectByID(id uint64) (*models.Thread, error) {
	thread := &models.Thread{}
	err := rep.db.QueryRow(`
		SELECT id, title, author, forum, message, votes, slug, created
		FROM threads
		WHERE id=$1`, id).Scan(&thread.ID, &thread.Title, &thread.Author,
		&thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	if err != nil {
		return nil, err
	}
	return thread, nil
}

func (rep *ThreadPgRepository) SelectBySlug(slug string) (*models.Thread, error) {
	thread := &models.Thread{}
	err := rep.db.QueryRow(`
		SELECT id, title, author, forum, message, votes, slug, created
		FROM threads
		WHERE slug=$1`, slug).Scan(&thread.ID, &thread.Title, &thread.Author,
		&thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	if err != nil {
		return nil, err
	}
	return thread, nil
}

func (rep *ThreadPgRepository) SelectPostsByID(id uint64) ([]*models.Post, error) {
	rows, err := rep.db.Query(`
		SELECT p.id, p.parent, p.author, p.message,
		       p.isedited, p.forum, p.thread, p.created
		FROM threads t
		LEFT OUTER JOIN posts p on t.id = p.thread
		WHERE t.id=$1`, id)
	if err != nil {
		return nil, err
	}

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		err := rows.Scan(&post.Id, &post.Parent,
			&post.Author, &post.Message, &post.IsEdited,
			&post.Forum, &post.Thread, &post.Created)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (rep *ThreadPgRepository) SelectPostsBySlug(slug string) ([]*models.Post, error) {
	rows, err := rep.db.Query(`
		SELECT p.id, p.parent, p.author, p.message,
		       p.isedited, p.forum, p.thread, p.created
		FROM threads t
		LEFT OUTER JOIN posts p on t.id = p.thread
		WHERE t.id=$1`, slug)
	if err != nil {
		return nil, err
	}

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		err := rows.Scan(&post.Id, &post.Parent,
			&post.Author, &post.Message, &post.IsEdited,
			&post.Forum, &post.Thread, &post.Created)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

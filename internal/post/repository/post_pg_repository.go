package repository

import (
	"context"
	"database/sql"
	"github.com/technopark_database/internal/models"
	"github.com/technopark_database/internal/post"
)

type PostPgRepository struct {
	db *sql.DB
}

func NewPostPgRepository(db *sql.DB) post.PostRepository {
	return &PostPgRepository{db: db}
}

func (rep *PostPgRepository) InsertMany(posts []*models.Post) error {
	tx, err := rep.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO posts(parent, author, message, 
		isedited, forum, thread, created)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`)
	if err != nil {
		return err
	}

	for _, post := range posts {
		_, err := stmt.Exec(post.Parent, post.Author, post.Message,
			post.IsEdited, post.Forum, post.Thread, post.Created)
		if err != nil {
			return err
		}
	}

	if _, err = stmt.Exec(); err != nil {
		return err
	}
	if err = stmt.Close(); err != nil {
		return err
	}

	return nil
}

func (rep *PostPgRepository) Update(post *models.Post) error {
	tx, err := rep.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE posts
		SET message=$1, isedited=true
		WHERE id=$2`, post.Message, post.Id)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (rep *PostPgRepository) SelectByID(id uint64) (*models.Post, error) {
	post := &models.Post{}
	err := rep.db.QueryRow(`
		SELECT id, parent, author, message, isedited, forum, thread, created
		FROM posts
		WHERE id=$1`, id).Scan(&post.Id, &post.Parent, &post.Author, &post.Message,
		&post.IsEdited, &post.Forum, &post.Thread, &post.Created)
	if err != nil {
		return nil, err
	}
	return post, err
}
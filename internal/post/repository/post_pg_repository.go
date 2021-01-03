package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/technopark_database/internal/helpers/gears"
	"github.com/technopark_database/internal/models"
	"github.com/technopark_database/internal/post"
	"strings"
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

	stmt, err := rep.db.Prepare(`
		INSERT INTO posts(parent, author, message, 
		isedited, forum, thread, created)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`)
	if err != nil {
		return err
	}

	for _, post := range posts {
		err := stmt.QueryRow(post.Parent, post.Author, post.Message,
			post.IsEdited, post.Forum, post.Thread, post.Created).
			Scan(&post.ID)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				logrus.Info(rollbackErr)
			}
			return err
		}
	}

	if err = stmt.Close(); err != nil {
		return err
	}

	if err := tx.Commit(); tx != nil {
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
		WHERE id=$2`, post.Message, post.ID)
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
		WHERE id=$1`, id).Scan(&post.ID, &post.Parent, &post.Author, &post.Message,
		&post.IsEdited, &post.Forum, &post.Thread, &post.Created)
	if err != nil {
		return nil, err
	}
	return post, err
}

func (rep *PostPgRepository) selectPostsTree(threadID uint64, since uint64,
	pagination *models.Pagination) ([]*models.Post, error) {
	var values []interface{}

	selectQuery := `
		SELECT id, parent, author, message, isedited, forum, thread, created
		FROM posts
		WHERE thread=$1`
	values = append(values, threadID)

	var sortQuery string
	if pagination.Desc {
		sortQuery = "ORDER BY path DESC"
	} else {
		sortQuery = "ORDER BY path"
	}

	var pgntQuery string
	if pagination.Limit != 0 {
		pgntQuery = "LIMIT $2"
		values = append(values, pagination.Limit)
	}

	var filterQuery string
	if since != 0 {
		ind := len(values) + 1
		subSelectQuery := fmt.Sprintf("(SELECT path FROM posts WHERE id=$%d)", ind)

		var subFilterQuery string
		if pagination.Desc {
			subFilterQuery = "AND path <"
		} else {
			subFilterQuery = "AND path >"
		}

		filterQuery = strings.Join([]string{
			subFilterQuery,
			subSelectQuery,
		}, " ")
		values = append(values, since)
	}

	resultQuery := strings.Join([]string{
		selectQuery,
		filterQuery,
		sortQuery,
		pgntQuery,
	}, " ")

	rows, err := rep.db.Query(resultQuery, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		err := rows.Scan(&post.ID, &post.Parent, &post.Author, &post.Message, &post.IsEdited,
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

func getSelectParentsQuery(threadID uint64,
	since uint64, pgnt *models.Pagination) (string, []interface{}) {

	var values []interface{}

	selectQuery := `
		SELECT id
		FROM posts
		WHERE thread=$1
		AND parent=0`
	values = append(values, threadID)

	var sortQuery string
	if pgnt.Desc {
		sortQuery = "ORDER BY id DESC"
	} else {
		sortQuery = "ORDER BY id"
	}

	var pgntQuery string
	if pgnt.Limit != 0 {
		pgntQuery = "LIMIT $2"
		values = append(values, pgnt.Limit)
	}

	var subSelectQuery string
	if since != 0 {
		subSelectQuery = `
		SELECT path[1]
		FROM posts
		WHERE id=$3`

		var filterQuery string
		if pgnt.Desc {
			filterQuery = "AND path[1] <"
		} else {
			filterQuery = "AND path[1] >"
		}
		subSelectQuery = fmt.Sprintf("%s (%s)", filterQuery, subSelectQuery)
		values = append(values, since)
	}

	resultQuery := strings.Join([]string{
		selectQuery,
		subSelectQuery,
		sortQuery,
		pgntQuery,
	}, " ")

	return resultQuery, values
}

func (rep *PostPgRepository) selectPostsParentTree(threadID uint64,
	since uint64, pgnt *models.Pagination) ([]*models.Post, error) {
	subSelectQuery, values := getSelectParentsQuery(threadID, since, pgnt)

	selectQuery := `
		SELECT id, parent, author, message, isedited, forum, thread, created
		FROM posts
		WHERE path[1] IN`

	var sortQuery string
	if pgnt.Desc {
		sortQuery = "ORDER BY path[1] DESC, path, id"
	} else {
		sortQuery = "ORDER BY path, id"
	}

	resultQuery := strings.Join([]string{
		selectQuery,
		fmt.Sprintf("(%s)", subSelectQuery),
		sortQuery,
	}, " ")

	rows, err := rep.db.Query(resultQuery, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		err := rows.Scan(&post.ID, &post.Parent, &post.Author, &post.Message, &post.IsEdited,
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

func (rep *PostPgRepository) selectPostsFlat(threadID uint64, since uint64,
	pagination *models.Pagination) ([]*models.Post, error) {
	query := `
		SELECT id, parent, author, message,
       			isedited, forum, thread, created
		FROM posts
		WHERE thread=$1`
	var values []interface{}
	values = append(values, threadID)

	query, values = gears.AddPagination(query, values, pagination, since, 2)

	rows, err := rep.db.Query(query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		err := rows.Scan(&post.ID, &post.Parent,
			&post.Author, &post.Message, &post.IsEdited, &post.Forum,
			&post.Thread, &post.Created)
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

func (rep *PostPgRepository) SelectPosts(threadID uint64, sort string, since uint64, pagination *models.Pagination) ([]*models.Post, error) {
	switch sort {
	case "tree":
		return rep.selectPostsTree(threadID, since, pagination)
	case "parent_tree":
		return rep.selectPostsParentTree(threadID, since, pagination)
	default:
		return rep.selectPostsFlat(threadID, since, pagination)
	}
}

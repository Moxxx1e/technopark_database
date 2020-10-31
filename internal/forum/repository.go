package forum

import "github.com/technopark_database/internal/models"

type ForumRepository interface {
	Insert(forum *models.Forum) error
	Select(slug string) (*models.Forum, error)
	SelectUsers(slug string) ([]*models.User, error)
	SelectThreads(slug string) ([]*models.Thread, error)
}

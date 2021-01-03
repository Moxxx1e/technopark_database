package forum

import (
	"github.com/technopark_database/internal/models"
)

type ForumRepository interface {
	Insert(forum *models.Forum) error
	InsertUserForum(nickname string, slug string) error
	Select(slug string) (*models.Forum, error)
	SelectCounts(slug string) (int, int, error)
	SelectFull(slug string) (*models.Forum, error)
	SelectUserForum(nickname string, slug string) (string, string, error)
	SelectUsers(slug string, limit int, since string, desc bool) ([]*models.User, error)
	SelectThreads(slug string, limit int, since string, desc bool) ([]*models.Thread, error)
}

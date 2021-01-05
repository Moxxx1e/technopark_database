package thread

import "github.com/technopark_database/internal/models"

type ThreadRepository interface {
	Insert(thread *models.Thread) error
	UpdateByID(thread *models.Thread) error
	UpdateBySlug(thread *models.Thread) error
	SelectByID(id uint64) (*models.Thread, error)
	SelectBySlug(slug string) (*models.Thread, error)
	SelectPostsByID(id uint64) ([]*models.Post, error)
	SelectPostsBySlug(slug string) ([]*models.Post, error)
}

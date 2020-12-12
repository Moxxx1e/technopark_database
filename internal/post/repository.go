package post

import "github.com/technopark_database/internal/models"

type PostRepository interface {
	InsertMany(posts []*models.Post) error
	Update(post *models.Post) error
	SelectByID(id uint64) (*models.Post, error)
}

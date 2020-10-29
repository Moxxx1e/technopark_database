package user

import "github.com/technopark_database/internal/models"

type UserRepository interface {
	Insert(user *models.User) error
	Select(nickname string) (*models.User, error)
	Update(nickname string, user *models.User) error
	SelectByEmail(email string) (*models.User, error)
}

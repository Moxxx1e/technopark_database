package user

import "github.com/technopark_database/internal/models"

type UserUseCase interface {
	Create(user *models.User) error
	UpdateUserInfo(nickname string, user *models.User) error
	GetUserInfo(nickname string) (*models.User, error)
}

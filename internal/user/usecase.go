package user

import (
	"github.com/technopark_database/internal/helpers/errors"
	"github.com/technopark_database/internal/models"
)

type UserUseCase interface {
	Create(user *models.User) *errors.Error
	UpdateUserInfo(nickname string, user *models.User) *errors.Error
	GetUserInfo(nickname string) (*models.User, *errors.Error)
}

package user

import (
	"github.com/technopark_database/internal/helpers/errors"
	"github.com/technopark_database/internal/models"
)

type UserUseCase interface {
	Create(user *models.User) ([]*models.User, *errors.Error)
	Change(user *models.User) (*models.User, *errors.Error)
	GetUserInfo(nickname string) (*models.User, *errors.Error)
	IsExist(nickname string) (bool, *errors.Error)
}

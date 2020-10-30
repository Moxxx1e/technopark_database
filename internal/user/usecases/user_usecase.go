package usecases

import (
	"database/sql"
	"github.com/technopark_database/internal/models"
	"github.com/technopark_database/internal/user"
)

type UserUseCase struct {
	rep user.UserRepository
}

func NewUserUseCase(repository user.UserRepository) *UserUseCase {
	return &UserUseCase{rep: repository}
}

func (uc *UserUseCase) Create(user *models.User) error {
	return uc.rep.Insert(user)
}

func (uc *UserUseCase) UpdateUserInfo(nickname string, user *models.User) error {
	return uc.rep.Update(nickname, user)
}

func (uc *UserUseCase) GetUserInfo(nickname string) (*models.User, error) {
	info, err :=  uc.rep.Select(nickname)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return info, nil
}

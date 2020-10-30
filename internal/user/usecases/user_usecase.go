package usecases

import (
	"database/sql"
	"github.com/technopark_database/internal/consts"
	"github.com/technopark_database/internal/helpers/errors"
	"github.com/technopark_database/internal/models"
	"github.com/technopark_database/internal/user"
)

type UserUseCase struct {
	rep user.UserRepository
}

func NewUserUseCase(repository user.UserRepository) *UserUseCase {
	return &UserUseCase{rep: repository}
}

func (uc *UserUseCase) Create(user *models.User) *errors.Error {
	err := uc.rep.Insert(user)
	if err != nil {
		return errors.New(consts.CodeInternalServerError, err)
	}
	return nil
}

func (uc *UserUseCase) UpdateUserInfo(nickname string, user *models.User) *errors.Error {
	err := uc.rep.Update(nickname, user)
	if err != nil {
		return errors.New(consts.CodeInternalServerError, err)
	}
	return nil
}

func (uc *UserUseCase) GetUserInfo(nickname string) (*models.User, *errors.Error) {
	info, err :=  uc.rep.Select(nickname)
	if err == sql.ErrNoRows {
		return nil, errors.Get(consts.CodeUserDoesNotExist)
	} else if err != nil {
		return nil, errors.New(consts.CodeInternalServerError, err)
	}
	return info, nil
}

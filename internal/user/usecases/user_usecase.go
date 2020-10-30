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
	if err := uc.IsUserConflicts(user.Nickname, user); err != nil {
		return err
	}

	err := uc.rep.Insert(user)
	if err != nil {
		return errors.New(consts.CodeInternalServerError, err)
	}
	return nil
}

func (uc *UserUseCase) UpdateUserInfo(nickname string, userWithNewInfo *models.User) *errors.Error {
	dbUser, customError := uc.IsEmailExists(userWithNewInfo.Email)
	if customError != nil {
		// if updated email == previous email
		if dbUser.Nickname == nickname {
			return nil
		}
		return customError
	}

	err := uc.rep.Update(nickname, userWithNewInfo)
	if err != nil {
		return errors.New(consts.CodeInternalServerError, err)
	}
	return nil
}

func (uc *UserUseCase) GetUserInfo(nickname string) (*models.User, *errors.Error) {
	info, err := uc.rep.Select(nickname)
	if err == sql.ErrNoRows {
		return nil, errors.Get(consts.CodeUserDoesNotExist)
	} else if err != nil {
		return nil, errors.New(consts.CodeInternalServerError, err)
	}
	return info, nil
}

func (uc *UserUseCase) IsNicknameExists(nickname string) (*models.User, *errors.Error) {
	dbUser, err := uc.rep.Select(nickname)
	if err != sql.ErrNoRows {
		return dbUser, errors.Get(consts.CodeUserNicknameConflicts)
	}
	return nil, nil
}

func (uc *UserUseCase) IsEmailExists(email string) (*models.User, *errors.Error) {
	dbUser, err := uc.rep.SelectByEmail(email)
	if err != sql.ErrNoRows {
		return dbUser, errors.Get(consts.CodeUserEmailConflicts)
	}
	return nil, nil
}

func (uc *UserUseCase) IsUserConflicts(nickname string, user *models.User) *errors.Error {
	if _, err := uc.IsNicknameExists(nickname); err != nil {
		return err
	}
	if _, err := uc.IsEmailExists(user.Email); err != nil {
		return err
	}
	return nil
}

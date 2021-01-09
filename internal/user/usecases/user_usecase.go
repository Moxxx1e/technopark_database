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

func (uc *UserUseCase) IsExist(nickname string) (bool, *errors.Error) {
	_, err := uc.rep.Select(nickname)
	if err == sql.ErrNoRows {
		return false, errors.Get(consts.CodeUserDoesNotExist)
	} else if err != nil {
		return false, errors.New(consts.CodeInternalServerError, err)
	}
	return true, nil
}

func NewUserUseCase(repository user.UserRepository) user.UserUseCase {
	return &UserUseCase{rep: repository}
}

func (uc *UserUseCase) Create(user *models.User) ([]*models.User, *errors.Error) {
	nicknameUserConflict, customErr := uc.isNicknameExist(user.Nickname)
	if customErr != nil {
		return nil, customErr
	}
	conflictedUsers := []*models.User{}
	if nicknameUserConflict != nil {
		conflictedUsers = append(conflictedUsers, nicknameUserConflict)
	}

	emailUserConflict, customErr := uc.IsEmailAlreadyExist(user.Email)
	if customErr != nil {
		return nil, customErr
	}
	if emailUserConflict != nil {
		if len(conflictedUsers) == 0 {
			conflictedUsers = append(conflictedUsers, emailUserConflict)
		} else if len(conflictedUsers) == 1 &&
			conflictedUsers[0].ID != emailUserConflict.ID {
			conflictedUsers = append(conflictedUsers, emailUserConflict)
		}
	}

	if len(conflictedUsers) != 0 {
		return conflictedUsers, errors.Get(consts.CodeUserEmailConflicts)
	}

	err := uc.rep.Insert(user)
	if err != nil {
		return nil, errors.New(consts.CodeInternalServerError, err)
	}
	return []*models.User{user}, nil
}

func (uc *UserUseCase) isNicknameExist(nickname string) (*models.User, *errors.Error) {
	dbUser, err := uc.rep.Select(nickname)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Get(consts.CodeInternalServerError)
	}

	return dbUser, nil
}

// returns user with this email
// if user exists
// otherwise return nil
func (uc *UserUseCase) IsEmailAlreadyExist(email string) (*models.User, *errors.Error) {
	dbUser, err := uc.rep.SelectByEmail(email)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Get(consts.CodeInternalServerError)
	}

	return dbUser, nil
}

func (uc *UserUseCase) Change(user *models.User) (*models.User, *errors.Error) {
	dbUser, customErr := uc.GetUserInfo(user.Nickname)
	if customErr != nil {
		return nil, customErr
	}

	if user.Email == "" {
		user.Email = dbUser.Email
	} else {
		emailUser, customErr := uc.IsEmailAlreadyExist(user.Email)
		if customErr != nil {
			return nil, customErr
		}
		if emailUser != nil {
			return emailUser, errors.Get(consts.CodeUserEmailConflicts)
		}
	}

	if user.Fullname == "" {
		user.Fullname = dbUser.Fullname
	}
	if user.About == "" {
		user.About = dbUser.About
	}

	err := uc.rep.Update(user)
	if err != nil {
		return nil, errors.New(consts.CodeInternalServerError, err)
	}
	return user, nil
}

func (uc *UserUseCase) GetUserInfo(nickname string) (*models.User, *errors.Error) {
	dbUser, err := uc.rep.Select(nickname)
	if err == sql.ErrNoRows {
		return dbUser, errors.Get(consts.CodeUserDoesNotExist)
	} else if err != nil {
		return nil, errors.New(consts.CodeInternalServerError, err)
	}
	return dbUser, nil
}

func removeDuplicateValues(stringSlice []string) []string {
	keys := make(map[string]bool)
	result := []string{}

	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			result = append(result, entry)
		}
	}
	return result
}

func (uc *UserUseCase) CheckNicknames(nicknames []string) *errors.Error {
	uniqNicknames := removeDuplicateValues(nicknames)

	countNicknames, err := uc.rep.SelectCountNicknames(uniqNicknames)
	if err != nil {
		return errors.New(consts.CodeInternalServerError, err)
	}
	if countNicknames != len(uniqNicknames) {
		return errors.Get(consts.CodeUserDoesNotExist)
	}

	return nil
}

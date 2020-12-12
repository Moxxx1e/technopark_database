package usecases

import (
	"database/sql"
	"github.com/technopark_database/internal/consts"
	"github.com/technopark_database/internal/forum"
	"github.com/technopark_database/internal/helpers/errors"
	"github.com/technopark_database/internal/models"
	"github.com/technopark_database/internal/user"
)

type ForumUseCase struct {
	rep         forum.ForumRepository
	userUseCase user.UserUseCase
}

func NewForumUseCase(rep forum.ForumRepository, userUseCase user.UserUseCase) forum.ForumUseCase {
	return &ForumUseCase{rep: rep, userUseCase: userUseCase}
}

func (uc *ForumUseCase) Create(forum *models.Forum) (*models.Forum, *errors.Error) {
	dbUser, customErr := uc.userUseCase.GetUserInfo(forum.User)
	if customErr != nil {
		return nil, customErr
	}
	forum.User = dbUser.Nickname

	dbForum, err := uc.rep.Select(forum.Slug)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.New(consts.CodeInternalServerError, err)
	}
	if dbForum != nil {
		return dbForum, errors.Get(consts.CodeForumAlreadyExist)
	}

	if err := uc.rep.Insert(forum); err != nil {
		return nil, errors.New(consts.CodeInternalServerError, err)
	}
	return forum, nil
}

func (uc *ForumUseCase) GetDetails(slug string) (*models.Forum, *errors.Error) {
	forum, err := uc.rep.Select(slug)
	if err == sql.ErrNoRows {
		return nil, errors.Get(consts.CodeForumDoesNotExist)
	} else if err != nil {
		return nil, errors.New(consts.CodeInternalServerError, err)
	}

	return forum, nil
}

func (uc *ForumUseCase) GetUsers(slug string, pagination *models.Pagination) ([]*models.User, *errors.Error) {
	// TODO: после поста
	panic("implement me!")
}

func (uc *ForumUseCase) GetThreads(slug string, pagination *models.Pagination) ([]*models.Thread, *errors.Error) {
	if pagination.Limit == 0 {
		pagination.Limit = 100
	}

	if _, err := uc.GetDetails(slug); err != nil {
		return nil, err
	}

	threads, err := uc.rep.SelectThreads(slug, pagination.Limit, pagination.Since, pagination.Desc)
	if err != nil {
		return nil, errors.New(consts.CodeInternalServerError, err)
	}
	if len(threads) == 0 {
		return []*models.Thread{}, nil
	}
	return threads, nil
}

func (uc *ForumUseCase) IsExist(slug string) (*models.Forum, *errors.Error) {
	dbForum, err := uc.rep.Select(slug)
	if err == sql.ErrNoRows {
		return dbForum, errors.Get(consts.CodeForumDoesNotExist)
	} else if err != nil {
		return nil, errors.New(consts.CodeInternalServerError, err)
	}

	return nil, nil
}

//TODO: isUserExist

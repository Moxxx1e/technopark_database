package usecases

import (
	"database/sql"
	"github.com/technopark_database/internal/consts"
	"github.com/technopark_database/internal/helpers/errors"
	"github.com/technopark_database/internal/models"
	"github.com/technopark_database/internal/post"
)

type PostUseCase struct {
	rep post.PostRepository
}

func NewPostUseCase(rep post.PostRepository) post.PostUseCase {
	return &PostUseCase{rep: rep}
}

func (p PostUseCase) CreateMany(posts []*models.Post) *errors.Error {
	panic("implement me")
}

func (p PostUseCase) ChangeByID(id uint64, message string) (*models.Post, *errors.Error) {
	post, err := p.rep.SelectByID(id)
	if err == sql.ErrNoRows {
		return nil, errors.Get(consts.CodePostDoesNotExist)
	} else if err != nil {
		return nil, errors.New(consts.CodeInternalServerError, err)
	}

	post.IsEdited = true
	post.Message = message
	if err := p.rep.Update(post); err != nil {
		return nil, errors.New(consts.CodeInternalServerError, err)
	}

	return post, nil
}

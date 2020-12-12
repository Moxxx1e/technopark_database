package usecases

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"github.com/technopark_database/internal/consts"
	"github.com/technopark_database/internal/helpers/errors"
	"github.com/technopark_database/internal/models"
	"github.com/technopark_database/internal/post"
	"github.com/technopark_database/internal/thread"
	"strconv"
	"time"
)

type PostUseCase struct {
	threadUseCase thread.ThreadUsecase
	rep           post.PostRepository
}

func NewPostUseCase(threadUseCase thread.ThreadUsecase, rep post.PostRepository) post.PostUseCase {
	return &PostUseCase{rep: rep, threadUseCase: threadUseCase}
}

func (p *PostUseCase) CreateMany(slugOrID string, posts []*models.Post) ([]*models.Post, *errors.Error) {
	if len(posts) == 0 {
		return []*models.Post{}, nil
	}

	id, err := strconv.ParseUint(slugOrID, 10, 64)
	var thread *models.Thread
	var customErr *errors.Error
	if err != nil {
		logrus.Info("slug: ", slugOrID)
		thread, customErr = p.threadUseCase.GetBySlug(slugOrID)
		if customErr != nil {
			return nil, customErr
		}
	} else {
		thread, customErr = p.threadUseCase.GetByID(id)
		if customErr != nil {
			return nil, customErr
		}
	}

	createdTime := time.Now()
	existedPosts, customErr := p.threadUseCase.GetPostsByID(id)
	if customErr != nil {
		return nil, customErr
	}

	for _, post := range posts {
		post.Forum = thread.Forum
		post.Created = createdTime
		post.Thread = thread.ID

		if post.Parent != 0 {
			parentFlag := false
			for _, existedPost := range existedPosts {
				if post.Parent == existedPost.ID {
					parentFlag = true
					break
				}
			}
			if parentFlag != true {
				return nil, errors.Get(consts.CodeParentPostDoesNotExistInThread)
			}
		}
	}

	err = p.rep.InsertMany(posts)
	if err != nil {
		return nil, errors.New(consts.CodeInternalServerError, err)
	}
	return posts, nil
}

func (p *PostUseCase) ChangeByID(id uint64, message string) (*models.Post, *errors.Error) {
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

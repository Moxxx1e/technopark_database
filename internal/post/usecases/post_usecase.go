package usecases

import (
	"database/sql"
	"github.com/technopark_database/internal/consts"
	"github.com/technopark_database/internal/forum"
	"github.com/technopark_database/internal/helpers/errors"
	"github.com/technopark_database/internal/models"
	"github.com/technopark_database/internal/post"
	"github.com/technopark_database/internal/thread"
	"github.com/technopark_database/internal/user"
	"strconv"
)

type PostUseCase struct {
	threadUseCase thread.ThreadUsecase
	rep           post.PostRepository
	forumUseCase  forum.ForumUseCase
	userUseCase   user.UserUseCase
}

func NewPostUseCase(threadUseCase thread.ThreadUsecase,
	rep post.PostRepository,
	forumUseCase forum.ForumUseCase,
	userUseCase user.UserUseCase) post.PostUseCase {
	return &PostUseCase{
		rep:           rep,
		threadUseCase: threadUseCase,
		forumUseCase:  forumUseCase,
		userUseCase:   userUseCase,
	}
}

func (uc *PostUseCase) CreateMany(slugOrID string, posts []*models.Post) ([]*models.Post, *errors.Error) {
	thread, customErr := uc.GetThreadBySlugOrID(slugOrID)
	if customErr != nil {
		return nil, customErr
	}

	if len(posts) == 0 {
		return []*models.Post{}, nil
	}

	var nicknames []string
	for _, post := range posts {
		nicknames = append(nicknames, post.Author)
		post.Forum = thread.Forum
		//post.Created = createdTime
		post.Thread = thread.ID
	}

	customErr = uc.userUseCase.CheckNicknames(nicknames)
	if customErr != nil {
		return nil, customErr
	}

	err := uc.rep.InsertMany(posts)
	if err != nil {
		if err.Error() == "pq: Parent post does not exist in thread" {
			return nil, errors.Get(consts.CodeParentPostDoesNotExistInThread)
		}
		return nil, errors.New(consts.CodeInternalServerError, err)
	}

	return posts, nil
}

func (uc *PostUseCase) GetThreadBySlugOrID(slugOrID string) (*models.Thread, *errors.Error) {
	id, err := strconv.ParseUint(slugOrID, 10, 64)
	var thread *models.Thread
	var customErr *errors.Error
	if err != nil {
		thread, customErr = uc.threadUseCase.GetBySlug(slugOrID)
		if customErr != nil {
			return nil, customErr
		}
	} else {
		thread, customErr = uc.threadUseCase.GetByID(id)
		if customErr != nil {
			return nil, customErr
		}
	}
	return thread, nil
}

func (uc *PostUseCase) GetPosts(slugOrID string, sort string, since uint64,
	pagination *models.Pagination) ([]*models.Post, *errors.Error) {
	thread, customErr := uc.GetThreadBySlugOrID(slugOrID)
	if customErr != nil {
		return nil, customErr
	}

	posts, err := uc.rep.SelectPosts(thread.ID, sort, since, pagination)
	if err != nil {
		return nil, errors.New(consts.CodeInternalServerError, err)
	}
	if posts == nil {
		return []*models.Post{}, nil
	}

	return posts, nil
}

func (uc *PostUseCase) ChangeByID(id uint64, message string) (*models.Post, *errors.Error) {
	post, err := uc.rep.SelectByID(id)
	if err == sql.ErrNoRows {
		return nil, errors.Get(consts.CodePostDoesNotExist)
	} else if err != nil {
		return nil, errors.New(consts.CodeInternalServerError, err)
	}

	if message == "" || message == post.Message {
		return post, nil
	}

	post.IsEdited = true
	post.Message = message
	if err := uc.rep.Update(post); err != nil {
		return nil, errors.New(consts.CodeInternalServerError, err)
	}

	return post, nil
}

func (uc *PostUseCase) GetPostInfo(id uint64, related *models.Related) (*models.PostDetails, *errors.Error) {
	postDetails := &models.PostDetails{}

	post, err := uc.rep.SelectByID(id)
	if err == sql.ErrNoRows {
		return nil, errors.Get(consts.CodePostDoesNotExist)
	} else if err != nil {
		return nil, errors.New(consts.CodeInternalServerError, err)
	}
	postDetails.Post = post

	if related.User {
		user, customErr := uc.userUseCase.GetUserInfo(post.Author)
		if customErr != nil {
			return nil, customErr
		}
		postDetails.Author = user
	}

	if related.Forum {
		forum, customErr := uc.forumUseCase.GetFullDetails(post.Forum)
		if customErr != nil {
			return nil, customErr
		}
		postDetails.Forum = forum
	}

	if related.Thread {
		thread, customErr := uc.threadUseCase.GetByID(post.Thread)
		if customErr != nil {
			return nil, customErr
		}
		postDetails.Thread = thread
	}

	return postDetails, nil
}

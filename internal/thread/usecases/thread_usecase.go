package usecases

import (
	"database/sql"
	"github.com/technopark_database/internal/consts"
	"github.com/technopark_database/internal/forum"
	"github.com/technopark_database/internal/helpers/errors"
	"github.com/technopark_database/internal/models"
	"github.com/technopark_database/internal/thread"
	"github.com/technopark_database/internal/user"
	"github.com/technopark_database/internal/vote"
)

type ThreadUseCase struct {
	rep          thread.ThreadRepository
	userUseCase  user.UserUseCase
	forumUseCase forum.ForumUseCase
	voteUseCase  vote.VoteUseCase
}

func NewThreadUseCase(rep thread.ThreadRepository, userUseCase user.UserUseCase,
	forumUseCase forum.ForumUseCase, voteUseCase vote.VoteUseCase) thread.ThreadUsecase {
	return &ThreadUseCase{rep: rep,
		userUseCase:  userUseCase,
		forumUseCase: forumUseCase,
		voteUseCase:  voteUseCase}
}

func (th *ThreadUseCase) Create(thread *models.Thread) (*models.Thread, *errors.Error) {
	existedForum, customErr := th.forumUseCase.GetDetails(thread.Forum)
	if customErr != nil {
		return nil, customErr
	}
	thread.Forum = existedForum.Slug

	author, customErr := th.userUseCase.GetUserInfo(thread.Author)
	if customErr != nil {
		return nil, customErr
	}
	thread.Author = author.Nickname

	if thread.Slug != "" {
		existedThread, customErr := th.GetBySlug(thread.Slug)
		if customErr != nil && customErr != errors.Get(consts.CodeThreadDoesNotExist) {
			return nil, customErr
		} else if existedThread != nil {
			return existedThread, errors.Get(consts.CodeThreadAlreadyExist)
		}
	}

	err := th.rep.Insert(thread)
	if err != nil {
		return nil, errors.New(consts.CodeInternalServerError, err)
	}

	return thread, nil
}

func (th *ThreadUseCase) CreateVote(thread *models.Thread,
	user *models.User, vote int) *models.Vote {
	voteModel := &models.Vote{
		ThreadID: thread.ID,
		UserID:   user.ID,
		Likes: vote == 1,
	}
	return voteModel
}

func (th *ThreadUseCase) CreateVoteByID(id uint64, nickname string, vote int) (*models.Thread, *errors.Error) {
	thread, customErr := th.GetByID(id)
	if customErr != nil {
		return nil, customErr
	}

	user, customErr := th.userUseCase.GetUserInfo(nickname)
	if customErr != nil {
		return nil, customErr
	}

	voteModel := th.CreateVote(thread, user, vote)
	// don't need to update votes field in thread
	// because there is a trigger in db
	_, customErr = th.voteUseCase.Create(voteModel)
	if customErr != nil {
		return nil, customErr
	}
	// TODO: неоптимально
	thread, customErr = th.GetByID(id)
	if customErr != nil {
		return nil, customErr
	}

	return thread, nil
}

func (th *ThreadUseCase) CreateVoteBySlug(slug string, nickname string, vote int) (*models.Thread, *errors.Error) {
	thread, customErr := th.GetBySlug(slug)
	if customErr != nil {
		return nil, customErr
	}

	user, customErr := th.userUseCase.GetUserInfo(nickname)
	if customErr != nil {
		return nil, customErr
	}

	voteModel := th.CreateVote(thread, user, vote)
	// don't need to update votes field in thread
	// because there is a trigger in db
	_, customErr = th.voteUseCase.Create(voteModel)

	// TODO: неоптимально
	thread, customErr = th.GetBySlug(slug)
	if customErr != nil {
		return nil, customErr
	}


	return thread, nil
}

func (th *ThreadUseCase) IsForumExist(title string) (bool, *errors.Error) {
	panic("")
}

func (th *ThreadUseCase) ChangeByID(id uint64, title, message string) (*models.Thread, *errors.Error) {
	thread, customErr := th.GetByID(id)
	if customErr != nil {
		return nil, customErr
	}
	thread.Title, thread.Message = title, message

	err := th.rep.UpdateByID(thread)
	if err != nil {
		return nil, errors.New(consts.CodeInternalServerError, err)
	}
	return thread, nil
}

func (th *ThreadUseCase) ChangeBySlug(slug string, title, message string) (*models.Thread, *errors.Error) {
	thread, customErr := th.GetBySlug(slug)
	if customErr != nil {
		return nil, customErr
	}
	thread.Title, thread.Message = title, message

	err := th.rep.UpdateBySlug(thread)
	if err != nil {
		return nil, errors.New(consts.CodeInternalServerError, err)
	}
	return thread, nil
}

func (th *ThreadUseCase) GetByID(id uint64) (*models.Thread, *errors.Error) {
	thread, err := th.rep.SelectByID(id)
	if err == sql.ErrNoRows {
		return nil, errors.Get(consts.CodeThreadDoesNotExist)
	} else if err != nil {
		return nil, errors.New(consts.CodeInternalServerError, err)
	}
	return thread, nil
}

func (th *ThreadUseCase) GetBySlug(slug string) (*models.Thread, *errors.Error) {
	thread, err := th.rep.SelectBySlug(slug)
	if err == sql.ErrNoRows {
		return nil, errors.Get(consts.CodeThreadDoesNotExist)
	} else if err != nil {
		return nil, errors.New(consts.CodeInternalServerError, err)
	}
	return thread, nil
}

func (th *ThreadUseCase) GetPostsByID(id uint64) ([]*models.Post, *errors.Error) {
	posts, err := th.rep.SelectPostsByID(id)
	if err == sql.ErrNoRows {
		return []*models.Post{}, nil
	} else if err != nil {
		return nil, errors.New(consts.CodeInternalServerError, err)
	}
	return posts, nil
}

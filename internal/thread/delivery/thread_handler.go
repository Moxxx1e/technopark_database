package delivery

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/technopark_database/internal/helpers/errors"
	"github.com/technopark_database/internal/models"
	"github.com/technopark_database/internal/thread"
	reader "github.com/technopark_database/tools/requestReader"
	"net/http"
	"strconv"
	"time"
)

type ThreadHandler struct {
	threadUseCase thread.ThreadUsecase
}

type Message struct {
	Message string `json:"message"`
}

func NewThreadHandler(threadUseCase thread.ThreadUsecase) *ThreadHandler {
	return &ThreadHandler{threadUseCase: threadUseCase}
}

func (th *ThreadHandler) Configure(e *echo.Echo) {
	// Проверено
	e.POST("/api/forum/:forum_slug/create", th.CreateThreadHandler())
	e.GET("/api/thread/:slug_or_id/details", th.GetDetailsHandler())
	e.POST("/api/thread/:slug_or_id/vote", th.VoteHandler())
	e.POST("/api/thread/:slug_or_id/details", th.ChangeThreadHandler())
}

func (th *ThreadHandler) CreateThreadHandler() echo.HandlerFunc {
	type Request struct {
		Title   string    `json:"title"`
		Author  string    `json:"author"`
		Message string    `json:"message"`
		Slug    string    `json:"slug" validate:"omitempty"`
		Created time.Time `json:"created"`
	}
	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logrus.Error(err.DebugMessage)
			return cntx.JSON(err.HTTPCode, Message{Message: err.UserMessage})
		}

		forumSlug := cntx.Param("forum_slug")


		thread := &models.Thread{
			Title:   req.Title,
			Author:  req.Author,
			Message: req.Message,
			Forum:   forumSlug,
			Slug:    req.Slug,
			Created: req.Created,
		}
		err := th.threadUseCase.Create(thread)
		if err != nil {
			logrus.Error(err.DebugMessage)
			return cntx.JSON(err.HTTPCode, Message{Message: err.UserMessage})
		}

		return cntx.JSON(http.StatusCreated, thread)
	}
}

func (th *ThreadHandler) GetDetailsHandler() echo.HandlerFunc {
	return func(cntx echo.Context) error {
		slugOrID := cntx.Param("slug_or_id")

		threadDetails := &models.Thread{}
		var customErr *errors.Error

		id, err := strconv.ParseUint(slugOrID, 10, 64)
		if err != nil {
			threadDetails, customErr = th.threadUseCase.GetBySlug(slugOrID)
			if customErr != nil {
				logrus.Error(customErr.DebugMessage)
				return cntx.JSON(customErr.HTTPCode, Message{customErr.UserMessage})
			}
		} else {
			threadDetails, customErr = th.threadUseCase.GetByID(id)
			if customErr != nil {
				logrus.Error(customErr.DebugMessage)
				return cntx.JSON(customErr.HTTPCode, Message{customErr.UserMessage})
			}
		}

		return cntx.JSON(http.StatusOK, threadDetails)
	}
}

func (th *ThreadHandler) VoteHandler() echo.HandlerFunc {
	type Request struct {
		Nickname string `json:"nickname"`
		Vote     int    `json:"vote"`
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logrus.Error(err.DebugMessage)
			return cntx.JSON(err.HTTPCode, Message{Message: err.UserMessage})
		}

		slugOrID := cntx.Param("slug_or_id")

		var threadDetails *models.Thread
		var customErr *errors.Error

		id, err := strconv.ParseUint(slugOrID, 10, 64)
		if err != nil {
			threadDetails, customErr = th.threadUseCase.CreateVoteBySlug(slugOrID, req.Nickname, req.Vote)
			if customErr != nil {
				logrus.Error(customErr.DebugMessage)
				return cntx.JSON(customErr.HTTPCode, Message{customErr.UserMessage})
			}
		} else {
			threadDetails, customErr = th.threadUseCase.CreateVoteByID(id, req.Nickname, req.Vote)
			if customErr != nil {
				logrus.Error(customErr.DebugMessage)
				return cntx.JSON(customErr.HTTPCode, Message{customErr.UserMessage})
			}
		}

		return cntx.JSON(http.StatusOK, threadDetails)
	}
}

func (th *ThreadHandler) ChangeThreadHandler() echo.HandlerFunc {
	type Request struct {
		Title   string `json:"title"`
		Message string `json:"message"`
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logrus.Error(err.DebugMessage)
			return cntx.JSON(err.HTTPCode, Message{Message: err.UserMessage})
		}

		slugOrID := cntx.Param("slug_or_id")

		var threadDetails *models.Thread
		var customErr *errors.Error

		id, err := strconv.ParseUint(slugOrID, 10, 64)
		if err != nil {
			threadDetails, customErr = th.threadUseCase.ChangeBySlug(slugOrID, req.Title, req.Message)
			if customErr != nil {
				logrus.Error(customErr.DebugMessage)
				return cntx.JSON(customErr.HTTPCode, Message{customErr.UserMessage})
			}
		} else {
			threadDetails, customErr = th.threadUseCase.ChangeByID(id, req.Title, req.Message)
			if customErr != nil {
				logrus.Error(customErr.DebugMessage)
				return cntx.JSON(customErr.HTTPCode, Message{customErr.UserMessage})
			}
		}

		return cntx.JSON(http.StatusOK, threadDetails)
	}
}

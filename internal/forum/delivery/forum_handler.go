package delivery

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/technopark_database/internal/consts"
	"github.com/technopark_database/internal/forum"
	"github.com/technopark_database/internal/helpers/errors"
	"github.com/technopark_database/internal/models"
	reader "github.com/technopark_database/tools/requestReader"
	"net/http"
)

type ForumHandler struct {
	forumUseCase forum.ForumUseCase
}

func NewForumHandler(usecase forum.ForumUseCase) *ForumHandler {
	return &ForumHandler{forumUseCase: usecase}
}

func (fh *ForumHandler) Configure(e *echo.Echo) {
	e.POST("/api/forum/create", fh.CreateHandler())
	e.GET("/api/forum/:slug/details", fh.GetInfo())
	e.GET("/api/forum/:slug/threads", fh.GetThreads())
	// TODO: проверить
	//e.GET("/api/forum/:slug/users", fh.GetUsers())
}

type Message struct {
	Message string `json:"message"`
}

func (fh *ForumHandler) CreateHandler() echo.HandlerFunc {
	type Request struct {
		Title string `json:"title"`
		User  string `json:"user"`
		Slug  string `json:"slug"`
	}
	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logrus.Info(err.DebugMessage)
			return cntx.JSON(err.HTTPCode, Message{Message: err.UserMessage})
		}

		forum := &models.Forum{
			Title: req.Title,
			User:  req.User,
			Slug:  req.Slug,
		}

		createdForum, err := fh.forumUseCase.Create(forum)
		if err == errors.Get(consts.CodeForumAlreadyExist) {
			return cntx.JSON(err.HTTPCode, createdForum)
		} else if err == errors.Get(consts.CodeUserDoesNotExist) {
			return cntx.JSON(err.HTTPCode, Message{Message: err.UserMessage})
		} else if err != nil {
			logrus.Error(err.DebugMessage)
			return cntx.JSON(err.HTTPCode, Message{Message: err.UserMessage})
		}

		return cntx.JSON(http.StatusCreated, forum)
	}
}

func (fh *ForumHandler) GetInfo() echo.HandlerFunc {
	return func(cntx echo.Context) error {
		slug := cntx.Param("slug")

		forum, err := fh.forumUseCase.GetDetails(slug)
		if err != nil {
			logrus.Error(err.DebugMessage)
			return cntx.JSON(err.HTTPCode, Message{err.UserMessage})
		}

		return cntx.JSON(http.StatusOK, forum)
	}
}

func (fh *ForumHandler) GetUsers() echo.HandlerFunc {
	type Request struct {
		models.Pagination
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logrus.Info(err.DebugMessage)
			return cntx.JSON(err.HTTPCode, Message{Message: err.UserMessage})
		}

		slug := cntx.Param("slug")

		users, err := fh.forumUseCase.GetUsers(slug, &req.Pagination)
		if err != nil {
			logrus.Error(err.DebugMessage)
			return cntx.JSON(err.HTTPCode, Message{err.UserMessage})
		}

		return cntx.JSON(http.StatusOK, users)
	}
}

func (fh *ForumHandler) GetThreads() echo.HandlerFunc {
	type Request struct {
		models.Pagination
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logrus.Info(err.DebugMessage)
			return cntx.JSON(err.HTTPCode, Message{Message: err.UserMessage})
		}

		slug := cntx.Param("slug")

		threads, err := fh.forumUseCase.GetThreads(slug, &req.Pagination)
		if err != nil {
			logrus.Error(err.DebugMessage)
			return cntx.JSON(err.HTTPCode, Message{err.UserMessage})
		}

		return cntx.JSON(http.StatusOK, threads)
	}
}

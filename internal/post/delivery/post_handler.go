package delivery

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	//"github.com/sirupsen/logrus"
	"github.com/technopark_database/internal/consts"
	"github.com/technopark_database/internal/helpers/errors"
	"github.com/technopark_database/internal/models"
	"io/ioutil"
	"strings"

	"github.com/technopark_database/internal/post"
	reader "github.com/technopark_database/tools/requestReader"
	"net/http"
	"strconv"
)

type PostHandler struct {
	postUseCase post.PostUseCase
}

func NewPostHandler(postUseCase post.PostUseCase) *PostHandler {
	return &PostHandler{postUseCase: postUseCase}
}

func (ph *PostHandler) Configure(e *echo.Echo) {
	e.POST("/api/thread/:slug_or_id/create", ph.CreatePostsHandler())
	e.GET("/api/thread/:slug_or_id/posts", ph.GetPosts())
	e.POST("/api/post/:id/details", ph.ChangeHandler())
	e.GET("/api/post/:id/details", ph.GetPostDetails())
}

type Message struct {
	Message string `json:"message"`
}

func (ph *PostHandler) CreatePostsHandler() echo.HandlerFunc {
	type Request struct {
		Posts []*models.Post `json:"posts"`
	}
	return func(ctx echo.Context) error {
		body, err := ioutil.ReadAll(ctx.Request().Body)
		if err != nil {
			customErr := errors.New(consts.CodeInternalServerError, err)
			//logrus.Error(customErr.DebugMessage)
			return ctx.JSON(customErr.HTTPCode, Message{Message: customErr.UserMessage})
		}

		req := []*models.Post{}
		if err := json.Unmarshal(body, &req); err != nil {
			customErr := errors.New(consts.CodeInternalServerError, err)
			//logrus.Error(customErr.DebugMessage)
			return ctx.JSON(customErr.HTTPCode, Message{Message: customErr.UserMessage})
		}

		slugOrID := ctx.Param("slug_or_id")

		createdPosts, customErr := ph.postUseCase.CreateMany(slugOrID, req)
		if customErr != nil {
			//logrus.Error(customErr.DebugMessage)
			return ctx.JSON(customErr.HTTPCode, Message{Message: customErr.UserMessage})
		}
		return ctx.JSON(http.StatusCreated, createdPosts)
	}
}

func (ph *PostHandler) ChangeHandler() echo.HandlerFunc {
	type Request struct {
		Message string `json:"message"`
	}
	return func(ctx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(ctx).Read(req); err != nil {
			//logrus.Error(err.DebugMessage)
			return ctx.JSON(err.HTTPCode, Message{Message: err.UserMessage})
		}

		strID := ctx.Param("id")
		id, _ := strconv.ParseUint(strID, 10, 64)

		post, err := ph.postUseCase.ChangeByID(id, req.Message)
		if err != nil {
			//logrus.Error(err.DebugMessage)
			return ctx.JSON(err.HTTPCode, Message{Message: err.UserMessage})
		}

		return ctx.JSON(http.StatusOK, post)
	}
}

func (ph *PostHandler) GetPosts() echo.HandlerFunc {
	type Request struct {
		Sort  string `query:"sort"`
		Since uint64 `query:"since"`
		models.Pagination
	}
	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			//logrus.Error(err.DebugMessage)
			return cntx.JSON(err.HTTPCode, Message{Message: err.UserMessage})
		}

		slugOrID := cntx.Param("slug_or_id")

		posts, customErr := ph.postUseCase.GetPosts(slugOrID, req.Sort, req.Since, &req.Pagination)
		if customErr != nil {
			//logrus.Error(customErr.DebugMessage)
			return cntx.JSON(customErr.HTTPCode, Message{Message: customErr.UserMessage})
		}
		return cntx.JSON(http.StatusOK, posts)
	}
}

func (ph *PostHandler) GetPostDetails() echo.HandlerFunc {
	type Request struct {
		models.Related
	}
	return func(cntx echo.Context) error {
		related := cntx.QueryParam("related")

		relatedModel := &models.Related{}

		if strings.Contains(related, "user") {
			relatedModel.User = true
		}
		if strings.Contains(related, "thread") {
			relatedModel.Thread = true
		}
		if strings.Contains(related, "forum") {
			relatedModel.Forum = true
		}

		strID := cntx.Param("id")
		id, _ := strconv.ParseUint(strID, 10, 64)

		posts, customErr := ph.postUseCase.GetPostInfo(id, relatedModel)
		if customErr != nil {
			//logrus.Error(customErr.DebugMessage)
			return cntx.JSON(customErr.HTTPCode, Message{Message: customErr.UserMessage})
		}
		return cntx.JSON(http.StatusOK, posts)
	}
}

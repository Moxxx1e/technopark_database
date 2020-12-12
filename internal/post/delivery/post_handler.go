package delivery

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

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
	// TODO: множественный инсёрт, потом
	e.POST("/api/v1/thread/:slug_or_id/create", ph.CreatePostsHandler())
	//
	e.POST("/api/v1/post/:id/details", ph.ChangeHandler())
}

type Message struct {
	Message string `json:"message"`
}

func (ph *PostHandler) CreatePostsHandler() echo.HandlerFunc {
	//type Request struct {
	//	Posts []*models.Post `json:"posts"`
	//}
	return func(ctx echo.Context) error {
		//	req := &Request{}
		//	if err := reader.NewRequestReader(ctx).Read(req); err != nil {
		//		logrus.Error(err.DebugMessage)
		//		return ctx.JSON(err.HTTPCode, Message{Message: err.UserMessage})
		//	}
		//
		//	err := ph.postUseCase
		//
		//	return ctx.JSON{http.StatusOK}
		return nil
	}
}

func (ph *PostHandler) ChangeHandler() echo.HandlerFunc {
	type Request struct {
		Message string `json:"message"`
	}
	return func(ctx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(ctx).Read(req); err != nil {
			logrus.Error(err.DebugMessage)
			return ctx.JSON(err.HTTPCode, Message{Message: err.UserMessage})
		}

		strID := ctx.Param("id")
		id, _ := strconv.ParseUint(strID, 10, 64)

		post, err := ph.postUseCase.ChangeByID(id, req.Message)
		if err != nil {
			logrus.Error(err.DebugMessage)
			return ctx.JSON(err.HTTPCode, Message{Message: err.UserMessage})
		}

		return ctx.JSON(http.StatusOK, post)
	}
}

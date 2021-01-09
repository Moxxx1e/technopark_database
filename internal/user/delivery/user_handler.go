package delivery

import (
	"github.com/labstack/echo/v4"
	//"github.com/sirupsen/logrus"
	"github.com/technopark_database/internal/consts"
	"github.com/technopark_database/internal/helpers/errors"
	"github.com/technopark_database/internal/models"
	"github.com/technopark_database/internal/user"
	reader "github.com/technopark_database/tools/requestReader"
	"net/http"
)

type UserHandler struct {
	userUseCase user.UserUseCase
}

func NewUserHandler(userUseCase user.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: userUseCase}
}

func (uh *UserHandler) Configure(e *echo.Echo) {
	e.POST("/api/user/:nickname/create", uh.CreateProfileHandler())
	e.GET("/api/user/:nickname/profile", uh.GetProfileHandler())
	e.POST("/api/user/:nickname/profile", uh.ChangeProfileHandler())
}

type Message struct {
	Message string `json:"message"`
}

func (uh *UserHandler) CreateProfileHandler() echo.HandlerFunc {
	type CreateRequest struct {
		Fullname string `json:"fullname"`
		About    string `json:"about"`
		Email    string `json:"email"`
	}

	return func(context echo.Context) error {
		nickname := context.Param("nickname")

		req := &CreateRequest{}
		if err := reader.NewRequestReader(context).Read(req); err != nil {
			//logrus.Info(err.DebugMessage)
			return context.JSON(err.HTTPCode, Message{Message: err.UserMessage})
		}

		newUser := &models.User{
			Nickname: nickname,
			Fullname: req.Fullname,
			About:    req.About,
			Email:    req.Email,
		}

		users, err := uh.userUseCase.Create(newUser)
		if err == errors.Get(consts.CodeUserEmailConflicts) ||
			err == errors.Get(consts.CodeUserNicknameConflicts) {
			return context.JSON(err.HTTPCode, users)
		} else if err != nil {
			return context.JSON(err.HTTPCode, Message{Message: err.UserMessage})
		}

		return context.JSON(http.StatusCreated, newUser)
	}
}

func (uh *UserHandler) GetProfileHandler() echo.HandlerFunc {
	return func(context echo.Context) error {
		nickname := context.Param("nickname")

		dbUser, err := uh.userUseCase.GetUserInfo(nickname)
		if err != nil {
			//logrus.Info(err.DebugMessage)
			return context.JSON(err.HTTPCode, Message{Message: err.UserMessage})
		}

		return context.JSON(http.StatusOK, dbUser)
	}
}

func (uh *UserHandler) ChangeProfileHandler() echo.HandlerFunc {
	type ChangeRequest struct {
		Fullname string `json:"fullname"`
		About    string `json:"about"`
		Email    string `json:"email"`
	}

	return func(context echo.Context) error {
		nickname := context.Param("nickname")

		req := &ChangeRequest{}
		if err := reader.NewRequestReader(context).Read(req); err != nil {
			//logrus.Info(err.DebugMessage)
			return context.JSON(err.HTTPCode, Message{Message: err.UserMessage})
		}

		updateUser := &models.User{
			Nickname: nickname,
			Fullname: req.Fullname,
			About:    req.About,
			Email:    req.Email,
		}

		user, customErr := uh.userUseCase.Change(updateUser)
		if customErr != nil {
			//logrus.Info(customErr.DebugMessage)
			return context.JSON(customErr.HTTPCode, Message{Message: customErr.UserMessage})
		}

		return context.JSON(http.StatusOK, user)
	}
}

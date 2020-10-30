package delivery

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/technopark_database/internal/models"
	"github.com/technopark_database/internal/user/usecases"
	reader "github.com/technopark_database/tools/requestReader"
	"net/http"
)

type UserHandler struct {
	userUseCase *usecases.UserUseCase
}

func NewUserHandler(userUseCase *usecases.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: userUseCase}
}

func (uh *UserHandler) Configure(e *echo.Echo) {
	e.POST("api/v1/user/:nickname/create", uh.CreateProfileHandler())
	e.GET("api/v1/user/:nickname/profile", uh.GetProfileHandler())
	e.POST("api/v1/user/:nickname/profile", uh.ChangeProfileHandler())
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
			logrus.Info(err.DebugMessage)
			return context.JSON(err.HTTPCode, Message{Message: err.UserMessage})
		}

		profile := &models.User{
			Nickname: nickname,
			Fullname: req.Fullname,
			About:    req.About,
			Email:    req.Email,
		}

		if err := uh.userUseCase.Create(profile); err != nil {
			logrus.Info(err.DebugMessage)
			return context.JSON(err.HTTPCode, Message{Message: err.UserMessage})
		}

		return context.JSON(http.StatusCreated, profile)
	}
}

func (uh *UserHandler) GetProfileHandler() echo.HandlerFunc {
	type GetProfileResponse struct {
		Message string `json:"message"`
	}

	return func(context echo.Context) error {
		nickname := context.Param("nickname")

		dbProfile, err := uh.userUseCase.GetUserInfo(nickname)
		if err != nil {
			logrus.Info(err.DebugMessage)
			return context.JSON(err.HTTPCode, Message{Message: err.UserMessage})
		}

		return context.JSON(http.StatusCreated, dbProfile)
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
			logrus.Info(err.DebugMessage)
			return context.JSON(err.HTTPCode, Message{Message: err.UserMessage})
		}

		profile := &models.User{
			Nickname: nickname,
			Fullname: req.Fullname,
			About:    req.About,
			Email:    req.Email,
		}

		if err := uh.userUseCase.UpdateUserInfo(nickname, profile); err != nil {
			logrus.Info(err.DebugMessage)
			return context.JSON(err.HTTPCode, Message{Message: err.UserMessage})
		}

		return context.JSON(http.StatusOK, profile)
	}
}

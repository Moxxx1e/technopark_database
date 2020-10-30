package delivery

import (
	"github.com/labstack/echo/v4"
	"github.com/technopark_database/internal/models"
	"github.com/technopark_database/internal/user/usecases"
	reader "github.com/technopark_database/tools/requestReader"
	"log"
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
			log.Println(err)
			return context.JSON(http.StatusBadRequest, err)
		}

		profile := &models.User{
			Nickname: nickname,
			Fullname: req.Fullname,
			About:    req.About,
			Email:    req.Email,
		}

		dbProfile, err := uh.userUseCase.GetUserInfo(nickname)
		if err != nil {
			log.Println(err)
			return context.JSON(http.StatusInternalServerError, err)
		}
		if dbProfile != nil {
			return context.JSON(http.StatusConflict, dbProfile)
		}

		if err := uh.userUseCase.Create(profile); err != nil {
			log.Println(err)
			return context.JSON(http.StatusInternalServerError, err)
		}

		return context.JSON(http.StatusCreated, profile)
	}
}

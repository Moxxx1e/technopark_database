package delivery

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/technopark_database/internal/service"
	"net/http"
)

type ServiceHandler struct {
	serviceUseCase service.ServiceUseCase
}

func NewServiceHandler(serviceUseCase service.ServiceUseCase) *ServiceHandler {
	return &ServiceHandler{serviceUseCase: serviceUseCase}
}

func (sh *ServiceHandler) Configure(e *echo.Echo) {
	e.POST("/api/service/clear", sh.ClearHandler())
	e.GET("/api/service/status", sh.GetStatusHandler())
}

func (sh *ServiceHandler) ClearHandler() echo.HandlerFunc {
	return func(cntx echo.Context) error {
		err := sh.serviceUseCase.Delete()
		if err != nil {
			logrus.Error(err.DebugMessage)
			return cntx.JSON(err.HTTPCode, err.UserMessage)
		}
		return cntx.JSON(http.StatusOK, "")
	}
}

func (sh *ServiceHandler) GetStatusHandler() echo.HandlerFunc {
	return func(cntx echo.Context) error {
		serviceStatus, err := sh.serviceUseCase.GetStatus()
		if err != nil {
			logrus.Error(err.DebugMessage)
			return cntx.JSON(err.HTTPCode, err.UserMessage)
		}
		return cntx.JSON(http.StatusOK, serviceStatus)
	}
}

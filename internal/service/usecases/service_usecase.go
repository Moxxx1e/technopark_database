package usecases

import (
	"github.com/technopark_database/internal/consts"
	"github.com/technopark_database/internal/helpers/errors"
	"github.com/technopark_database/internal/models"
	"github.com/technopark_database/internal/service"
)

type ServiceUseCase struct {
	rep service.ServiceRepository
}

func (su *ServiceUseCase) GetStatus() (*models.ServiceStatus, *errors.Error) {
	serviceStatus, err := su.rep.GetStatus()
	if err != nil {
		return nil, errors.New(consts.CodeInternalServerError, err)
	}
	return serviceStatus, nil
}

func NewServiceUseCase(rep service.ServiceRepository) service.ServiceUseCase {
	return &ServiceUseCase{rep: rep}
}

func (su *ServiceUseCase) Delete() *errors.Error {
	err := su.rep.Delete()
	if err != nil {
		return errors.Get(consts.CodeCantDeleteDatabase)
	}
	return nil
}

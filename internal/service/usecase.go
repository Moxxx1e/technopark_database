package service

import (
	"github.com/technopark_database/internal/helpers/errors"
	"github.com/technopark_database/internal/models"
)

type ServiceUseCase interface {
	Delete() *errors.Error
	GetStatus() (*models.ServiceStatus, *errors.Error)
}

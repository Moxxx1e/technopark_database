package service

import "github.com/technopark_database/internal/models"

type ServiceRepository interface {
	Delete() error
	GetStatus() (*models.ServiceStatus, error)
}

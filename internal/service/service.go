package service

import (
	"github.com/EMus88/Market/internal/models"
	"github.com/EMus88/Market/internal/repository"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

type Repository interface {
	//auth methods
	SaveUser(user *models.User) (string, error)
	GetUser(user *models.User) (string, string, error)
	CheckUser(id uuid.UUID) (string, error)
	AddCategory(m *models.Category) error
	AddProduct(m *models.ProductDTO) error
	ChangeVisible(v *models.Visible) error
	GetCatalog() ([]models.ProductDTO, error)
	GetByCategory(productName string, category string) ([]models.ProductDTO, error)
	GetByAllCategories(productName string) ([]models.ProductDTO, error)
}

type Service struct {
	Repository
	Auth
	logger *logrus.Logger
}

func NewService(r *repository.Repository, logger *logrus.Logger) *Service {
	return &Service{
		Repository: r,
		Auth:       *NewAuth(r),
		logger:     logger,
	}
}

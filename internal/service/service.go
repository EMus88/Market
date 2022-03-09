package service

import (
	"JWT_auth/internal/model"
	"JWT_auth/internal/repository"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

type Repository interface {
	//auth methods
	SaveUser(user *model.User) (string, error)
	GetUser(user *model.User) (string, string, error)
	CheckUser(id uuid.UUID) (string, error)
	AddCategory(m *model.Category) error
	AddProduct(m *model.ProductDTO) error
	ChangeVisible(v *model.Visible) error
	GetCatalog() ([]model.ProductDTO, error)
	GetByCategory(productName string, category string) ([]model.ProductDTO, error)
	GetByAllCategories(productName string) ([]model.ProductDTO, error)
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

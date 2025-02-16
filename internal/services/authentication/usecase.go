package authentication

import (
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UseCase struct {
	DB         *gorm.DB
	Log        *logrus.Logger
	Validate   *validator.Validate
	Repository *Repository
}

func NewUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate, repository *Repository) *UseCase {
	return &UseCase{
		DB:         db,
		Log:        logger,
		Validate:   validate,
		Repository: repository,
	}
}

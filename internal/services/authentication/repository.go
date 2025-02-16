package authentication

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Repository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewRepository(log *logrus.Logger) *Repository {
	return &Repository{
		Log: log,
	}
}

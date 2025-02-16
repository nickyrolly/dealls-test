package testdb

import (
	"os"

	"github.com/glebarez/sqlite"
	"github.com/nickyrolly/dealls-test/internal/services/user"
	"github.com/nickyrolly/dealls-test/internal/testutil/testmodels"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// NewTestDB creates a new test database connection
func NewTestDB() (*gorm.DB, error) {
	// Use in-memory SQLite for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate test tables
	err = db.AutoMigrate(
		&user.Entity{},
		&testmodels.UserProfile{},
		&testmodels.UserPhoto{},
		&testmodels.UserPreference{},
		&testmodels.UserMatch{},
		&testmodels.UserLike{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// NewTestLogger creates a new test logger
func NewTestLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)
	return logger
}

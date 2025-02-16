package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/nickyrolly/dealls-test/internal/delivery/http/authentication"
	authService "github.com/nickyrolly/dealls-test/internal/services/authentication"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupAuthTest(t *testing.T) (*gorm.DB, *logrus.Logger, func()) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)

	cleanup := func() {
		sqlDB, err := db.DB()
		if err != nil {
			t.Errorf("Failed to get database instance: %v", err)
			return
		}
		err = sqlDB.Close()
		if err != nil {
			t.Errorf("Failed to close database: %v", err)
		}
	}

	return db, log, cleanup
}

func TestAuthEndpoints(t *testing.T) {
	// Setup
	db, log, cleanup := setupAuthTest(t)
	defer cleanup()

	// Initialize service and controller
	authSvc := authService.NewService(db, log)
	authCtrl := authentication.NewController(authSvc, log)

	// Test SignUp endpoint
	t.Run("SignUp", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/signup", nil)
		w := httptest.NewRecorder()

		authCtrl.SignUp(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code) // Should fail due to invalid request body
	})

	// Test Login endpoint
	t.Run("Login", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
		w := httptest.NewRecorder()

		authCtrl.Login(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code) // Should fail due to invalid request body
	})
}

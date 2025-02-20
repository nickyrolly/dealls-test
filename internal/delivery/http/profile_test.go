package http

import (
	"testing"

	"github.com/glebarez/sqlite"
	profileService "github.com/nickyrolly/dealls-test/internal/services/profile"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func setupProfileTest(t *testing.T) (*gorm.DB, *logrus.Logger, func()) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)

	// Auto-migrate the database
	err = db.AutoMigrate(&profileService.UserProfile{}, &profileService.UserPhoto{}, &profileService.UserPreference{}, &profileService.UserMatch{}, &profileService.UserLike{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

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

// func TestProfileEndpoints(t *testing.T) {
// 	// Setup
// 	db, log, cleanup := setupProfileTest(t)
// 	defer cleanup()

// 	// Initialize service and controller
// 	profileSvc := profileService.NewService(db, log)
// 	profileCtrl := profile.NewController(log, profileSvc)

// 	// Create test user ID
// 	testUserID := "123e4567-e89b-12d3-a456-426614174000"

// 	// Test GetProfile endpoint
// 	t.Run("GetProfile", func(t *testing.T) {
// 		req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
// 		// Set user ID in context
// 		ctx := context.WithValue(req.Context(), "user_id", testUserID)
// 		req = req.WithContext(ctx)
// 		w := httptest.NewRecorder()

// 		profileCtrl.HandleGetProfile(w, req)
// 		assert.Equal(t, http.StatusNotFound, w.Code)
// 	})

// 	// Test UpdateProfile endpoint
// 	t.Run("UpdateProfile", func(t *testing.T) {
// 		reqBody := `{"height": 180, "weight": 75, "occupation": "Software Engineer"}`
// 		req := httptest.NewRequest(http.MethodPut, "/api/v1/profile", strings.NewReader(reqBody))
// 		// Set user ID in context
// 		ctx := context.WithValue(req.Context(), "user_id", testUserID)
// 		req = req.WithContext(ctx)
// 		w := httptest.NewRecorder()

// 		profileCtrl.HandleUpdateProfile(w, req)
// 		assert.Equal(t, http.StatusOK, w.Code)
// 	})
// }

package profile

import (
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/nickyrolly/dealls-test/internal/services/user"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Run migrations
	err = db.AutoMigrate(&user.Entity{}, &UserProfile{}, &UserPhoto{}, &UserPreference{}, &UserLike{}, &UserMatch{})
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
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

	return db, cleanup
}

func createTestUser(db *gorm.DB) (*user.Entity, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &user.Entity{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
		FirstName:    "Test",
		LastName:     nil,
		DateOfBirth:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		Gender:       "male",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// func TestGetUserProfile(t *testing.T) {
// 	// Setup
// 	db, cleanup := setupTestDB(t)
// 	defer cleanup()

// 	log := logrus.New()
// 	service := NewService(db, log)

// 	// Create test user
// 	testUser, err := createTestUser(db)
// 	require.NoError(t, err)

// 	// Test getting non-existent profile
// 	profile, err := service.GetUserProfile(testUser.ID)
// 	require.NoError(t, err)
// 	assert.Nil(t, profile)

// 	// Create profile
// 	height := 180
// 	weight := 75
// 	education := "Bachelor's"
// 	aboutMe := "Test profile"
// 	testProfile := &UserProfile{
// 		UserID:    testUser.ID,
// 		Height:    &height,
// 		Weight:    &weight,
// 		Education: &education,
// 		AboutMe:   &aboutMe,
// 	}
// 	err = db.Create(testProfile).Error
// 	require.NoError(t, err)

// 	// Test getting existing profile
// 	profile, err = service.GetUserProfile(testUser.ID)
// 	require.NoError(t, err)
// 	assert.NotNil(t, profile)
// 	assert.Equal(t, testUser.ID, profile.UserID)
// 	assert.Equal(t, *testProfile.Height, *profile.Height)
// 	assert.Equal(t, *testProfile.AboutMe, *profile.AboutMe)
// }

// func TestUpdateUserProfile(t *testing.T) {
// 	// Setup
// 	db, cleanup := setupTestDB(t)
// 	defer cleanup()

// 	log := logrus.New()
// 	service := NewService(db, log)

// 	// Create test user
// 	testUser, err := createTestUser(db)
// 	require.NoError(t, err)

// 	// Create initial profile
// 	height := 180
// 	weight := 75
// 	education := "Bachelor's"
// 	aboutMe := "Test profile"
// 	initialProfile := &UserProfile{
// 		UserID:    testUser.ID,
// 		Height:    &height,
// 		Weight:    &weight,
// 		Education: &education,
// 		AboutMe:   &aboutMe,
// 	}
// 	err = db.Create(initialProfile).Error
// 	require.NoError(t, err)

// 	// Update profile
// 	newHeight := 185
// 	newAboutMe := "Updated profile"
// 	updates := map[string]interface{}{
// 		"height":   &newHeight,
// 		"about_me": &newAboutMe,
// 	}

// 	err = service.UpdateUserProfile(testUser.ID, updates)
// 	require.NoError(t, err)

// 	// Verify update
// 	profile, err := service.GetUserProfile(testUser.ID)
// 	require.NoError(t, err)
// 	assert.NotNil(t, profile)
// 	assert.Equal(t, newHeight, *profile.Height)
// 	assert.Equal(t, newAboutMe, *profile.AboutMe)
// }

func TestCreateLike(t *testing.T) {
	// Setup
	db, cleanup := setupTestDB(t)
	defer cleanup()

	log := logrus.New()
	service := NewService(db, log)

	// Create test users
	user1, err := createTestUser(db)
	require.NoError(t, err)

	// Create another test user with a different email
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	require.NoError(t, err)

	user2 := &user.Entity{
		ID:           uuid.New(),
		Email:        "test2@example.com",
		PasswordHash: string(hashedPassword),
		FirstName:    "Test2",
		LastName:     nil,
		DateOfBirth:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		Gender:       "female",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err = db.Create(user2).Error
	require.NoError(t, err)

	// Test creating like
	like := &UserLike{
		LikerID: user1.ID,
		LikedID: user2.ID,
	}

	match, err := service.CreateLike(like)
	require.NoError(t, err)
	assert.Nil(t, match) // No match yet as user2 hasn't liked user1

	// Test mutual like (creates match)
	like2 := &UserLike{
		LikerID: user2.ID,
		LikedID: user1.ID,
	}

	match, err = service.CreateLike(like2)
	require.NoError(t, err)
	assert.NotNil(t, match)
	assert.True(t, match.User1ID == user1.ID && match.User2ID == user2.ID ||
		match.User1ID == user2.ID && match.User2ID == user1.ID)
}

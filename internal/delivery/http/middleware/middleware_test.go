package middleware

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
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

func setupTestMiddleware(t *testing.T) (*Middleware, *redis.Pool, func()) {
	// Setup test database
	db, cleanup := setupTestDB(t)

	// Setup Redis pool
	redisPool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}

	// Create middleware
	mw := NewMiddleware(logrus.New(), db, redisPool, "test-secret")

	return mw, redisPool, cleanup
}

func TestBasicAuth(t *testing.T) {
	mw, _, cleanup := setupTestMiddleware(t)
	defer cleanup()

	// Create test handler
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create test server
	handler := mw.BasicAuth(nextHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	// Create test cases
	tests := []struct {
		name           string
		auth           string
		expectedStatus int
	}{
		{
			name:           "No auth header",
			auth:           "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid auth format",
			auth:           "Invalid",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid base64",
			auth:           "Basic invalid-base64",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid credentials format",
			auth:           fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte("invalid"))),
			expectedStatus: http.StatusUnauthorized,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", server.URL, nil)
			if tt.auth != "" {
				req.Header.Set("Authorization", tt.auth)
			}

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestJWT(t *testing.T) {
	mw, _, cleanup := setupTestMiddleware(t)
	defer cleanup()

	// Create test handler
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create test server
	handler := mw.JWT(nextHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	// Create test cases
	tests := []struct {
		name           string
		auth           string
		expectedStatus int
	}{
		{
			name:           "No auth header",
			auth:           "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid auth format",
			auth:           "Invalid",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid token",
			auth:           "Bearer invalid-token",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", server.URL, nil)
			if tt.auth != "" {
				req.Header.Set("Authorization", tt.auth)
			}

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestCORS(t *testing.T) {
	mw, _, cleanup := setupTestMiddleware(t)
	defer cleanup()

	// Create test handler
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create test server
	handler := mw.CORS(nextHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	// Create test cases
	tests := []struct {
		name           string
		method         string
		expectedStatus int
	}{
		{
			name:           "OPTIONS request",
			method:         "OPTIONS",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GET request",
			method:         "GET",
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, server.URL, nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
		})
	}
}

package middleware

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/nickyrolly/dealls-test/internal/common"
	"github.com/nickyrolly/dealls-test/internal/services/user"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"errors"
)

// Middleware handles all HTTP middleware functionality
type Middleware struct {
	log       *logrus.Logger
	db        *gorm.DB
	redisPool *redis.Pool
	jwtSecret string
}

// NewMiddleware creates a new Middleware instance
func NewMiddleware(log *logrus.Logger, db *gorm.DB, redisPool *redis.Pool, jwtSecret string) *Middleware {
	return &Middleware{
		log:       log,
		db:        db,
		redisPool: redisPool,
		jwtSecret: jwtSecret,
	}
}

type contextKey string

const (
	userContextKey      contextKey = "user"
	userIDContextKey    contextKey = "user_id"
	claimsContextKey    contextKey = "claims"
	requestIDContextKey contextKey = "request_id"
)

// Session middleware checks if the user has a valid session
func (m *Middleware) Session(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get session token from cookie
		cookie, err := r.Cookie("session_token")
		if err != nil {
			m.log.Warnf("No session token found: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Get session from Redis
		conn := m.redisPool.Get()
		defer conn.Close()

		userID, err := redis.String(conn.Do("GET", cookie.Value))
		if err != nil {
			m.log.Warnf("Invalid session token: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Set user ID in context
		context.Set(r, userIDContextKey, userID)
		defer context.Clear(r)

		next.ServeHTTP(w, r)
	})
}

// BasicAuth middleware handles basic authentication
func (m *Middleware) BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get Authorization header
		auth := r.Header.Get("Authorization")
		if auth == "" {
			m.log.Warn("No Authorization header")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if it's Basic auth
		const prefix = "Basic "
		if !strings.HasPrefix(auth, prefix) {
			m.log.Warn("Not Basic auth")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Decode base64 credentials
		payload, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
		if err != nil {
			m.log.Warnf("Invalid base64: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Split into username and password
		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 {
			m.log.Warn("Invalid auth format")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		email := pair[0]
		password := pair[1]

		// Find user in database
		var user user.Entity
		if err := m.db.Where("email = ?", email).First(&user).Error; err != nil {
			m.log.Warnf("User not found: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Verify password
		if !user.VerifyPassword(password) {
			m.log.Warn("Invalid password")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Set user in context
		context.Set(r, userContextKey, user)
		defer context.Clear(r)

		next.ServeHTTP(w, r)
	})
}

// JWT middleware handles JWT authentication
func (m *Middleware) JWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get Authorization header

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.log.Warn("No Authorization header")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if it's Bearer auth
		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			m.log.Warn("Not Bearer auth")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Get token
		tokenString := authHeader[len(prefix):]

		fmt.Println("====== token string ======", tokenString)
		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			fmt.Println("====== token ======", token)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(m.jwtSecret), nil
		})

		if err != nil {
			m.log.Warnf("Invalid token: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Get claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			m.log.Warn("Invalid claims")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		fmt.Println("--- JWT Success")
		// Set claims in context

		context.Set(r, "user", claims["user_id"])
		defer context.Clear(r)

		next.ServeHTTP(w, r)
	})
}

// AuthenticateCredentials middleware checks credentials and manages authentication
func (m *Middleware) AuthenticateCredentials(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("[AuthenticateCredentials]")
		var credentials struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		// Decode credentials from request body
		if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
			m.log.WithError(err).Error("Failed to decode credentials")
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Check Redis cache first
		redisKey := fmt.Sprintf("auth:%s", credentials.Email)

		conn := m.redisPool.Get()
		defer conn.Close()

		userJSON, redisErr := redis.String(conn.Do("GET", redisKey))
		if redisErr == nil {
			fmt.Println("[AuthenticateCredentials] User found in Redis")
			// Found in Redis, unmarshal and verify password
			var user user.Entity
			if unmarshalErr := json.Unmarshal([]byte(userJSON), &user); unmarshalErr == nil {
				// Verify password
				if user.VerifyPassword(credentials.Password) {
					// Set user in context
					context.Set(r, "user", user.ID.String())
					defer context.Clear(r)

					next.ServeHTTP(w, r)
					return
				}
			}
		}

		// Check database if not found in Redis
		fmt.Println("[AuthenticateCredentials] User Not found in Redis, Check DB")
		var user user.Entity
		if err := m.db.Where("email = ?", credentials.Email).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, "Invalid credentials", http.StatusUnauthorized)
				return
			}
			m.log.WithError(err).Error("Database error")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Verify password
		if !user.VerifyPassword(credentials.Password) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		fmt.Println("[AuthenticateCredentials] User found in DB")

		// Cache user in Redis
		userJSON1, _ := json.Marshal(user)
		// m.redisPool.Set(redisKey, userJSON, 24*time.Hour)

		// _, redisErr = conn.Do("SETEX", redisKey, 24*time.Hour, userJSON1) // 900 seconds = 15 minutes
		_, redisErr = conn.Do("SETEX", redisKey, 900, userJSON1) // 900 seconds = 15 minutes
		if redisErr != nil {
			// Log Redis error but don't fail the request
			fmt.Printf("Redis cache error: %v\n", redisErr)
		}

		fmt.Println("[AuthenticateCredentials] User cached in Redis")

		// Set user in context
		context.Set(r, "user", user.ID.String())
		defer context.Clear(r)

		next.ServeHTTP(w, r)
	})
}

// CORS middleware handles Cross-Origin Resource Sharing
func (m *Middleware) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RequireJSON middleware ensures the request has JSON content type
func (m *Middleware) RequireJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check Content-Type header
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			m.log.Warnf("Invalid Content-Type: %s", contentType)
			http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Recover middleware handles panic recovery
func (m *Middleware) Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				m.log.Errorf("Panic recovered: %v", err)
				response := common.ErrorResponse{
					Error: "Internal server error",
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// Logger middleware logs HTTP requests
func (m *Middleware) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response wrapper to capture the status code
		rw := &responseWriter{w, http.StatusOK}

		// Set request ID in context
		requestID := uuid.New().String()
		context.Set(r, requestIDContextKey, requestID)
		defer context.Clear(r)

		next.ServeHTTP(rw, r)

		// Log request details
		m.log.WithFields(logrus.Fields{
			"method":     r.Method,
			"path":       r.URL.Path,
			"status":     rw.status,
			"duration":   time.Since(start),
			"request_id": requestID,
		}).Info("HTTP request")
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

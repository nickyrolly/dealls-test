# Dating App Architecture Documentation

## Overview
The Dating App is built using a clean architecture pattern, separating concerns into distinct layers for better maintainability and testability.

## Project Structure

```
dealls-test/
├── cmd/
│   ├── api/         # API entry point
│   ├── app/         # Application entry point
│   └── migrate/     # Database migrations
├── internal/
│   ├── common/      # Shared utilities
│   ├── config/      # Configuration management
│   ├── delivery/    # HTTP handlers and middleware
│   │   └── http/
│   ├── mock/        # Test mocks
│   └── services/    # Business logic
│       ├── authentication/
│       ├── profile/
│       └── user/
└── docs/            # Documentation
```

## Architecture Layers

### 1. Delivery Layer (HTTP)
- Handles HTTP requests and responses
- Implements middleware for authentication and request validation
- Routes requests to appropriate service methods
- Located in `internal/delivery/http`

### 2. Service Layer
- Contains business logic
- Implements domain rules and validations
- Manages transactions and data consistency
- Located in `internal/services`

### 3. Data Layer
- Uses GORM as ORM for database operations
- Implements repository pattern for data access
- Supports SQLite for testing and PostgreSQL for production
- Models are defined in respective service packages

## Key Components

### Authentication
- JWT-based authentication
- Secure password hashing using bcrypt
- Session management with Redis
- Middleware for protecting routes

### Profile Management
- CRUD operations for user profiles
- Photo management with primary photo support
- Profile search and filtering

### Matching System
- Like/Pass functionality
- Match creation on mutual likes
- Match status management

## Database Schema

### User Profiles
```sql
CREATE TABLE user_profiles (
    id UUID PRIMARY KEY,
    user_id UUID UNIQUE,
    name VARCHAR(255),
    bio TEXT,
    age INTEGER,
    gender VARCHAR(50),
    location VARCHAR(255),
    height INTEGER,
    weight INTEGER,
    occupation VARCHAR(255),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

### User Photos
```sql
CREATE TABLE user_photos (
    id UUID PRIMARY KEY,
    user_id UUID,
    url VARCHAR(255),
    is_primary BOOLEAN,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES user_profiles(user_id)
);
```

### User Matches
```sql
CREATE TABLE user_matches (
    id UUID PRIMARY KEY,
    user1_id UUID,
    user2_id UUID,
    matched_at TIMESTAMP,
    status VARCHAR(50),
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    FOREIGN KEY (user1_id) REFERENCES user_profiles(user_id),
    FOREIGN KEY (user2_id) REFERENCES user_profiles(user_id)
);
```

### User Likes
```sql
CREATE TABLE user_likes (
    id UUID PRIMARY KEY,
    liker_id UUID,
    liked_id UUID,
    created_at TIMESTAMP,
    FOREIGN KEY (liker_id) REFERENCES user_profiles(user_id),
    FOREIGN KEY (liked_id) REFERENCES user_profiles(user_id)
);
```

## Testing Strategy

### Unit Tests
- Service layer tests with mocked dependencies
- Controller tests with mocked services
- Located alongside the code they test

### Integration Tests
- End-to-end API tests
- Uses in-memory SQLite database
- Tests authentication flow and data persistence

## Dependencies

### Core Dependencies
- `github.com/google/uuid`: UUID generation
- `github.com/sirupsen/logrus`: Logging
- `gorm.io/gorm`: ORM
- `github.com/gomodule/redigo`: Redis client
- `github.com/golang-jwt/jwt`: JWT authentication

### Test Dependencies
- `github.com/stretchr/testify`: Testing framework
- `github.com/glebarez/sqlite`: SQLite driver for testing

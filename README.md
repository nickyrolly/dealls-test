# Dating App API

A modern dating application API built with Go, featuring user authentication, profile management, and matching functionality.

## Prerequisites

Before running the application, make sure you have the following installed:

- Go 1.21 or higher
- SQLLite3
- Redis 7 or higher

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/nickyrolly/dealls-test.git
cd dealls-test
```

### 2. Environment Setup

Create a `.env` file in the root directory:

```bash

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT Configuration
JWT_SECRET_KEY=your_secret_key
JWT_EXPIRATION_HOURS=24

# Server Configuration
SERVER_PORT=8080
```

### 3. Running the Application

There are several ways to run the application:

#### Using Make (Recommended)

1. Initialize the project:
```bash
make init
```
This command will:
- Build the application
- Initialize the database
- Run necessary migrations
- Set up dependencies

2. Install Dependencies:
```bash
go mod tidy
go mod vendor
```

3. Run the application:
```bash
make run
```


### 4. Testing

Run the test suite:
```bash
make test
```

### 5. API Documentation

The API documentation is available in the `postman` directory. Import the collection into Postman to explore the available endpoints.

The server will start on `http://localhost:8080` by default.

## API Documentation

The API documentation is available in Postman format. Import the following files into Postman:

- `postman/Dating_App_API.postman_collection.json`
- `postman/Dating_App_API.postman_environment.json`

### Key Endpoints

1. Authentication
   - POST `/api/signup` - User registration
   - POST `/api/login` - User login

2. Profile Management
   - GET `/api/v1/profile` - Get user profile
   - GET `/api/v1/profile/discovery` - Get user profiles discovery (matching profiles by gender preferences)
   - PUT `/api/v1/profile` - Update profile
   - PUT `/api/v1/profile/preferences` - Update preferences

3. Swipe (Pass or Like)
   - POST `/api/v1/swipe` - Swipe action (like or pass)

4. Subscription (Pass or Like)
   - PUT `/api/v1/subscription` - Subscribe to premium package - unlock 10 swipe limiter

## Development

### Project Structure

```
.
├── cmd/
│   ├── app/        # Main application entry point
│   └── migrate/    # Database migration tool
├── internal/
│   ├── delivery/   # HTTP handlers and middleware
│   ├── common/     # Business domain models and interfaces
│   ├── config/     # Configuration and environment setup
│   ├── mock/       # Mock Unit Testing
│   ├── testutil/   # Testing helpers
│   └── services/   # Business logic implementation
├── common/         # Common configuration
│   ├── common.go   # Common functions
├── scripts/        # Helper scripts
└── postman/        # Postman collection and environment
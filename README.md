# Dating App API

A modern dating application API built with Go, featuring user authentication, profile management, and matching functionality.

## Prerequisites

Before running the application, make sure you have the following installed:

- Go 1.21 or higher
- PostgreSQL 14 or higher
- Redis 7 or higher
- Docker (optional, for containerized deployment)

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/nickyrolly/dealls-test.git
cd dealls-test
```

### 2. Environment Setup

Create a `.env` file in the root directory:

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=dating_app
DB_SSL_MODE=disable

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

### 3. Database Setup

#### Using Local PostgreSQL

1. Create the database:
```bash
createdb dating_app
```

2. Run migrations:
```bash
go run cmd/migrate/main.go
```

#### Using Docker

```bash
docker-compose up -d db redis
```

### 4. Running the Application

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

2. Run the application:
```bash
./app-dealls-test
```

#### Manual Setup

1. Install dependencies:
```bash
go mod tidy
go mod vendor
```

2. Build the application:
```bash
go build -o app-dealls-test cmd/app/*.go
```

3. Run the application:
```bash
./app-dealls-test
```

#### Using Docker

1. Build the Docker image:
```bash
docker build -t dealls-test .
```

2. Run the container:
```bash
docker run -p 8080:8080 --env-file .env dealls-test
```

### 5. Testing

Run the test suite:
```bash
make test
```

### 6. API Documentation

The API documentation is available in the `postman` directory. Import the collection into Postman to explore the available endpoints.

The server will start on `http://localhost:8080` by default.

## API Documentation

The API documentation is available in Postman format. Import the following files into Postman:

- `postman/Dating_App_API.postman_collection.json`
- `postman/Dating_App_API.postman_environment.json`

### Key Endpoints

1. Authentication
   - POST `/api/v1/auth/signup` - User registration
   - POST `/api/v1/auth/login` - User login

2. Profile Management
   - GET `/api/v1/profile` - Get user profile
   - PUT `/api/v1/profile` - Update profile
   - POST `/api/v1/profile/photos` - Upload profile photo

3. Matching System
   - GET `/api/v1/matches` - Get user matches
   - POST `/api/v1/profile/{id}/like` - Like a profile
   - POST `/api/v1/profile/{id}/pass` - Pass a profile

## Development

### Project Structure

```
.
├── cmd/
│   ├── api/        # Main application entry point
│   └── migrate/    # Database migration tool
├── internal/
│   ├── delivery/   # HTTP handlers and middleware
│   ├── domain/     # Business domain models and interfaces
│   ├── repository/ # Data access layer
│   └── services/   # Business logic implementation
├── pkg/            # Shared packages
└── scripts/        # Helper scripts
```

### Running Tests

```bash
go test ./...
```

### Code Generation

```bash
go generate ./...
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

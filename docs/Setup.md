# Dating App Setup Guide

## Prerequisites

- Go 1.19 or later
- PostgreSQL 12 or later
- Redis 6 or later
- Make (optional, for using Makefile commands)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/nickyrolly/dealls-test.git
cd dealls-test
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
cp .env.example .env
```

Edit the `.env` file with your configuration:
```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=dating_app

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT
JWT_SECRET=your_jwt_secret
JWT_EXPIRY=24h

# Server
PORT=8080
ENV=development
```

4. Run database migrations:
```bash
go run cmd/migrate/main.go
```

## Running the Application

### Development Mode
```bash
go run cmd/api/main.go
```

### Production Mode
```bash
go build -o api cmd/api/main.go
./api
```

## Testing

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Running Specific Tests
```bash
# Run tests in a specific package
go test ./internal/services/profile/...

# Run a specific test
go test ./... -run TestProfileEndpoints
```

## Development Workflow

1. Create a new branch for your feature:
```bash
git checkout -b feature/your-feature-name
```

2. Make your changes and write tests

3. Run tests and ensure they pass:
```bash
go test ./...
```

4. Format your code:
```bash
go fmt ./...
```

5. Commit your changes:
```bash
git add .
git commit -m "Description of your changes"
```

6. Push your changes and create a pull request:
```bash
git push origin feature/your-feature-name
```

## API Documentation

The API documentation is available at `/docs/API.md`. You can also access the live API documentation when running the server:

- Swagger UI: `http://localhost:8080/swagger/index.html`
- OpenAPI JSON: `http://localhost:8080/swagger/doc.json`

## Troubleshooting

### Common Issues

1. Database Connection Issues
```
Error: failed to connect to database
Solution: Check your database credentials in .env and ensure PostgreSQL is running
```

2. Redis Connection Issues
```
Error: failed to connect to Redis
Solution: Verify Redis is running and check your Redis configuration in .env
```

3. Permission Issues
```
Error: permission denied
Solution: Ensure your database user has the necessary permissions
```

### Getting Help

If you encounter any issues:

1. Check the logs for detailed error messages
2. Refer to the troubleshooting section in this guide
3. Open an issue on GitHub with:
   - Description of the problem
   - Steps to reproduce
   - Error messages
   - Your environment details

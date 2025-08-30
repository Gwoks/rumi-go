# RUMI Backend - Go API Server

A clean architecture Go backend service for the RUMI application with JWT authentication, MySQL database, and role-based access control.

## 🏗️ Architecture

This project follows clean architecture principles with the following structure:

```
rumi-go/
├── cmd/server/                 # Application entry point
├── internal/
│   ├── config/                # Configuration management
│   ├── database/              # Database connection
│   ├── domain/                # Business logic layer
│   │   ├── entity/           # Domain entities
│   │   ├── repository/       # Repository interfaces
│   │   └── usecase/          # Use case interfaces
│   ├── handler/              # HTTP handlers
│   ├── infrastructure/       # Infrastructure implementations
│   │   ├── repository/       # Repository implementations
│   │   └── usecase/          # Use case implementations
│   └── middleware/           # HTTP middleware
├── migrations/               # Database migrations
├── pkg/                     # Shared packages
│   ├── auth/               # Authentication utilities
│   ├── response/           # API response utilities
│   └── validator/          # Validation utilities
└── scripts/                # Utility scripts
```

## 🚀 Features

- **Clean Architecture**: Separation of concerns with clear boundaries
- **JWT Authentication**: Secure token-based authentication
- **Role-Based Access Control**: Admin and User roles
- **MySQL Database**: With sqlx for raw SQL queries
- **CORS Support**: Configurable cross-origin resource sharing
- **Environment Configuration**: Using .env files
- **Password Hashing**: Bcrypt for secure password storage
- **Graceful Shutdown**: Proper server shutdown handling
- **Database Migrations**: Schema versioning and migration system

## 🛠️ Prerequisites

- Go 1.21 or later
- MySQL 8.0 or later
- Git

## 📦 Installation

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd rumi-go
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Set up environment variables:**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Set up MySQL database:**
   ```sql
   CREATE DATABASE rumi_db;
   CREATE USER 'rumi_user'@'localhost' IDENTIFIED BY 'your_password';
   GRANT ALL PRIVILEGES ON rumi_db.* TO 'rumi_user'@'localhost';
   FLUSH PRIVILEGES;
   ```

5. **Run database migrations:**
   ```bash
   go run scripts/migrate.go
   ```

## 🏃 Running the Application

### Development Mode
```bash
go run cmd/server/main.go
```

### Production Mode
```bash
# Build the application
go build -o bin/server cmd/server/main.go

# Run the binary
./bin/server
```

The server will start on `http://localhost:8080` by default.

## 📚 API Endpoints

### Authentication Endpoints (matching Angular API service)

| Method | Endpoint                  | Description           | Auth Required |
|--------|---------------------------|-----------------------|---------------|
| POST   | `/api/v1/auth/login`      | User login           | No            |
| POST   | `/api/v1/auth/signup`     | User registration    | No            |
| POST   | `/api/v1/auth/logout`     | User logout          | Yes           |
| POST   | `/api/v1/auth/refresh`    | Refresh JWT token    | Yes           |
| POST   | `/api/v1/auth/validate`   | Validate JWT token   | No            |
| GET    | `/api/v1/auth/profile`    | Get user profile     | Yes           |

### System Endpoints

| Method | Endpoint                  | Description           | Auth Required |
|--------|---------------------------|-----------------------|---------------|
| GET    | `/health`                 | Health check         | No            |

### Example Requests

#### Register User
```bash
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "name": "John Doe",
    "phone": "+1234567890",
    "password": "password123"
  }'
```

#### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

#### Access Protected Endpoint
```bash
curl -X GET http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer <your-jwt-token>"
```

## 🔧 Environment Variables

| Variable                   | Description                    | Default              |
|---------------------------|--------------------------------|----------------------|
| `DB_HOST`                 | MySQL host                     | `localhost`          |
| `DB_PORT`                 | MySQL port                     | `3306`               |
| `DB_USER`                 | MySQL username                 | `root`               |
| `DB_PASSWORD`             | MySQL password                 | **Required**         |
| `DB_NAME`                 | MySQL database name            | `rumi_db`            |
| `SERVER_PORT`             | Server port                    | `8080`               |
| `SERVER_HOST`             | Server host                    | `0.0.0.0`            |
| `JWT_SECRET`              | JWT signing secret             | **Required in prod** |
| `JWT_EXPIRY_HOURS`        | JWT token expiry hours         | `24`                 |
| `ENVIRONMENT`             | Environment (dev/prod)         | `development`        |
| `BCRYPT_COST`             | Bcrypt hashing cost            | `12`                 |
| `CORS_ALLOWED_ORIGINS`    | CORS allowed origins           | `*`                  |

## 👤 Default Users

The migration creates two default users:

| Email            | Password   | Role    |
|------------------|------------|---------|
| admin@rumi.id    | admin123   | admin   |
| user@rumi.id     | admin123   | user    |

## 🔒 Security Features

- **Password Hashing**: Uses bcrypt with configurable cost
- **JWT Security**: HS256 algorithm with configurable secret
- **Role-Based Access**: Admin and User roles with middleware protection
- **Input Validation**: Request validation using Gin binding
- **CORS Protection**: Configurable CORS policies

## 🏗️ Development

### Code Structure

- **Entities**: Define business objects and validation rules
- **Repositories**: Define data access interfaces
- **Use Cases**: Implement business logic
- **Handlers**: Handle HTTP requests and responses
- **Middleware**: Provide cross-cutting concerns (auth, CORS, logging)

### Adding New Features

1. Define entities in `internal/domain/entity/`
2. Create repository interfaces in `internal/domain/repository/`
3. Create use case interfaces in `internal/domain/usecase/`
4. Implement repositories in `internal/infrastructure/repository/`
5. Implement use cases in `internal/infrastructure/usecase/`
6. Create handlers in `internal/handler/`
7. Add routes in `cmd/server/main.go`

## 🚀 Deployment

### Using Docker (optional)
```bash
# Build image
docker build -t rumi-backend .

# Run container
docker run -p 8080:8080 --env-file .env rumi-backend
```

### Environment Setup
- Set `ENVIRONMENT=production`
- Use strong `JWT_SECRET`
- Configure proper database credentials
- Set up proper CORS origins
- Use reverse proxy (nginx) for HTTPS

## 🤝 Contributing

1. Follow Go code conventions
2. Write tests for new features
3. Update documentation
4. Use meaningful commit messages

## 📝 License

This project is licensed under the MIT License.

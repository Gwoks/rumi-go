# rumi-go

A Go microservice following Clean Architecture principles with MySQL database, RBAC authentication, and comprehensive API endpoints.

## 🚀 Quick Start

### Prerequisites

- Go 1.23+
- MySQL 8.0+
- Make (optional, for using Makefile commands)

### Installation

1. **Clone the repository**
   ```bash
   git clone <your-repo-url>
   cd rumi-go
   ```

2. **Install Go dependencies**
   ```bash
   go mod download
   ```

3. **Install development tools**
   ```bash
   make install-tools
   ```

4. **Setup configuration**
   ```bash
   # Copy example config
   cp config.yaml.tmpl config.yaml
   
   # Edit config.yaml with your database settings
   # Update database host, port, user, password as needed
   ```

5. **Run database migrations**
   ```bash
   make migrate-mysql-up
   ```

6. **Start the server**
   ```bash
   make run
   ```

The server will start on port 8080 (configurable in `config.yaml`).

## 🗄️ Database Migrations

### Migration Commands

#### Using Makefile (Recommended)
```bash
# Run all pending migrations
make migrate-mysql-up

# Rollback the last migration
make migrate-mysql-down

# Check migration status
make migrate-mysql-status

# Create new migration
make migrate-create FILENAME=your_migration_name
```

#### Using CLI directly
```bash
# Run all pending migrations
go run cmd/server/main.go migrate:mysql:up

# Rollback the last migration
go run cmd/server/main.go migrate:mysql:down

# Check migration status
go run cmd/server/main.go migrate:mysql:status

# Create new migration
go run cmd/server/main.go migrate:create --filename your_migration_name
```

### Current Migrations

1. **001_create_users_table** - Creates the users table with RBAC support
2. **002_insert_admin_user** - Inserts the default admin user

### Default Admin User

After running migrations, you'll have an admin user with:
- **Email**: `admin@rumi.play`
- **Password**: `adminrumi`
- **Role**: `admin`

## 🏗️ Project Structure

```
rumi-go/
├── cmd/                          # Application entry points
│   ├── server/                   # Main server binary
│   │   └── main.go              # Server entry point
│   └── migrate/                  # Migration binary
│       └── main.go              # Migration entry point
├── internal/                     # Internal application code
│   ├── config/                   # Configuration management
│   │   └── config.go            # Viper-based config loader
│   ├── database/                 # Database configuration
│   │   └── config.go            # Database config structs
│   ├── infrastructure/           # External dependencies
│   │   ├── infrastructure.go     # Infrastructure initialization
│   │   └── database/             # Database layer
│   │       ├── databasestore.go  # Database store interface
│   │       └── user.go          # User data access
│   ├── migration/                # Database migrations
│   │   ├── migration.go         # Migration interface
│   │   └── runner.go            # MySQL migration runner
│   ├── model/                    # Data models
│   │   └── user.go              # User model and interfaces
│   ├── presenter/                # Interface adapters
│   │   ├── console/              # CLI handlers
│   │   │   ├── console.go       # Console interface
│   │   │   ├── server.go        # HTTP server setup
│   │   │   ├── migration.go     # Migration commands
│   │   │   └── config.go        # Configuration commands
│   │   └── rest/                 # HTTP handlers
│   │       └── handlers/         # REST API handlers
│   │           ├── base.go       # Base handler structure
│   │           ├── create.go     # User creation
│   │           ├── login.go      # User login
│   │           ├── get.go        # User retrieval
│   │           ├── update.go     # User updates
│   │           └── delete.go     # User deletion
│   ├── usecase/                  # Business logic layer
│   │   └── user/                 # User domain
│   │       ├── user.go          # Base user usecase
│   │       ├── create.go        # User creation logic
│   │       ├── login.go         # User login logic
│   │       ├── get.go           # User retrieval logic
│   │       ├── update.go        # User update logic
│   │       └── delete.go        # User deletion logic
│   └── utils/                    # Internal utilities
│       └── sqlx.go              # SQLX database utilities
├── migrations/                    # Database migration files
│   ├── 001_create_users_table.up.sql
│   ├── 001_create_users_table.down.sql
│   ├── 002_insert_admin_user.up.sql
│   └── 002_insert_admin_user.down.sql
├── docs/                         # Generated Swagger documentation
├── config.yaml                   # Configuration file
├── config.yaml.tmpl              # Configuration template
├── go.mod                        # Go module file
├── go.work                       # Go workspace file
├── Makefile                      # Build and development commands
├── .mockery.yaml                 # Mockery configuration
└── README.md                     # This file
```

## 🛠️ Development Commands

### Build and Run
```bash
# Compile the application
make compile

# Run the server
make run

# Run with debug logging
env LOG_LEVEL=debug make run
```

### Testing
```bash
# Run all tests
make test

# Generate mocks
make generate-mock

# Run tests with coverage
make cover
```

### Code Quality
```bash
# Run linter
make lint

# Run pre-commit hooks
pre-commit run --all-files
```

### Configuration
```bash
# Show current configuration
make config-show

# Copy config template
make copy-config
```

## 🔧 Configuration

The application uses Viper for configuration management with the following sources (in order of precedence):

1. **Environment variables** (with `RUMI_` prefix)
2. **Configuration file** (`config.yaml`)
3. **Default values**

### Configuration File Structure

```yaml
app:
  name: "rumi-go"
  env: "local"
  api_prefix: "rumi-go"

server:
  port: 8080

database:
  driver: mysql
  name: rumi
  host: localhost
  port: 3306
  user: root
  password: root
  max_open: 50
  max_idle: 10
  max_life_time: 3m
  max_idle_time: 3m
  statement_timeout: 3s

redis:
  address: localhost
  port: "6379"
  password: ""
  db: 0
```

### Environment Variables

You can override any configuration using environment variables:

```bash
export RUMI_DATABASE_PASSWORD=your_password
export RUMI_SERVER_PORT=9090
export RUMI_DATABASE_PORT=3306
```

## 🌐 API Endpoints

### Health Check
- `GET /ping` - Health check endpoint

### User Management
- `POST /rumi-go/internal/v1/user` - Create user
- `POST /rumi-go/internal/v1/user/login` - User login
- `GET /rumi-go/internal/v1/user?id={id}` - Get user by ID
- `GET /rumi-go/internal/v1/user/email?email={email}` - Get user by email
- `GET /rumi-go/internal/v1/users` - List users
- `PUT /rumi-go/internal/v1/user?id={id}` - Update user
- `DELETE /rumi-go/internal/v1/user?id={id}` - Delete user

### Swagger Documentation
- `GET /swagger/index.html` - Interactive API documentation

## 🔐 Authentication & Authorization

The application implements:
- **JWT-based authentication** for API endpoints
- **RBAC (Role-Based Access Control)** with roles: admin, user, manager
- **Password hashing** using bcrypt
- **Session management** via database

## 🧪 Testing

The project includes comprehensive testing:
- **Unit tests** for all business logic
- **Mock generation** using Mockery
- **Table-driven tests** for edge cases
- **HTTP testing** using httptest
- **Test fixtures** and helpers

## 🐳 Docker Support

### Build Docker Image
```bash
docker build -t rumi-go .
```

### Run with Docker Compose
```bash
docker-compose up -d
```

## 📚 Additional Documentation

- **API Examples**: See `curl-samples.md` for API usage examples
- **Migration Details**: Database schema and migration information
- **Architecture**: Clean Architecture implementation details

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## 📄 License

[Add your license information here]

## 🆘 Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Ensure MySQL is running
   - Check database credentials in `config.yaml`
   - Verify database exists and is accessible

2. **Port Already in Use**
   ```bash
   # Kill process using port 8080
   lsof -ti:8080 | xargs kill -9
   ```

3. **Migration Errors**
   - Check database connection
   - Verify migration files are in `migrations/` directory
   - Use `make migrate-mysql-status` to check status

4. **Configuration Issues**
   - Use `make config-show` to verify loaded configuration
   - Check `config.yaml` syntax
   - Verify environment variables are set correctly

### Getting Help

- Check the logs for detailed error messages
- Use `make config-show` to verify configuration
- Ensure all dependencies are installed
- Verify database connectivity
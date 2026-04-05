# Dependencies & Requirements

## Go Version
- **Required**: Go 1.24.0 or higher
- **Recommended**: Latest stable Go version

## External Dependencies

### Core Framework
- **github.com/go-chi/chi/v5** (v5.2.5)
  - HTTP router and middleware
  - Lightweight and composable
  - Used for: REST API routing

### Database
- **github.com/go-sql-driver/mysql** (v1.9.3)
  - MySQL driver for Go
  - Pure Go implementation
  - Used for: Database connectivity

### Database Migrations
- **github.com/golang-migrate/migrate/v4** (v4.19.1)
  - Database migration tool
  - Supports versioning and rollback
  - Used for: Schema management

### Logging
- **go.uber.org/zap** (v1.27.1)
  - Structured logging library
  - High performance
  - Used for: Application logging

### Utility
- **go.uber.org/multierr** (v1.10.0)
  - Error handling utilities
  - Dependency of zap
  - Used for: Error aggregation

## Installation

### Install Go Dependencies
```bash
go mod download
go mod tidy
```

### Install Database Tools

#### Option 1: Using Homebrew (macOS)
```bash
brew install golang-migrate
```

#### Option 2: Using Go
```bash
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

#### Option 3: Docker
```bash
# Already included in docker-compose.yml
docker-compose up mysql
```

### Install MySQL Client (optional, for manual testing)
```bash
# macOS
brew install mysql-client

# Ubuntu/Debian
sudo apt-get install mysql-client

# Or use docker
docker run -it --rm mysql:8.0 mysql -h host.docker.internal -u root -p
```

## System Requirements

### Minimum
- **CPU**: 1 core
- **Memory**: 256 MB
- **Disk**: 100 MB

### Recommended
- **CPU**: 2+ cores
- **Memory**: 1 GB+
- **Disk**: 1 GB+

## Operating Systems
- ✅ Linux (Ubuntu, Debian, CentOS, etc.)
- ✅ macOS (Intel and Apple Silicon)
- ✅ Windows (WSL2 recommended)

## Database Requirements

### MySQL Version
- **Required**: 8.0+
- **Recommended**: 8.0.27+ (LTS)

### MySQL Connection
- **User**: root (or any user with permissions)
- **Password**: root (default, change in production)
- **Host**: localhost:3306
- **Database**: booking_api

### MySQL Features Required
- ✅ InnoDB storage engine
- ✅ Foreign key constraints
- ✅ UNIQUE constraints
- ✅ CHECK constraints
- ✅ Transactions
- ✅ TIMESTAMP support

### MySQL Configuration
No special configuration required. Default MySQL 8.0 settings work fine.

## Docker & Container Setup

### Docker Version
- **Docker**: 20.10+
- **Docker Compose**: 1.29+

### Docker Images Used
```yaml
services:
  mysql:
    image: mysql:8.0
    # Automatically pulls from Docker Hub
```

### Container Resources
- **Memory**: 512 MB minimum for MySQL
- **Disk**: 1 GB for database persistence

## Development Tools (Optional)

### Code Quality
```bash
# Formatter
go fmt ./...

# Linter
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Vet tool
go vet ./...
```

### Testing
```bash
# Built-in testing
go test ./...

# With coverage
go test -cover ./...

# With verbosity
go test -v ./...
```

### Profiling
```bash
# CPU profiling
go test -cpuprofile=cpu.prof ./...

# Memory profiling
go test -memprofile=mem.prof ./...
```

## Port Requirements

- **Application Port**: 8080 (configurable via PORT env var)
- **MySQL Port**: 3306 (default)

Make sure these ports are available or configure them in `.env`

## Environment Setup

### Create .env file
```bash
cp .env.example .env
```

### Edit .env with your settings
```env
DATABASE_URL=root:root@tcp(localhost:3306)/booking_api?parseTime=true
PORT=8080
ENVIRONMENT=development
LOG_LEVEL=info
```

## Verification Checklist

Before running the application:

- [ ] Go 1.24+ installed
- [ ] MySQL 8.0+ running
- [ ] Dependencies downloaded (`go mod download`)
- [ ] Database created (`CREATE DATABASE booking_api`)
- [ ] Migrations applied (`make migrate-up`)
- [ ] Application builds (`make build`)
- [ ] Tests pass (`make test`)
- [ ] Server starts (`make run`)

## Troubleshooting Dependencies

### "package not found" error
```bash
go mod download
go mod tidy
go mod verify
```

### MySQL connection refused
```bash
# Check MySQL is running
docker-compose ps

# Or check local MySQL
ps aux | grep mysql

# Verify credentials in DATABASE_URL
```

### Port already in use
```bash
# Find process using port
lsof -i :8080

# Kill process
kill -9 PID

# Or use different port
PORT=8081 make run
```

### Go version mismatch
```bash
# Check current Go version
go version

# If using Go version manager (like gvm or asdf)
# Update to 1.24.0 or later
```

## Dependency Licensing

All dependencies use permissive licenses:
- **Chi**: MIT License
- **MySQL Driver**: Mozilla Public License 2.0
- **golang-migrate**: BSD 3-Clause License
- **Zap**: MIT License

See individual package repositories for detailed license information.

## Update Dependencies

### Check for updates
```bash
go list -u -m all
```

### Update all dependencies
```bash
go get -u ./...
go mod tidy
```

### Update specific dependency
```bash
go get -u github.com/package/name
```

## Production Considerations

### Database
- Use managed MySQL service (AWS RDS, Google Cloud SQL, etc.)
- Enable automated backups
- Use encryption at rest and in transit
- Set up replication/failover

### Logging
- Configure centralized logging (ELK, Datadog, etc.)
- Set appropriate log levels
- Monitor error rates

### Monitoring
- Add Prometheus metrics
- Set up alerts
- Monitor database performance
- Track API response times

### Security
- Use strong database passwords
- Enable SSL/TLS for connections
- Implement rate limiting
- Add authentication/authorization

## Docker Production Setup

For production deployments with Docker:

```dockerfile
FROM golang:1.24 AS builder
WORKDIR /app
COPY . .
RUN go build -o booking-api cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/booking-api .
EXPOSE 8080
CMD ["./booking-api"]
```

Build and run:
```bash
docker build -t booking-api:latest .
docker run -e DATABASE_URL="..." -p 8080:8080 booking-api:latest
```

## Support Matrix

| Component | Version | Status |
|-----------|---------|--------|
| Go | 1.24+ | ✅ Tested |
| MySQL | 8.0+ | ✅ Tested |
| Chi | 5.2.5 | ✅ Tested |
| Zap | 1.27.1 | ✅ Tested |
| Docker | 20.10+ | ✅ Tested |
| Docker Compose | 1.29+ | ✅ Tested |

## Getting Help

1. Check the [README.md](README.md) for API documentation
2. Review [QUICKSTART.md](QUICKSTART.md) for setup help
3. Check [ARCHITECTURE.md](ARCHITECTURE.md) for design details
4. Review test files for usage examples
5. Check error logs for specific issues

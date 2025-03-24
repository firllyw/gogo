# User Service

A simple microservices-based application using Go for managing users.

## Features

- RESTful API for user management (Create, Read, Update, Delete)
- PostgreSQL database integration
- Structured logging
- Clean architecture with separation of concerns
- Docker support for easy deployment
- Optional OAuth authentication

## Project Structure

```
user-service/
├── cmd/
│   └── server/
│       └── main.go             # Application entry point
├── config/
│   └── config.go               # Configuration management
├── internal/
│   ├── api/
│   │   ├── handlers/           # HTTP request handlers
│   │   ├── middleware/         # HTTP middleware
│   │   └── routes/             # API route definitions
│   ├── domain/                 # Business domain models and interfaces
│   ├── repository/             # Data access implementations
│   └── service/                # Business logic
├── migrations/                 # Database migration scripts
├── pkg/                        # Shared packages
│   ├── db/                     # Database connections
│   └── validator/              # Input validation
├── .env                        # Environment variables (not in git)
├── .env.example                # Example environment variables
├── Dockerfile                  # Docker configuration
├── docker-compose.yml          # Docker Compose configuration
└── README.md                   # Project documentation
```

## Prerequisites

- Go 1.21 or higher
- PostgreSQL
- Docker and Docker Compose (optional, for containerized deployment)

## Getting Started

### Local Development

1. Clone the repository:

```sh
git clone https://github.com/firllyw/gogo.git
cd gogo
```

2. Install dependencies:

```sh
go mod download
```

3. Set up your environment variables:

```sh
cp .env.example .env
# Edit .env file with your configuration
```

4. Start a PostgreSQL instance:

5. Run the application:

```sh
go run cmd/server/main.go
```

## API Endpoints

### User Management

| Method | Endpoint         | Description        |
|--------|------------------|--------------------|
| POST   | /api/v1/users    | Create a new user  |
| GET    | /api/v1/users/:id| Get user by ID     |
| PUT    | /api/v1/users/:id| Update user        |
| DELETE | /api/v1/users/:id| Delete user        |

### Health Check

| Method | Endpoint  | Description   |
|--------|-----------|---------------|
| GET    | /health   | Service health|

## API Examples

### Create a User

```sh
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jake Peralta",
    "email": "peralta@b99.com"
  }'
```

### Get a User

```sh
curl -X GET http://localhost:8080/api/v1/users/123e4567-e89b-12d3-a456-426614174000
```

### Update a User

```sh
curl -X PUT http://localhost:8080/api/v1/users/123e4567-e89b-12d3-a456-426614174000 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Amy Santiago",
    "email": "amy@b99.com"
  }'
```

### Delete a User

```sh
curl -X DELETE http://localhost:8080/api/v1/users/123e4567-e89b-12d3-a456-426614174000
```

## OAuth Authentication

This service includes OAuth 2.0 authentication with JWT tokens. To configure it:

1. Set `OAUTH_ENABLED=true` in your `.env` file.
2. Set the OAuth provider credentials:
   ```
   OAUTH_CLIENT_ID=your-client-id
   OAUTH_CLIENT_SECRET=your-client-secret
   OAUTH_REDIRECT_URL=http://localhost:8080/api/v1/auth/callback
   ```
3. Set a strong JWT secret key:
   ```
   JWT_SECRET=your-secure-random-key
   ```

### Authentication Flow

1. Users visit `/api/v1/auth/login` to initiate the OAuth flow
2. After successful authentication, they are redirected to the callback URL
3. The callback handler exchanges the authorization code for an access token
4. User information is retrieved and a JWT token is generated
5. The JWT token must be included in subsequent API requests as a Bearer token:
   ```
   Authorization: Bearer your-jwt-token
   ```

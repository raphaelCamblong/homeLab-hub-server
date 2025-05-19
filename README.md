# Homelab Server

![Go](https://img.shields.io/badge/Go-1.20-blue.svg)
![License](https://img.shields.io/badge/License-MIT-green.svg)

## Overview

Homelab Server is a Go-based application designed to manage and serve various functionalities for your home lab. It is built following the principles of Clean Architecture to ensure maintainability, testability, and scalability.

## Table of Contents

- [Features](#features)
- [Getting Started](#getting-started)
- [Building and Running](#building-and-running)
- [Directory Structure](#directory-structure)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Features

- **Clean Architecture**: Ensures a scalable and maintainable codebase.
- **Fast and Lightweight**: Built with Go, known for its performance and efficiency.
- **Extensible**: Easily add new features and modules.
- **Secure**: Implements best practices for security in a home lab environment.

## Getting Started

### Prerequisites

Ensure you have the following installed:

- Go 1.20 or later
- Docker (for containerized builds)
- `just` (a handy command runner)

### Installation

Clone the repository:

```sh
git clone https://github.com/yourusername/homelab_server.git
cd homelab_server
```

## API:
```
Public Routes:
/api/v1/auth/
  ├── POST /login
  └── POST /register

Protected Routes (requires JWT):
/api/v1/auth/
  └── POST /logout

/api/v1/users/
  ├── GET /                         (requires users:read)
  ├── GET /:id                      (requires users:read)
  ├── GET /:id/roles               (requires roles:read)
  ├── POST /:id/roles/:roleId      (requires roles:write)
  └── DELETE /:id/roles/:roleId    (requires roles:write)

/api/v1/api-keys/
  ├── POST /
  ├── GET /
  └── DELETE /:keyId

/api/v1/audit-logs                 (requires audit:read)

/api/v1/services
  ├──GET    /          - List all services (requires services:read)
  ├──GET    /:id      - Get a specific service (requires services:read)
  ├──POST   /         - Create a new service (requires services:write)
  ├──PUT    /:id      - Update a service (requires services:write)
  └──DELETE /:id      - Delete a service (requires services:write)
```

## OpenAPI Documentation

This project uses [swaggo/swag](https://github.com/swaggo/swag) to generate OpenAPI documentation.

 backend api available at: http://127.0.0.1:8080/swagger/index.html#/

### Installation

1. Install swag CLI:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

2. Add Swagger dependencies to your project:
```bash
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

### Generate Documentation

1. Add Swagger annotations to your handlers (already done)
2. Generate Swagger docs:
```bash
swag init -g cmd/server/main.go
```

This will create a `docs` directory containing the generated Swagger documentation.

### Access Documentation UI

After starting the server, access the Swagger UI at:
```
http://localhost:8080/swagger/index.html
```

### Example Swagger Integration

In your main.go:
```go
import (
    "github.com/gin-gonic/gin"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
    _ "homelab.com/homelab-server/docs" // generated swagger docs
)

func main() {
    r := gin.Default()
    
    // Swagger documentation endpoint
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    
    // ... rest of your setup
}
```

### Updating Documentation

After making changes to the API or annotations, regenerate the documentation:
```bash
swag init -g cmd/server/main.go
```
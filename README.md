# Health API

The Health API is a RESTful service built with [Gin](https://github.com/gin-gonic/gin) in Go, designed to manage and provide health-related data. It offers endpoints for user authentication, patient records, and other health-related functionalities.

## Table of Contents

- [Features](#features)
- [Technologies Used](#technologies-used)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Running the Application](#running-the-application)
- [API Endpoints](#api-endpoints)
- [Project Structure](#project-structure)
- [Contributing](#contributing)
- [License](#license)

## Features

- **User Authentication**: Secure login and registration system.
- **Patient Management**: CRUD operations for patient records.
- **Middleware Support**: Includes JWT authentication middleware.
- **Database Migrations**: Laravel-style migration and seeder commands.
- **Structured Project Layout**: Organized into controllers, services, models, and routes.

## Technologies Used

- **Go**: Programming language.
- **Gin**: Web framework for Go.
- **GORM**: ORM library for Go.
- **JWT**: JSON Web Tokens for authentication.

## Prerequisites

- [Go](https://golang.org/doc/install) 1.16 or higher
- [MySQL](https://www.mysql.com/downloads/) or any compatible database

## Installation

1. **Clone the repository**:

   ```bash
   git clone https://github.com/mkawsar/health-api.git
2. **Install dependencies**:

   ```bash
   cd health-api
   go mod tidy

## Configuration

1. Environment Variables: Copy the **.env.example** file to **.env** and update the values accordingly.

   ```bash
   cp .env.example .env
   ```

## Docker

Run the API with MongoDB and Redis using Docker Compose:

```bash
# Copy env and start (create .env from .env.example first)
cp .env.example .env
docker compose up --build -d
```

The API is at [http://localhost:8080](http://localhost:8080). MongoDB is on **port 27018** (host) so it doesn’t conflict with local MongoDB; connect with `mongodb://localhost:27018`. Redis is on **port 6380** (host) so it doesn’t conflict with local Redis; connect with `localhost:6380`.

**One image for dev and deploy:** The same Dockerfile builds one image. Use `.env` (or compose) to set `MODE=debug` for development and `MODE=release` for deployment. Optional: `docker build --build-arg TARGET=development -t health-api .` for faster dev builds (no binary stripping).

## Database Migrations

This project includes Laravel-style migration and seeder commands for managing your database.

### Migration Commands

**Run Migrations:**
```bash
docker exec health-api go run cmd/migrate/main.go -command migrate
```

### Seeder Commands

**Run All Seeders:**
```bash
docker exec health-api go run cmd/seed/main.go 
```

**Run Specific Seeder:**
```bash
docker exec health-api go run cmd/seed/main.go -seeder user_seeder
```

For more detailed information, see [MIGRATION_COMMANDS.md](MIGRATION_COMMANDS.md).

## 📁 Project Directory Structure

```bash
health-api/ 
  
    ├── cmd/ # CLI commands (migrate, seed)
    │   ├── migrate/ # Migration command
    │   └── seed/ # Seeder command
    
    ├── controllers/ # HTTP handlers 
  
    ├── docs/ # API documentation (Swagger, Postman, etc.) 
    
    ├── middlewares/ # Custom middleware (e.g., JWT auth) 
    
    ├── migrations/ # Database migration SQL files
    
    ├── models/ # Database models (GORM) 
    
    ├── routes/ # API route groupings 
    
    ├── seeders/ # Database seeder files
    
    ├── services/ # Business logic 
    
    ├── utils/ # Helper functions and utilities 
    
    ├── .env.example # Environment variable example 
    
    ├── .gitignore # Git ignored files list 
    
    ├── go.mod # Go module definition 
    
    ├── go.sum # Go module checksums 
    
    ├── Makefile # Make commands for migrations and seeders
    
    ├── MIGRATION_COMMANDS.md # Detailed migration documentation
    
    └── main.go # Entry point of the application
```

## Contributing
Contributions are welcome! Please follow these steps:

    1. Fork the repository.
    2. Create a new branch for your feature or bug fix.
    3. Make your changes and commit them.
    4. Push your changes to your fork.
    5. Create a pull request to the main repository.

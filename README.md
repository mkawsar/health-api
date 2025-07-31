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

## Running the Application
To run the application in development mode with live-reloading, you can use [Air](https://github.com/cosmtrek/air):
1. **Install Air**:

   ```bash
   go install github.com/cosmtrek/air@v1.40.4
   ```
2. **Start Air**:

   ```bash
   air
   ```
3. **Start the Application**:

   ```bash
   go run main.go
   ```
4. **Access the Application**:

   Open your browser and navigate to [http://localhost:8080](http://localhost:8080) to access the application.

## 📁 Project Directory Structure

```bash
health-api/ 
  
    ├── config/ # App configuration (DB, env) 
    
    ├── controllers/ # HTTP handlers 
  
    ├── docs/ # API documentation (Swagger, Postman, etc.) 
    
    ├── middlewares/ # Custom middleware (e.g., JWT auth) 
    
    ├── models/ # Database models (GORM) 
    
    ├── routes/ # API route groupings 
    
    ├── services/ # Business logic 
    
    ├── utils/ # Helper functions and utilities 
    
    ├── .env.example # Environment variable example 
    
    ├── .gitignore # Git ignored files list 
    
    ├── go.mod # Go module definition 
    
    ├── go.sum # Go module checksums 
    
    └── main.go # Entry point of the application
```

## Contributing
Contributions are welcome! Please follow these steps:

    1. Fork the repository.
    2. Create a new branch for your feature or bug fix.
    3. Make your changes and commit them.
    4. Push your changes to your fork.
    5. Create a pull request to the main repository.

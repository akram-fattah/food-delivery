# Food Delivery API

A RESTful API for a food delivery application built with Go and PostgreSQL.

## Overview

This project is a backend API for a food delivery service that handles user management, authentication, and order processing. It uses Go for the backend logic and PostgreSQL for data persistence.

## Tech Stack

- **Go** - Backend programming language
- **PostgreSQL** - Relational database
- **godotenv** - Environment variable management
- **pq** - PostgreSQL driver for Go

## Prerequisites

Before you begin, ensure you have the following installed:

- [Go](https://golang.org/doc/install) (version 1.26.1 or later)
- [PostgreSQL](https://www.postgresql.org/download/) (version 12 or later)
- Git

## Project Structure

```
food-delivery/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   └── database/
│       └── db.go                # Database connection logic
├── delivery_db.sql              # Database schema
├── go.mod                       # Go module dependencies
├── go.sum                       # Go module checksums
├── .gitignore                   # Git ignore rules
└── README.md                    # This file
```

## Database Setup

### 1. Create the Database

Connect to PostgreSQL and create the database:

```sql
CREATE DATABASE delivery_db;
```

Or using the command line:

```bash
# For macOS/Linux
psql -U postgres -c "CREATE DATABASE delivery_db;"

# For Windows
psql -U postgres -c "CREATE DATABASE delivery_db;"
```

### 2. Run the Schema

Execute the SQL script to create the necessary tables:

```bash
# For macOS/Linux
psql -U postgres -d delivery_db -f delivery_db.sql

# For Windows
psql -U postgres -d delivery_db -f delivery_db.sql
```

### Database Schema

The current schema includes:

**Users Table:**
| Field | Type | Description |
|-------|------|-------------|
| id | BIGSERIAL | Primary key, auto-increment |
| username | VARCHAR(50) | Unique username |
| name | VARCHAR(100) | Full name |
| email | VARCHAR(150) | Unique email address |
| password | VARCHAR(255) | Hashed password |
| phone | VARCHAR(20) | Phone number (optional) |
| address | TEXT | User address (optional) |
| role | VARCHAR(20) | User role: 'user', 'admin', or 'Delivery' |
| is_verified | BOOLEAN | Email verification status |
| verification_code | VARCHAR(100) | Email verification code |
| created_at | TIMESTAMP | Account creation time |
| updated_at | TIMESTAMP | Last update time (auto-updated) |

## Environment Configuration

### 1. Create Environment File

Create a `.env` file in the root directory of the project:

```bash
touch .env
```

### 2. Configure Environment Variables

Add the following variables to your `.env` file:

```env
# Application Environment
APP_ENV=development

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=delivery_db
```

**Note:** Replace `your_password` with your actual PostgreSQL password. If you're using the default `postgres` user without a password, you can leave `DB_PASSWORD` empty or omit it.

### Environment Variables Reference

| Variable | Description | Required |
|----------|-------------|----------|
| APP_ENV | Application environment (development/production) | No (defaults to development) |
| DB_HOST | PostgreSQL host address | Yes |
| DB_PORT | PostgreSQL port | Yes |
| DB_USER | PostgreSQL username | Yes |
| DB_PASSWORD | PostgreSQL password | No |
| DB_NAME | Database name | Yes |

## Installation & Running

### 1. Clone the Repository

```bash
git clone https://github.com/akram-fattah/food-delivery.git
cd food-delivery
```

### 2. Install Dependencies

```bash
go mod download
```

Or:

```bash
go mod tidy
```

### 3. Run the Application

```bash
go run cmd/api/main.go
```

If everything is configured correctly, you should see:

```
2024/XX/XX XX:XX:XX Database connected successfully
Hello, Akram!
```

## Building for Production

To build the application for production:

```bash
go build -o food-delivery-api cmd/api/main.go
```

Then run the binary:

```bash
./food-delivery-api
```

## Features

### Current Implementation

- Database connection management with connection pooling
- Environment-based configuration
- Structured logging
- Connection timeout handling (5 seconds)
- Database health check on startup

### Database Features

- Auto-incrementing primary keys
- Automatic timestamp updates via triggers
- Data validation constraints (check constraints, unique constraints)
- Role-based user management

## API Development Status

This project is currently in active development. The database layer is fully configured and ready for API endpoint implementation.

## Troubleshooting

### Database Connection Issues

**Problem:** `Database connection failed` or `Database unreachable`

**Solutions:**
1. Verify PostgreSQL is running: `sudo service postgresql status` (Linux) or check Services (Windows)
2. Check your `.env` file has correct credentials
3. Ensure the database `delivery_db` exists
4. Verify the PostgreSQL port (default: 5432)

### Environment File Not Loading

**Problem:** `Warning: .env file not loaded`

**Solutions:**
1. Ensure the `.env` file is in the project root
2. Check file permissions
3. Verify the file is named exactly `.env` (not `.env.txt` or similar)

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/new-feature`)
3. Commit your changes (`git commit -m 'Add new feature'`)
4. Push to the branch (`git push origin feature/new-feature`)
5. Open a Pull Request

## License

This project is private and proprietary.

---

**Author:** Akram Fattah  
**Repository:** https://github.com/akram-fattah/food-delivery

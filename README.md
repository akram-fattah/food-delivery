# 🍔 Food Delivery API 🚀

A RESTful API for a food delivery application built with Go and PostgreSQL.

## Overview

This project is a backend API for a food delivery service that handles user management, authentication, meal management, category management, and order processing. It uses Go for the backend logic and PostgreSQL for data persistence.

## Tech Stack

- **Go** - Backend programming language
- **PostgreSQL** - Relational database
- **godotenv** - Environment variable management
- **pq** - PostgreSQL driver for Go
- **JWT** - JSON Web Tokens for authentication
- **bcrypt** - Password hashing

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
│       └── main.go              # 🚦 Application entry point & HTTP routes
├── internal/
│   ├── database/
│   │   ├── auth.go              # 🔐 Authentication DB logic
│   │   ├── category_repository.go # 📂 Category DB operations
│   │   ├── database.go          # 🗄️ Database connection logic
│   │   ├── meal_repoditory.go   # 🍽️ Meal DB operations
│   │   ├── refresh_token_repository.go # 🔄 Refresh token DB logic
│   │   └── user_repository.go   # 👤 User DB logic
│   ├── handlers/
│   │   ├── category.go          # 📂 Category endpoints
│   │   ├── login.go             # 🔑 /auth/login endpoint
│   │   ├── logout.go            # 🚪 /auth/logout endpoint
│   │   ├── meals.go             # 🍽️ Meal endpoints
│   │   ├── refresh.go           # ♻️ /auth/refresh endpoint
│   │   ├── register.go          # 📝 /auth/register endpoint
│   │   ├── reset_password.go    # 🔒 /auth/reset-password endpoint
│   │   ├── update_password.go   # 🛠️ /auth/update-password endpoint
│   │   ├── verifyEmail.go       # ✅ /auth/verify-email endpoint
│   │   └── verify_reset_code.go # 🔍 /auth/verify-reset-code endpoint
│   ├── helper/
│   │   ├── code_helper.go       # 🧩 Code generation utilities
│   │   ├── config.go            # ⚙️ Config helpers
│   │   ├── email_helper.go      # 📧 Email sending logic
│   │   ├── response_helper.go   # 📦 JSON/error helpers
│   │   ├── security_helper.go   # 🛡️ Password hashing & validation
│   │   └── user.go              # 👤 User validation helpers
│   └── models/
│       ├── category.go          # 📂 Category struct definition
│       ├── meal.go              # 🍽️ Meal struct definition
│       └── user.go              # 🧑‍💻 User struct definition
├── uploads/                     # 🖼️ Image storage for meal photos
├── delivery_db.sql              # 🗃️ Database schema
├── go.mod                       # 📦 Go module dependencies
├── go.sum                       # 📑 Go module checksums
├── .env                         # 🌱 Environment variables
└── README.md                    # 📖 This file
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

#### 👤 Users Table

| 🏷️ Field | 🗃️ Type | 📝 Description |
|----------|---------|----------------|
| id | BIGSERIAL | Primary key, auto-increment |
| username | VARCHAR(50) | Unique username |
| name | VARCHAR(100) | Full name |
| email | VARCHAR(150) | Unique email address |
| password | VARCHAR(255) | Hashed password |
| phone | VARCHAR(20) | Phone number (optional) |
| address | TEXT | User address (optional) |
| role | VARCHAR(20) | User role: 'user', 'admin', 'Delivery' |
| is_verified | BOOLEAN | Email verification status |
| verification_code | VARCHAR(100) | Email verification code |
| created_at | TIMESTAMP | Account creation time |
| updated_at | TIMESTAMP | Last update time (auto-updated) |
| verification_expires | TIMESTAMP | Verification code expiration (24h) |

#### 🔄 Refresh Tokens Table

| 🏷️ Field | 🗃️ Type | 📝 Description |
|----------|---------|----------------|
| id | SERIAL | Primary key, auto-increment |
| user_id | INTEGER | References users(id), owner of the token |
| token | TEXT | The refresh token value |
| expires_at | TIMESTAMP | Expiration date/time of the refresh token |
| created_at | TIMESTAMP | Creation timestamp (default: now) |

**How it works:**
- Each refresh token is linked to a user and is used to obtain new access tokens without re-authenticating.
- When a user logs out, the refresh token is deleted from this table.
- Expired tokens are rejected automatically.

#### 📂 Categories Table

| 🏷️ Field | 🗃️ Type | 📝 Description |
|----------|---------|----------------|
| id | BIGSERIAL | Primary key, auto-increment |
| category_name | VARCHAR(100) | Unique category name |
| created_at | TIMESTAMP | Creation timestamp |

#### 🍽️ Meals Table

| 🏷️ Field | 🗃️ Type | 📝 Description |
|----------|---------|----------------|
| id | BIGSERIAL | Primary key, auto-increment |
| name | VARCHAR(150) | Meal name |
| description | TEXT | Meal description |
| price | NUMERIC(10,2) | Meal price |
| image_url | TEXT | Path to meal image |
| is_available | BOOLEAN | Availability status (default: true) |
| category_id | BIGINT | References categories(id) |
| created_at | TIMESTAMP | Creation timestamp |

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

# JWT Configuration
JWT_SECRET=your_jwt_secret_key_here

# Email Configuration (SMTP)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_EMAIL=your_email@gmail.com
SMTP_PASSWORD=your_app_password
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
| JWT_SECRET | Secret key for JWT token signing | Yes |
| SMTP_HOST | SMTP server host | Yes (for email) |
| SMTP_PORT | SMTP server port | Yes (for email) |
| SMTP_EMAIL | SMTP email address | Yes (for email) |
| SMTP_PASSWORD | SMTP email password | Yes (for email) |

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

### 3. Create Uploads Directory

```bash
mkdir -p uploads
```

### 4. Run the Application

```bash
go run cmd/api/main.go
```

If everything is configured correctly, you should see:

```
2024/XX/XX XX:XX:XX Database connected successfully
Server running on port :8000
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

- ✅ Database connection management with connection pooling
- ✅ Environment-based configuration
- ✅ Structured logging
- ✅ Connection timeout handling (5 seconds)
- ✅ Database health check on startup
- ✅ JWT-based authentication with access & refresh tokens
- ✅ Email verification system
- ✅ Password reset functionality
- ✅ Meal management (CRUD operations)
- ✅ Category management (CRUD operations)
- ✅ Image upload for meals
- ✅ Search functionality for meals
- ✅ Filter meals by category

### Database Features

- Auto-incrementing primary keys
- Automatic timestamp updates via triggers
- Data validation constraints (check constraints, unique constraints)
- Role-based user management
- Foreign key constraints with cascade delete

---

## 🚦 API Endpoints

### 🔐 Authentication Endpoints

#### 📝 Register
- **POST** `/auth/register`
- Registers a new user and sends a verification code to their email.

**Request Body:**
```json
{
  "username": "johndoe",
  "name": "John Doe",
  "email": "john@example.com",
  "password": "SecurePass123!",
  "phone": "+1234567890",
  "address": "123 Main St, City"
}
```

**Response (201 Created):**
```json
{
  "message": "Registration successful! Please check your email"
}
```

**Response (400 Bad Request):**
```json
{
  "error": "Invalid data"
}
```

**Response (409 Conflict):**
```json
{
  "error": "Email already in use"
}
```

**cURL Example:**
```bash
curl -X POST http://localhost:8000/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "name": "John Doe",
    "email": "john@example.com",
    "password": "SecurePass123!",
    "phone": "+1234567890",
    "address": "123 Main St, City"
  }'
```

---

#### ✅ Verify Email
- **POST** `/auth/verify-email`
- Verifies the user's email using the code sent.

**Request Body:**
```json
{
  "email": "john@example.com",
  "code": "123456"
}
```

**Response (200 OK):**
```json
{
  "message": "Account activated successfully"
}
```

**Response (400 Bad Request):**
```json
{
  "error": "Invalid code or expired or account already activated"
}
```

**cURL Example:**
```bash
curl -X POST http://localhost:8000/auth/verify-email \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "code": "123456"
  }'
```

---

#### 🔑 Login
- **POST** `/auth/login`
- Authenticates the user and returns access token. Refresh token is set as HTTP-only cookie.

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "SecurePass123!"
}
```

**Response (200 OK):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "role": "user"
}
```

**Response (401 Unauthorized):**
```json
{
  "error": "Invalid email or password"
}
```

**Response (403 Forbidden):**
```json
{
  "error": "Account not activated. Please verify your email"
}
```

**cURL Example:**
```bash
curl -X POST http://localhost:8000/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "SecurePass123!"
  }'
```

---

#### 🚪 Logout
- **POST** `/auth/logout`
- Logs out the user by invalidating the refresh token.

**Request Body:**
```json
{
  "refresh_token": "dGhpcyBpcyBhIHJlZnJlc2ggdG9rZW4..."
}
```

**Response (200 OK):**
```json
{
  "message": "Logged out successfully"
}
```

**Response (400 Bad Request):**
```json
{
  "error": "Invalid data"
}
```

**cURL Example:**
```bash
curl -X POST http://localhost:8000/auth/logout \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "your_refresh_token_here"
  }'
```

---

#### ♻️ Refresh Token
- **POST** `/auth/refresh`
- Returns a new access token if the refresh token cookie is valid and not expired.

**Request:**
- Requires `refresh_token` HTTP-only cookie (sent automatically by browser)

**Response (200 OK):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response (401 Unauthorized):**
```json
{
  "error": "Token not found"
}
```

**Response (401 Unauthorized):**
```json
{
  "error": "Invalid or expired token"
}
```

**cURL Example:**
```bash
curl -X POST http://localhost:8000/auth/refresh \
  -H "Cookie: refresh_token=your_refresh_token_here"
```

---

#### 🔒 Reset Password (Request Code)
- **POST** `/auth/reset-password`
- Sends a password reset code to the user's email.

**Request Body:**
```json
{
  "email": "john@example.com"
}
```

**Response (200 OK):**
```json
{
  "message": "Password reset code has been sent to your email if registered."
}
```

**cURL Example:**
```bash
curl -X POST http://localhost:8000/auth/reset-password \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com"
  }'
```

---

#### 🔍 Verify Reset Code
- **POST** `/auth/verify-reset-code`
- Verifies the reset code before allowing password update.

**Request Body:**
```json
{
  "email": "john@example.com",
  "code": "123456"
}
```

**Response (200 OK):**
```json
{
  "message": "Code is valid"
}
```

**Response (400 Bad Request):**
```json
{
  "error": "Invalid data. Please ensure all fields are filled correctly."
}
```

**Response (400 Bad Request):**
```json
{
  "error": "Invalid or expired code."
}
```

**cURL Example:**
```bash
curl -X POST http://localhost:8000/auth/verify-reset-code \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "code": "123456"
  }'
```

---

#### 🛠️ Update Password
- **POST** `/auth/update-password`
- Updates the user's password if the code is valid and not expired.

**Request Body:**
```json
{
  "email": "john@example.com",
  "code": "123456",
  "password": "NewPassword@123",
  "confirmPassword": "NewPassword@123"
}
```

**Response (200 OK):**
```json
{
  "message": "Password has been updated."
}
```

**Validation Responses:**

- If passwords do not match (400):
```json
{
  "error": "Passwords do not match."
}
```

- If password length is invalid (400):
```json
{
  "error": "Password must be between 8 and 20 characters."
}
```

- If password is not strong (400):
```json
{
  "error": "Password must contain uppercase, lowercase, number, and special character."
}
```

- If code is invalid or expired (400):
```json
{
  "error": "Invalid or expired code."
}
```

**cURL Example:**
```bash
curl -X POST http://localhost:8000/auth/update-password \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "code": "123456",
    "password": "NewPassword@123",
    "confirmPassword": "NewPassword@123"
  }'
```

---

### 🍽️ Meal Endpoints

#### ➕ Create Meal
- **POST** `/create/meal`
- Creates a new meal item with image upload.

**Request:**
- Content-Type: `multipart/form-data`

**Form Fields:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | Yes | Meal name |
| description | string | Yes | Meal description |
| price | float | Yes | Meal price |
| category_id | integer | Yes | Category ID |
| image | file | Yes | Meal image (max 10MB) |

**Response (201 Created):**
```json
{
  "id": 1,
  "name": "Margherita Pizza",
  "description": "Classic Italian pizza with tomato sauce and mozzarella",
  "price": 12.99,
  "image_url": "uploads/1234567890_pizza.jpg",
  "is_available": true,
  "category_name": "Pizza",
  "created_at": "2024-01-15T10:30:00Z"
}
```

**Response (400 Bad Request):**
```json
{
  "error": "Failed to parse form data"
}
```

**Response (400 Bad Request):**
```json
{
  "error": "Invalid price"
}
```

**Response (400 Bad Request):**
```json
{
  "error": "Invalid category_id"
}
```

**Response (400 Bad Request):**
```json
{
  "error": "Image is required"
}
```

**Response (400 Bad Request):**
```json
{
  "error": "Category does not exist"
}
```

**cURL Example:**
```bash
curl -X POST http://localhost:8000/create/meal \
  -F "name=Margherita Pizza" \
  -F "description=Classic Italian pizza" \
  -F "price=12.99" \
  -F "category_id=1" \
  -F "image=@/path/to/pizza.jpg"
```

---

#### 📋 Get All Meals
- **GET** `/meals`
- Retrieves all meals from the database.

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "name": "Margherita Pizza",
    "description": "Classic Italian pizza",
    "price": 12.99,
    "image_url": "uploads/1234567890_pizza.jpg",
    "is_available": true,
    "category_name": "Pizza",
    "created_at": "2024-01-15T10:30:00Z"
  },
  {
    "id": 2,
    "name": "Chicken Burger",
    "description": "Grilled chicken with fresh veggies",
    "price": 8.99,
    "image_url": "uploads/1234567891_burger.jpg",
    "is_available": true,
    "category_name": "Burgers",
    "created_at": "2024-01-15T11:00:00Z"
  }
]
```

**cURL Example:**
```bash
curl -X GET http://localhost:8000/meals
```

---

#### 🔍 Get Meal by ID
- **GET** `/meals/{id}`
- Retrieves a specific meal by its ID.

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| id | integer | Meal ID |

**Response (200 OK):**
```json
{
  "id": 1,
  "name": "Margherita Pizza",
  "description": "Classic Italian pizza with tomato sauce and mozzarella",
  "price": 12.99,
  "image_url": "uploads/1234567890_pizza.jpg",
  "is_available": true,
  "category_name": "Pizza",
  "created_at": "2024-01-15T10:30:00Z"
}
```

**Response (400 Bad Request):**
```json
{
  "error": "Invalid ID"
}
```

**Response (404 Not Found):**
```json
{
  "error": "Meal not found"
}
```

**cURL Example:**
```bash
curl -X GET http://localhost:8000/meals/1
```

---

#### 🗑️ Delete Meal
- **DELETE** `/delete/meal/{id}`
- Deletes a meal by its ID.

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| id | integer | Meal ID |

**Response (204 No Content):**
- Empty response body

**Response (400 Bad Request):**
```json
{
  "error": "Invalid ID"
}
```

**Response (500 Internal Server Error):**
```json
{
  "error": "Failed to delete meal"
}
```

**cURL Example:**
```bash
curl -X DELETE http://localhost:8000/delete/meal/1
```

---

#### ✏️ Update Meal
- **PUT** `/update/meal/{id}`
- Updates a meal's information (partial updates supported).

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| id | integer | Meal ID |

**Request Body:**
```json
{
  "name": "Updated Pizza Name",
  "description": "Updated description",
  "price": 14.99,
  "image_url": "uploads/new_image.jpg",
  "is_available": false,
  "category_id": 2
}
```

**Note:** All fields are optional. Only provided fields will be updated.

**Response (200 OK):**
```json
{
  "id": 1,
  "name": "Updated Pizza Name",
  "description": "Updated description",
  "price": 14.99,
  "image_url": "uploads/new_image.jpg",
  "is_available": false,
  "category_name": "Pizza",
  "created_at": "2024-01-15T10:30:00Z"
}
```

**Response (400 Bad Request):**
```json
{
  "error": "Invalid ID"
}
```

**Response (400 Bad Request):**
```json
{
  "error": "Invalid data"
}
```

**Response (400 Bad Request):**
```json
{
  "error": "No data provided for update"
}
```

**Response (500 Internal Server Error):**
```json
{
  "error": "Failed to update meal"
}
```

**cURL Example:**
```bash
curl -X PUT http://localhost:8000/update/meal/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Pizza Name",
    "price": 14.99,
    "is_available": false
  }'
```

---

#### 📂 Get Meals by Category
- **GET** `/meals/by-category/{category_id}`
- Retrieves all meals belonging to a specific category.

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| category_id | integer | Category ID |

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "name": "Margherita Pizza",
    "description": "Classic Italian pizza",
    "price": 12.99,
    "image_url": "uploads/1234567890_pizza.jpg",
    "is_available": true,
    "category_name": "Pizza",
    "created_at": "2024-01-15T10:30:00Z"
  },
  {
    "id": 3,
    "name": "Pepperoni Pizza",
    "description": "Spicy pepperoni with cheese",
    "price": 14.99,
    "image_url": "uploads/1234567892_pepperoni.jpg",
    "is_available": true,
    "category_name": "Pizza",
    "created_at": "2024-01-15T12:00:00Z"
  }
]
```

**Response (400 Bad Request):**
```json
{
  "error": "Invalid ID"
}
```

**cURL Example:**
```bash
curl -X GET http://localhost:8000/meals/by-category/1
```

---

#### 🔎 Search Meals
- **GET** `/meals/search?q={query}`
- Searches for meals by name or description.

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| q | string | Yes | Search query string |

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "name": "Margherita Pizza",
    "description": "Classic Italian pizza",
    "price": 12.99,
    "image_url": "uploads/1234567890_pizza.jpg",
    "is_available": true,
    "category_name": "Pizza",
    "created_at": "2024-01-15T10:30:00Z"
  }
]
```

**Response (400 Bad Request):**
```json
{
  "error": "Query parameter 'q' is required"
}
```

**cURL Example:**
```bash
curl -X GET "http://localhost:8000/meals/search?q=pizza"
```

---

### 📂 Category Endpoints

#### ➕ Create Category
- **POST** `/create/categories`
- Creates a new category.

**Request Body:**
```json
{
  "category_name": "Pizza"
}
```

**Response (200 OK):**
```json
{
  "id": 1,
  "category_name": "Pizza",
  "created_at": "2024-01-15T10:00:00Z"
}
```

**Response (400 Bad Request):**
```json
{
  "error": "Invalid data"
}
```

**Response (400 Bad Request):**
```json
{
  "error": "Category name is required"
}
```

**Response (409 Conflict):**
```json
{
  "error": "Category already exists"
}
```

**cURL Example:**
```bash
curl -X POST http://localhost:8000/create/categories \
  -H "Content-Type: application/json" \
  -d '{
    "category_name": "Pizza"
  }'
```

---

#### 📋 Get All Categories
- **GET** `/categories`
- Retrieves all categories.

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "category_name": "Pizza",
    "created_at": "2024-01-15T10:00:00Z"
  },
  {
    "id": 2,
    "category_name": "Burgers",
    "created_at": "2024-01-15T10:05:00Z"
  },
  {
    "id": 3,
    "category_name": "Drinks",
    "created_at": "2024-01-15T10:10:00Z"
  }
]
```

**cURL Example:**
```bash
curl -X GET http://localhost:8000/categories
```

---

#### 🔍 Get Category by ID
- **GET** `/categories/{id}`
- Retrieves a specific category by its ID.

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| id | integer | Category ID |

**Response (200 OK):**
```json
{
  "id": 1,
  "category_name": "Pizza",
  "created_at": "2024-01-15T10:00:00Z"
}
```

**Response (400 Bad Request):**
```json
{
  "error": "Invalid ID"
}
```

**Response (404 Not Found):**
```json
{
  "error": "Category not found"
}
```

**cURL Example:**
```bash
curl -X GET http://localhost:8000/categories/1
```

---

#### ✏️ Update Category
- **PUT** `/update/categories/{id}`
- Updates a category's name.

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| id | integer | Category ID |

**Request Body:**
```json
{
  "category_name": "Italian Pizza"
}
```

**Response (200 OK):**
```json
{
  "id": 1,
  "category_name": "Italian Pizza",
  "created_at": "2024-01-15T10:00:00Z"
}
```

**Response (400 Bad Request):**
```json
{
  "error": "Invalid ID"
}
```

**Response (400 Bad Request):**
```json
{
  "error": "Invalid data"
}
```

**Response (400 Bad Request):**
```json
{
  "error": "category_name is required"
}
```

**Response (404 Not Found):**
```json
{
  "error": "Category not found or error occurred"
}
```

**cURL Example:**
```bash
curl -X PUT http://localhost:8000/update/categories/1 \
  -H "Content-Type: application/json" \
  -d '{
    "category_name": "Italian Pizza"
  }'
```

---

#### 🗑️ Delete Category
- **DELETE** `/delete/categories/{id}`
- Deletes a category by its ID.

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| id | integer | Category ID |

**Response (204 No Content):**
- Empty response body

**Response (400 Bad Request):**
```json
{
  "error": "Invalid ID"
}
```

**Response (404 Not Found):**
```json
{
  "error": "Category not found"
}
```

**Response (500 Internal Server Error):**
```json
{
  "error": "Error occurred during deletion"
}
```

**cURL Example:**
```bash
curl -X DELETE http://localhost:8000/delete/categories/1
```

---

### 🖼️ Static Files

#### Serve Uploaded Images
- **GET** `/uploads/{filename}`
- Serves uploaded meal images.

**cURL Example:**
```bash
curl -X GET http://localhost:8000/uploads/1234567890_pizza.jpg
```

---

## Password Reset & Update Flow

### 1. Request Password Reset
- **POST** `/auth/reset-password`
- Request body:
  ```json
  { "email": "user@example.com" }
  ```
- Always returns a generic message:
  ```json
  { "message": "Password reset code has been sent to your email." }
  ```
- If the email exists, a reset code is sent to the email (valid for 24 hours). If not, same message is returned for security.

### 2. Verify Reset Code (Optional)
- **POST** `/auth/verify-reset-code`
- Request body:
  ```json
  { "email": "user@example.com", "code": "123456" }
  ```
- This step is optional but recommended for better UX to validate the code before showing the password update form.

### 3. Update Password
- **POST** `/auth/update-password`
- Request body:
  ```json
  {
    "email": "user@example.com",
    "code": "123456",
    "password": "NewPassword@123",
    "confirmPassword": "NewPassword@123"
  }
  ```
- Always returns a generic message:
  ```json
  { "message": "Password has been updated." }
  ```

### 4. Security & UX Notes
- The reset code is only used in `/auth/update-password` (not in `/auth/reset-password`).
- All responses are unified to prevent user enumeration and brute-force attacks.
- The code is deleted after successful password update and cannot be reused.
- The user interface should:
  1. Ask for email/username to request reset.
  2. Ask for the code (from email) and new password (with confirmation) to update password.
- No sensitive info is leaked in any response.

---

## Authentication Flow

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   Register  │────▶│ Verify Email │────▶│    Login    │
└─────────────┘     └──────────────┘     └──────┬──────┘
                                                │
                                                ▼
                                        ┌──────────────┐
                                        │  Access Token│◀────┐
                                        │  (15 min)    │     │
                                        └──────────────┘     │
                                                │            │
                                                ▼            │
                                        ┌──────────────┐     │
                                        │ Refresh Token│─────┘
                                        │  (7 days)    │
                                        └──────────────┘
```

### Token Management

1. **Access Token**: Short-lived (15 minutes), used for authenticated requests
2. **Refresh Token**: Long-lived (7 days), stored in HTTP-only cookie, used to get new access tokens
3. **Token Refresh**: When access token expires, use `/auth/refresh` with the refresh token cookie to get a new access token
4. **Logout**: Invalidates the refresh token, requiring re-authentication

---

## Data Models

### User Model
```go
type User struct {
    ID               int       `json:"id"`
    Name             string    `json:"name"`
    Username         string    `json:"username"`
    Email            string    `json:"email"`
    Password         string    `json:"password"`
    Phone            string    `json:"phone"`
    Address          string    `json:"address"`
    Role             string    `json:"role"`
    IsVerified       bool      `json:"is_verified"`
    VerificationCode string    `json:"verification_code"`
    CreateAt         time.Time `json:"created_at"`
}
```

### Category Model
```go
type Category struct {
    ID           int       `json:"id"`
    CategoryName string    `json:"category_name"`
    CreatedAt    time.Time `json:"created_at"`
}
```

### Meal Model
```go
type Meal struct {
    ID           int       `json:"id"`
    Name         *string   `json:"name,omitempty"`
    Description  *string   `json:"description,omitempty"`
    Price        *float64  `json:"price,omitempty"`
    ImageURL     *string   `json:"image_url,omitempty"`
    IsAvailable  *bool     `json:"is_available,omitempty"`
    CategoryID   *int      `json:"-"`
    CreatedAt    time.Time `json:"created_at"`
    CategoryName *string   `json:"category_name,omitempty"`
}
```

---

## Validation Rules

### User Registration
- **Username**: Unique, required
- **Email**: Valid email format, unique, required
- **Password**: 
  - 8-20 characters
  - At least one uppercase letter
  - At least one lowercase letter
  - At least one number
  - At least one special character
- **Phone**: Optional, unique if provided

### Meal Creation
- **Name**: Required
- **Description**: Required
- **Price**: Required, positive number
- **Category ID**: Required, must exist in categories table
- **Image**: Required, max 10MB

### Category Creation
- **Category Name**: Required, unique

---

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

### Image Upload Issues

**Problem:** `Failed to save image` or `Failed to copy image`

**Solutions:**
1. Ensure the `uploads` directory exists and is writable
2. Check disk space availability
3. Verify file size is under 10MB limit

### Email Sending Issues

**Problem:** Verification or reset emails not being sent

**Solutions:**
1. Check SMTP configuration in `.env` file
2. For Gmail, use an App Password instead of your regular password
3. Ensure less secure app access is enabled (for non-Gmail providers)
4. Check spam/junk folders in recipient email

---

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/new-feature`)
3. Commit your changes (`git commit -m 'Add new feature'`)
4. Push to the branch (`git push origin feature/new-feature`)
5. Open a Pull Request

---

<p align="center">
  <img src="https://em-content.zobj.net/source/microsoft-teams/363/bento-box_1f96a.png" width="60"/>
  <img src="https://em-content.zobj.net/source/microsoft-teams/363/takeout-box_1f961.png" width="60"/>
  <img src="https://em-content.zobj.net/source/microsoft-teams/363/pizza_1f355.png" width="60"/>
  <img src="https://em-content.zobj.net/source/microsoft-teams/363/cooking_1f373.png" width="60"/>
</p>

## License

This project is private and proprietary.

---

**Author:** Akram Fattah  
**Repository:** https://github.com/akram-fattah/food-delivery

# 🍔 Food Delivery API 🚀

A professional, robust, and secure RESTful API for a food delivery application built with **Go** and **PostgreSQL**. This project follows clean architecture principles and implements modern security standards.

## 🌟 Key Features

- **Security & Middleware**: Implemented custom middleware for Authentication, Authorization (RBAC), and CORS.
- **Authentication**: JWT-based system with Access and Refresh tokens.
- **Email System**: Integrated email verification and password reset functionality.
- **Role-Based Access Control (RBAC)**: Distinct permissions for `user`, `admin`, and `delivery` roles.
- **Meal & Category Management**: Full CRUD operations with image upload support.
- **Order Processing**: Real-time order creation and status tracking.
- **WhatsApp Integration**: Automated bot for notifications and interactions.

## 🛠️ Tech Stack

| Technology | Purpose |
| :--- | :--- |
| **Go (Golang)** | High-performance backend logic |
| **PostgreSQL** | Reliable relational data storage |
| **JWT** | Secure stateless authentication |
| **Bcrypt** | Industry-standard password hashing |
| **CORS** | Cross-Origin Resource Sharing for frontend integration |
| **SMTP** | Automated email communication |

## 🏗️ Project Structure

```text
food-delivery/
├── cmd/api/main.go          # 🚦 Application entry point & Route configuration
├── internal/
│   ├── middleware/          # 🛡️ Security & Request processing (CORS, Auth, Admin)
│   ├── handlers/            # 🎮 HTTP request handlers (Controllers)
│   ├── database/            # 🗄️ Data access layer (Repositories)
│   ├── models/              # 📦 Data structures & Schema definitions
│   ├── helper/              # 🧩 Utility functions & Security helpers
│   └── whatsapp/            # 🤖 WhatsApp bot integration
├── uploads/                 # 🖼️ Static assets & Meal images
└── delivery_db.sql          # 🗃️ Database schema & Initial setup
```

## 🛡️ Security & Middleware

The project implements a robust middleware layer to ensure security and cross-platform compatibility:

### 1. CORS (Cross-Origin Resource Sharing)
Configured to allow secure communication with frontend applications (React, Vue, Mobile apps).
- **Allowed Origins**: `*` (Configurable for production)
- **Allowed Methods**: `GET, POST, PUT, DELETE, OPTIONS, PATCH`
- **Allowed Headers**: `Content-Type, Authorization`

### 2. Authentication Middleware
Protects sensitive routes by validating JWT tokens in the `Authorization` header. It supports the `Bearer <token>` format and injects user identity into the request context.

### 3. Authorization Middleware (Admin Only)
A specialized layer that restricts access to administrative endpoints (e.g., managing categories, updating roles) only to users with the `admin` role.

## 🚦 API Endpoints

### 🔐 Authentication
| Method | Endpoint | Description | Access |
| :--- | :--- | :--- | :--- |
| POST | `/auth/register` | Create new account | Public |
| POST | `/auth/login` | Authenticate & get tokens | Public |
| POST | `/auth/verify-email` | Verify email with code | Public |
| POST | `/auth/refresh` | Get new access token | Public |
| POST | `/auth/logout` | Invalidate session | **User** |
| PUT | `/auth/update-password` | Change password | **User** |

### 🍽️ Meal & Category Management
| Method | Endpoint | Description | Access |
| :--- | :--- | :--- | :--- |
| GET | `/categories` | List all categories | Public |
| POST | `/create/categories` | Create new category | **Admin** |
| GET | `/meals` | List all meals | Public |
| POST | `/create/meal` | Add new meal | **Admin** |
| DELETE | `/delete/meal/:id` | Remove a meal | **Admin** |

### 📦 Orders & Profile
| Method | Endpoint | Description | Access |
| :--- | :--- | :--- | :--- |
| GET | `/profile` | Get current user info | **User** |
| PATCH | `/profile/update` | Update profile details | **User** |
| POST | `/orders/create` | Place a new order | **User** |
| GET | `/orders` | View order history | **User** |
| PATCH | `/orders/update-status/:id` | Update order status | **Admin** |

## 🚀 Getting Started

### Prerequisites
- Go 1.26+
- PostgreSQL 12+
- SMTP Server (e.g., Gmail)

### Installation
1. **Clone & Install**:
   ```bash
   git clone https://github.com/akram-fattah/food-delivery.git
   cd food-delivery
   go mod tidy
   ```

2. **Database Setup**:
   ```bash
   psql -U postgres -c "CREATE DATABASE delivery_db;"
   psql -U postgres -d delivery_db -f delivery_db.sql
   ```

3. **Environment Config**:
   Create a `.env` file and fill in your credentials (DB, JWT, SMTP).

4. **Run**:
   ```bash
   go run cmd/api/main.go
   ```

---
*Developed with ❤️ by [Akram Fattah](https://github.com/akram-fattah)*

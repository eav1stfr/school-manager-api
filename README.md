# 📚 School Manager API

A secure and efficient RESTful backend server for managing **students**, **teachers**, and **executives** in an academic environment. This Go-based API supports full CRUD operations, secure authentication, and various middleware layers for production-readiness.

---

## ✨ Features

- ✅ **CRUD** operations for:
  - Students
  - Teachers (with endpoint to get number of their students)
  - Executives (execs)
- 🔒 **Authentication & Authorization**:
  - JWT login/logout system
  - Password hashing using `argon2` with salt
  - Reset password functionality
- 🛡️ **Middlewares**:
  - CORS
  - Rate Limiting
  - HTTP Parameter Pollution (HPP) protection
  - Response time measurement
  - Compression
  - Secure headers
- 🗄️ **Database**:
  - MariaDB with `sqlx` for SQL queries
  - Three main tables: `students`, `teachers`, `execs`
- ⚙️ TLS support with HTTPS
- 🧪 API tested with Postman

---

## 🏗️ Project Structure

```
School_Manager_Project/
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   ├── middlewares/
│   │   └── router/
│   ├── models/
│   ├── sqlconnect/
├── server/
└── utils/
```

---

## 🧪 Tech Stack

- **Language**: Go (Golang)
- **Database**: MariaDB
- **ORM/SQL**: [`github.com/jmoiron/sqlx`](https://github.com/jmoiron/sqlx)
- **Authentication**: JWT, Argon2 password hashing
- **Middleware**: Custom implementation for:
  - CORS
  - Rate Limiting
  - HPP
  - Compression
  - Security Headers
  - Response Timing
- **HTTPS**: TLS v1.2+
- **Environment Config**: `.env` + [`github.com/joho/godotenv`](https://github.com/joho/godotenv)

---

## 🔐 Security

- Passwords are hashed using **Argon2** with securely generated salts:
  
  ```go
  salt := make([]byte, 16)
  _, err := rand.Read(salt)
  hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
  ```

- JWTs are used for access control for executive users, with protected routes and middleware to exclude public login endpoints.

---

## 🧭 API Overview

| Resource | Method | Endpoint | Description |
|----------|--------|----------|-------------|
| Students | GET | `/students` | Get all students |
| Students | POST | `/students` | Create student |
| Students | PUT | `/students/:id` | Update student |
| Students | DELETE | `/students/:id` | Delete student |
| Teachers | GET | `/teachers` | Get all teachers |
| Teachers | GET | `/teachers/:id/students/count` | Get student count of a teacher |
| Execs | POST | `/execs/login` | Login (JWT) |
| Execs | POST | `/execs/logout` | Logout |
| Execs | PATCH | `/execs/reset-password` | Reset password |

---

## 🧪 Testing

All endpoints were tested using **Postman** collections during development.

---

## 🚀 Running the Project

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/school-manager-api.git
   ```

2. Navigate into the project:
   ```bash
   cd school-manager-api
   ```

3. Fill in your `.env` file and TLS certificates (`cert.pem`, `key.pem`)

4. Run the server:
   ```bash
   go run server/server.go
   ```

---


# 📚 School Manager API

A RESTful backend server built in Go for managing a school system. It supports secure and efficient CRUD operations for **Students**, **Teachers**, and **Executives (Execs)**. The system is equipped with secure authentication, middleware layers, and follows clean architecture principles.

---

## 🚀 Features

- Full **CRUD** operations for:
  - 👨‍🏫 Teachers
  - 🧑‍🎓 Students
  - 🧑‍💼 Execs (admin users)
- 🔐 **Authentication & Authorization** for Execs using **JWT**
  - `POST /execs/login`
  - `POST /execs/logout`
  - `POST /execs/reset-password`
- 👥 **Teacher-to-student relationship**:
  - `GET /teachers/{id}/students` — returns the number of students assigned to a teacher
- 🧂 **Secure password hashing** using **argon2** with salt & base64 encoding
- 📄 .env-based configuration
- 🔗 Clean project structure with modular packages

---

## ⚙️ Tech Stack

| Layer             | Technology                          |
|------------------|--------------------------------------|
| Language          | Go (Golang)                         |
| HTTP Server       | net/http                            |
| Database          | MariaDB + sqlx                      |
| Auth              | JWT (`github.com/golang-jwt/jwt`)   |
| Hashing           | Argon2 + base64                     |
| Environment       | `github.com/joho/godotenv`          |
| API Testing       | Postman                             |

---

## 🧰 Middlewares

- 🛡 **CORS** — Cross-Origin Resource Sharing
- 🔁 **Rate Limiting** — Prevents abuse (based on IP address)
- 📦 **Compression** — Response compression
- 🧨 **HPP (HTTP Parameter Pollution)** protection
- 🔑 **JWT Authorization** — Role-based route protection

---

## 🔐 Password Hashing

Passwords are hashed using Argon2 with a random salt. Here's the logic:

```go
func Hash(password string) (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", ErrorGeneratingSaltForHashing
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	hashBase64 := base64.StdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("%s.%s", saltBase64, hashBase64)
	return encodedHash, nil
}

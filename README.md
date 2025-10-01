# Social Media Backend API

A simple **Social Media Backend API** built with **Go** and **Gin framework**.  
It supports user authentication, posts, likes, comments, follows, and notifications.

## Design System 
<img width="2644" height="860" alt="image" src="https://github.com/user-attachments/assets/444de218-52f1-4c85-a77f-8c2b2bbd758d" />

#### Functional Requirements

  - User Account Managements
  - User Interactions (Post, Feed, Comment, and Social Relationships)
  - User Discovery & Communication (Notification & etc)

#### Non-Functional Requirements

1. **Performance**
    - Server response
    - Optimize Query
2. **Scalability**
    - Handle +1k user
    - Expand to 1k++
3. **Reliability**
    - Consistency data (Post, Like, Comment, etc)
    - Error Handling
4. **Availability**
    - Server Availablity up > 95%
    - Server downtime
  
## 🔧 Tech Stack

- [Go](https://go.dev/dl/) - Backend programming language
- [Gin Gonic](https://gin-gonic.com/) - HTTP web framework
- [PostgreSQL](https://www.postgresql.org/download/) - Primary database
- [Redis](https://redis.io/docs/latest/operate/oss_and_stack/install/) - Caching and session management
- [JWT](https://github.com/golang-jwt/jwt) - Authentication tokens
- [Swagger](https://swagger.io/) + [Swaggo](https://github.com/swaggo/swag) - API documentation
- [Docker](https://docs.docker.com/engine/install/) - Containerization

## 🗝️ Environment Variables

```bash
# Database Configuration
DBUSER=<your_database_user>
DBPASS=<your_database_password>
DBNAME=<your_database_name>
DBHOST=<your_database_host>
DBPORT=<your_database_port>

# JWT Configuration
JWT_SECRET=<your_secret_jwt>

# Redis Configuration
RDB_HOST=<your_redis_host>
RDB_PORT=<your_redis_port>

# Server Configuration
PORT=8080
```

## ⚙️ Installation

1. Clone the repository

```sh
$ git clone <repository-url>
```

2. Navigate to project directory

```sh
$ cd social-media-backend
```

3. Install dependencies

```sh
$ go mod tidy
```

4. Set up environment variables as described above

5. Run database migrations

```sh
$ make migrate-up
```

6. Start the application

```sh
$ go run ./cmd/main.go
```

## 📚 API Documentation

The API is organized into four main resource groups: **Authentication**, **Feed Management**, **Post Management**, and **User Management**.

### Authentication Endpoints

| Method | Endpoint         | Description       | Authentication |
| ------ | ---------------- | ----------------- | -------------- |
| POST   | `/auth/login`    | User login        | ❌             |
| POST   | `/auth/register` | User registration | ❌             |

### Feed Management

| Method | Endpoint | Description                  | Authentication  |
| ------ | -------- | ---------------------------- | --------------- |
| GET    | `/feed`  | Get user's personalized feed | ✅ Bearer Token |

### Post Management

| Method | Endpoint             | Description           | Authentication  |
| ------ | -------------------- | --------------------- | --------------- |
| POST   | `/post`              | Create a new post     | ✅ Bearer Token |
| POST   | `/post/{id}/comment` | Add comment to a post | ✅ Bearer Token |
| POST   | `/post/{id}/like`    | Like a post           | ✅ Bearer Token |
| DELETE | `/post/{id}/unlike`  | Unlike a post         | ✅ Bearer Token |

### User Management

| Method | Endpoint                   | Description               | Authentication  |
| ------ | -------------------------- | ------------------------- | --------------- |
| GET    | `/user/`                   | Get all users             | ❌ Bearer Token |
| GET    | `/user/notifications`      | Get user notifications    | ✅ Bearer Token |
| PATCH  | `/user/notifications/{id}` | Mark notification as read | ✅ Bearer Token |
| POST   | `/user/{id}/follow`        | Follow a user             | ✅ Bearer Token |
| DELETE | `/user/{id}/unfollow`      | Unfollow a user           | ✅ Bearer Token |

## 🐳 Docker Support

**Build Docker Image:**

```bash
$ docker build -t social-media-backend .
```

## 🔧 Development

**Run Tests:**

```bash
$ go test ./cmd/main.go
```

**Generate Swagger Docs:**

```bash
$ swag init -g ./cmd/main.go
```

## 📄 License

MIT License - see LICENSE file for details.

# Content Management System API

_Go Udacity project #1_

A backend API for a Content Management System (CMS), built in Go with [Gin](https://github.com/gin-gonic/gin) and [GORM](https://gorm.io/). The API manages **pages**, **posts**, and **media**, simulating real-world CMS backend tasks  CRUD operations, database migrations, relational data modeling, environment-specific configuration, and automated testing.

## Features

- **CRUD operations** across three resources: pages, posts, and media
- **Relational data modeling** — posts support a many-to-many relationship with media items
- **Search-style filtering** — filter pages and posts by title (partial match) or author; filter media by URL (partial match) or type
- **Field validation** — each model validates required fields and length constraints before being persisted
- **Environment-aware configuration** — auto-migration and Gin's release mode are toggled based on the `ENV` variable
- **Swagger/OpenAPI documentation** — API metadata defined via code comments
- **Unit and integration tests** — including a mocked database setup for isolated unit testing

## Tech Stack

| Component        | Technology                                                |
|-------------------|------------------------------------------------------------|
| Language          | Go 1.22 (toolchain 1.23.1)                                  |
| Web framework     | [Gin](https://github.com/gin-gonic/gin)                     |
| ORM               | [GORM](https://gorm.io/) with the PostgreSQL driver          |
| Database          | PostgreSQL                                                   |
| Config            | [godotenv](https://github.com/joho/godotenv) (autoloaded `.env`) |
| Testing           | Go's built-in `testing` package + [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock) |

## Prerequisites

- **[Go](https://go.dev/dl/)** 1.22 or later — verify with:
  ```bash
  go version
  ```
- **[PostgreSQL](https://www.postgresql.org/download/)** 13+ running locally, in Docker, or hosted remotely
- **`make`** (optional, but used by the provided test commands)

## Installation

**1. Clone the repository**

```bash
git clone <your-repo-url>
cd cms-backend
```

**2. Install Go dependencies**

```bash
go mod download
```

**3. Set up PostgreSQL**

Create a database for development:

```bash
createdb cms_dev
```

**4. Configure environment variables**

Create a `.env` file in the project root:

```bash
touch .env
```

Populate it with your database connection details


> `.env` is loaded automatically at startup (via `godotenv/autoload`) — no explicit `Load()` call is needed in code. If any value is missing, the server will fail to connect with a database connection error.

> `ENV` controls two things: when set to `development` (the default if unset), the server runs `AutoMigrate` on startup to create/update tables. When set to `production`, auto-migration is skipped and Gin runs in release mode (less verbose logging, no debug warnings).

## Running the Server

Start the API with:

```bash
go run main.go
```

On startup, the server will:

1. Load environment variables from `.env`
2. Connect to PostgreSQL
3. Run `AutoMigrate` for `Page`, `Post`, and `Media` models — **only if `ENV` is `development` or unset**
4. Start listening on `http://localhost:8080`

## Usage

All endpoints are served under the base path `/api/v1`.

### Pages — `/api/v1/pages`

| Method | Endpoint                | Description                              |
|--------|--------------------------|--------------------------------------------|
| GET    | `/api/v1/pages`          | List pages (filter by `title`, `author`)   |
| GET    | `/api/v1/pages/:id`      | Retrieve a single page by ID                |
| POST   | `/api/v1/pages`          | Create a new page                           |
| PUT    | `/api/v1/pages/:id`      | Update an existing page                     |
| DELETE | `/api/v1/pages/:id`      | Delete a page                               |

**Create a page:**

```bash
curl -X POST "http://localhost:8080/api/v1/pages" \
  -H "Content-Type: application/json" \
  -d '{"title": "About Us", "content": "We build great software."}'
```

`title` and `content` are required. `title` cannot exceed 255 characters, and `content` cannot be empty.

### Posts — `/api/v1/posts`

| Method | Endpoint                | Description                              |
|--------|--------------------------|--------------------------------------------|
| GET    | `/api/v1/posts`          | List posts (filter by `title`, `author`), includes associated media |
| GET    | `/api/v1/posts/:id`      | Retrieve a single post by ID                |
| POST   | `/api/v1/posts`          | Create a new post                           |
| PUT    | `/api/v1/posts/:id`      | Update an existing post                     |
| DELETE | `/api/v1/posts/:id`      | Delete a post                               |

**Create a post:**

```bash
curl -X POST "http://localhost:8080/api/v1/posts" \
  -H "Content-Type: application/json" \
  -d '{"title": "Launch Day", "content": "We just shipped v1.", "author": "Jane Doe"}'
```

`title` and `content` are required; `author` is optional. Each post can be associated with multiple `Media` items via a many-to-many relationship.

### Media — `/api/v1/media`

| Method | Endpoint                | Description                              |
|--------|--------------------------|--------------------------------------------|
| GET    | `/api/v1/media`          | List media (filter by `url`, `type`)        |
| GET    | `/api/v1/media/:id`      | Retrieve a single media item by ID          |
| POST   | `/api/v1/media`          | Create a new media item                     |
| DELETE | `/api/v1/media/:id`      | Delete a media item                         |

> Note: media does not currently support an update (`PUT`) endpoint.

**Create a media item:**

```bash
curl -X POST "http://localhost:8080/api/v1/media" \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com/image.png", "type": "image"}'
```

`url` is required (max 255 characters); `type` is required (max 50 characters).

### Error Responses

Errors are returned as JSON with a consistent shape:

```json
{
  "code": 404,
  "message": "Post could not be found"
}
```

## Testing

This project includes both unit tests (using a mocked database via `go-sqlmock`) and integration tests (against a real PostgreSQL database).

**Run all tests:**

```bash
make test
```

**Run unit tests only:**

```bash
make test-unit
```

**Run integration tests only:**

```bash
make test-integration
```

**Set up a test database and run integration tests in one step:**

```bash
make test-integration-full
```

This drops and recreates a `cms_test` database (using `postgres`/`postgres` credentials by default — adjust the `create-test-db` target in the `makefile` if your local setup differs), then runs the integration suite against it.

## License

This project is licensed under the MIT License — see the [LICENSE](./LICENSE) file for details.
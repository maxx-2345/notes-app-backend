# Notes App — Backend Plan

> **Stack:** Go (Gin) · PostgreSQL · GORM  
> **Repo:** `notes-backend/` — standalone, no frontend here  
> **Rule:** Complete and test each step before moving to the next.

---

## Folder Structure

```
notes-backend/
├── main.go
├── go.mod
├── go.sum
├── .env
├── .gitignore
├── config/
│   └── db.go
├── models/
│   └── note.go
├── handlers/
│   └── note_handler.go
└── routes/
    └── routes.go
```

> `middleware/` folder will be added in Phase 2 when JWT auth is introduced.

---

## Phase 1 — Project Setup

### Step 1: Create the project

```bash
mkdir notes-backend && cd notes-backend
go mod init notes-backend
```

---

### Step 2: Install dependencies

```bash
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/postgres
go get github.com/joho/godotenv
go get github.com/gin-contrib/cors
```

---

### Step 3: Setup PostgreSQL with Docker

Create `docker-compose.yml` in `notes-backend/`:

```yaml
version: '3.8'
services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_USER: admin
      POSTGRES_PASSWORD: secret
      POSTGRES_DB: notes_db
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data
volumes:
  pgdata:
```

Start the database:

```bash
docker compose up -d
```

Verify it is running:

```bash
docker ps
```

---

### Step 4: Create .env file

Create `.env` in `notes-backend/`:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=secret
DB_NAME=notes_db
```

Create `.gitignore`:

```
.env
```

---

## Phase 2 — Database Connection

### Step 5: Write config/db.go

This file is responsible for:

1. Loading `.env` using `godotenv`
2. Reading DB credentials from environment variables
3. Building the DSN (connection string)
4. Opening a GORM connection with `gorm.Open(postgres.Open(dsn))`
5. Storing the DB instance in a package-level variable so handlers can use it
6. Calling `db.AutoMigrate()` to auto-create tables — you will add models here as you create them

DSN format:

```
host=localhost user=admin password=secret dbname=notes_db port=5432 sslmode=disable
```

---

## Phase 3 — Note Model

### Step 6: Write models/note.go

Define the `Note` struct:

```go
type Note struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

Then go back to `config/db.go` and add `&models.Note{}` to `AutoMigrate`.

> Later when auth is added you will add `UserID uint` here and a `User` model.

---

## Phase 4 — Handlers (Business Logic)

### Step 7: Write handlers/note_handler.go

Write one function per API operation. Each function receives `*gin.Context` and uses the GORM DB instance from `config`.

| Function | What it does |
|---|---|
| `GetAllNotes` | Query all notes → return JSON array |
| `GetNoteByID` | Get `:id` from URL param → query → return JSON |
| `CreateNote` | Bind JSON body → insert into DB → return created note |
| `UpdateNote` | Get `:id` → bind JSON body → update in DB → return updated note |
| `DeleteNote` | Get `:id` → delete from DB → return success message |

**Response rules:**
- Always return JSON, never plain text
- On success: return appropriate status (`200`, `201`)
- On not found: return `404` with `{"error": "note not found"}`
- On bad input: return `400` with `{"error": "invalid input"}`
- On DB error: return `500` with `{"error": "internal server error"}`

---

## Phase 5 — Routes

### Step 8: Write routes/routes.go

Create a `SetupRoutes` function that takes a `*gin.Engine` and registers all routes under the `/api` prefix:

```
GET    /api/notes        → GetAllNotes
GET    /api/notes/:id    → GetNoteByID
POST   /api/notes        → CreateNote
PUT    /api/notes/:id    → UpdateNote
DELETE /api/notes/:id    → DeleteNote
```

---

## Phase 6 — Entry Point

### Step 9: Write main.go

`main.go` should do exactly these things in order:

1. Load `.env`
2. Call `config.ConnectDB()` to connect and run migrations
3. Create Gin engine with `gin.Default()`
4. Add CORS middleware (allow all origins for now — you will restrict this later)
5. Call `routes.SetupRoutes(r)`
6. Run on port `8080`

---

## Phase 7 — Test Everything ✅

### Step 10: Test all 5 endpoints

Run the server:

```bash
go run main.go
```

Test with curl one by one:

```bash
# 1. Create a note
curl -X POST http://localhost:8080/api/notes \
  -H "Content-Type: application/json" \
  -d '{"title":"My First Note","content":"Hello World"}'

# 2. Get all notes
curl http://localhost:8080/api/notes

# 3. Get one note by ID
curl http://localhost:8080/api/notes/1

# 4. Update a note
curl -X PUT http://localhost:8080/api/notes/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated Title","content":"Updated content"}'

# 5. Delete a note
curl -X DELETE http://localhost:8080/api/notes/1

# 6. Confirm it is gone
curl http://localhost:8080/api/notes/1
```

All 6 checks must pass before the backend is considered done.

---

## What Comes After This

Once this backend is fully working you will move to these in order:

| Next Phase | What |
|---|---|
| Frontend | Separate `notes-frontend/` Flutter project |
| Auth | Register, Login, JWT token generation |
| User-based notes | Each user sees only their own notes |

---

## Quick Reference

| Thing | Value |
|---|---|
| Server port | `8080` |
| DB name | `notes_db` |
| DB user | `admin` |
| DB password | `secret` |
| DB port | `5432` |
| API base URL | `http://localhost:8080/api` |

---

## Build Checklist

```
[ ] Step 1  — Project created, go.mod initialized
[ ] Step 2  — All dependencies installed
[ ] Step 3  — Docker + PostgreSQL running
[ ] Step 4  — .env created, .gitignore set
[ ] Step 5  — config/db.go written, DB connects
[ ] Step 6  — models/note.go written, table auto-created in DB
[ ] Step 7  — All 5 handlers written
[ ] Step 8  — Routes registered
[ ] Step 9  — main.go wired up, server starts
[ ] Step 10 — All 5 endpoints tested with curl ✅
```
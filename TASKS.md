# TASKS.md
## Kamus Banjar API — Community Edition Implementation Tasks

Status legend: `[ ]` todo · `[~]` in progress · `[x]` done

---

## Phase 1 — Auth Foundation

### TASK-001 · Add dependencies

Install required Go packages:

```bash
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto
go get github.com/google/uuid
go get github.com/gofiber/fiber/v2/middleware/limiter
```

Update `go.mod` and `go.sum`.

- [ ] TASK-001

---

### TASK-002 · Database migration — `users` table

Create file `database/migrations/002_users.sql`:

```sql
CREATE TABLE IF NOT EXISTS users (
    id          CHAR(36)     NOT NULL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    email       VARCHAR(255) NOT NULL UNIQUE,
    password    VARCHAR(255) NOT NULL,
    role        ENUM('admin','user') NOT NULL DEFAULT 'user',
    is_active   TINYINT(1)   NOT NULL DEFAULT 1,
    created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

Also create `database/migrations/003_refresh_tokens.sql`:

```sql
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id          CHAR(36)     NOT NULL PRIMARY KEY,
    user_id     CHAR(36)     NOT NULL,
    token_hash  VARCHAR(255) NOT NULL UNIQUE,
    expires_at  DATETIME     NOT NULL,
    created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

- [ ] TASK-002

---

### TASK-003 · Admin seed script

Create `database/seeds/003_admin_seed.sql` — placeholder only (admin is seeded at runtime from env vars, not from SQL). The seed logic runs in Go on startup.

In `main.go` `setup()`, after DB connection, call `seedAdmin(db)`:

```go
func seedAdmin(db *sql.DB) {
    email := os.Getenv("ADMIN_EMAIL")
    password := os.Getenv("ADMIN_PASSWORD")
    if email == "" || password == "" {
        return
    }
    // Check if admin already exists
    var count int
    db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
    if count > 0 {
        return
    }
    hash, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
    id := uuid.New().String()
    db.Exec("INSERT INTO users (id, name, email, password, role) VALUES (?, 'Admin', ?, ?, 'admin')",
        id, email, string(hash))
    log.Printf("Admin seeded: %s", email)
}
```

- [ ] TASK-003

---

### TASK-004 · User domain — model, repository interface, service

Create `domain/user/` package. Files:

**`domain/user/model.go`**
```go
type User struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Role      string    `json:"role"`
    IsActive  bool      `json:"is_active"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type RegisterRequest struct {
    Name     string `json:"name"  validate:"required,min=2,max=100"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
    Email    string `json:"email"    validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

type UpdateProfileRequest struct {
    Name        string `json:"name"         validate:"omitempty,min=2,max=100"`
    OldPassword string `json:"old_password" validate:"omitempty"`
    NewPassword string `json:"new_password" validate:"omitempty,min=8"`
}

type TokenPair struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
}
```

**`domain/user/repository.go`** — interface:
```go
type Repository interface {
    Create(u User, passwordHash string) error
    FindByEmail(email string) (User, string, error) // returns user + password_hash
    FindByID(id string) (User, error)
    Update(id string, name string) error
    UpdatePassword(id string, hash string) error
    SaveRefreshToken(id, tokenHash string, expiresAt time.Time) error
    RevokeRefreshToken(tokenHash string) error
    FindRefreshToken(tokenHash string) (userID string, err error)
    // admin
    ListUsers(page, limit int) ([]User, int, error)
    SetActive(id string, active bool) error
    SetRole(id string, role string) error
}
```

**`domain/user/repository_mysql.go`** — implement all interface methods using `*sql.DB`.

**`domain/user/service.go`** — interface + implementation:
```go
type Service interface {
    Register(req RegisterRequest) (User, error)
    Login(req LoginRequest) (TokenPair, error)
    RefreshTokens(refreshToken string) (TokenPair, error)
    Logout(refreshToken string) error
    Me(userID string) (User, error)
    UpdateProfile(userID string, req UpdateProfileRequest) (User, error)
    // admin
    ListUsers(page, limit int) ([]User, int, error)
    SetActive(userID string, active bool) error
    Promote(userID string) error
}
```

Service implementation notes:
- `Register`: hash password with `bcrypt.GenerateFromPassword([]byte(pw), 12)`, check email uniqueness (return `409` sentinel error if duplicate)
- `Login`: fetch user by email, `bcrypt.CompareHashAndPassword`, check `is_active`, generate JWT pair
- `RefreshTokens`: SHA-256 hash incoming token, look up in DB, validate expiry, generate new pair, delete old token hash, save new
- `Logout`: SHA-256 hash token, delete from DB
- Token generation helper: `generateTokenPair(userID, role string) (TokenPair, error)` — access token signed with `JWT_SECRET`, claims: `sub`, `role`, `exp`; refresh token is a random UUID stored hashed

- [ ] TASK-004

---

### TASK-005 · JWT middleware + role middleware

Create `middleware/auth.go`:

```go
func AuthMiddleware(secret string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        header := c.Get("Authorization")
        if !strings.HasPrefix(header, "Bearer ") {
            return fiber.ErrUnauthorized
        }
        tokenStr := strings.TrimPrefix(header, "Bearer ")
        token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
            if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fiber.ErrUnauthorized
            }
            return []byte(secret), nil
        })
        if err != nil || !token.Valid {
            return fiber.ErrUnauthorized
        }
        claims := token.Claims.(jwt.MapClaims)
        c.Locals("userID", claims["sub"])
        c.Locals("role", claims["role"])
        return c.Next()
    }
}

func RoleMiddleware(role string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        if c.Locals("role") != role {
            return fiber.ErrForbidden
        }
        return c.Next()
    }
}
```

Create `middleware/rate_limit.go` — wrap Fiber's built-in `limiter` middleware:

```go
func AuthRateLimit() fiber.Handler {
    return limiter.New(limiter.Config{
        Max:        10,
        Expiration: 1 * time.Minute,
        KeyGenerator: func(c *fiber.Ctx) string {
            return c.IP()
        },
        LimitReached: func(c *fiber.Ctx) error {
            return c.Status(429).JSON(fiber.Map{"error": "too many requests"})
        },
    })
}
```

- [ ] TASK-005

---

### TASK-006 · Auth HTTP handler + route registration

Create `domain/user/handler.go`:

```
POST /api/v1/auth/register  → handler.Register
POST /api/v1/auth/login     → handler.Login
POST /api/v1/auth/refresh   → handler.Refresh
POST /api/v1/auth/logout    → handler.Logout  (requires AuthMiddleware)
GET  /api/v1/auth/me        → handler.Me      (requires AuthMiddleware)
PUT  /api/v1/auth/me        → handler.UpdateProfile (requires AuthMiddleware)
```

Handler validation: parse body, validate required fields, return `400` with field errors on invalid input. Return `409` when email already registered.

In `main.go` `setup()`, wire new routes:
```go
userRepository := user.NewRepository(db)
userService := user.NewService(userRepository, jwtSecret, accessTTL, refreshTTL)
userHandler := user.NewHandler(userService)

auth := api.Group("/v1/auth")
auth.Post("/register", middleware.AuthRateLimit(), userHandler.Register)
auth.Post("/login",    middleware.AuthRateLimit(), userHandler.Login)
auth.Post("/refresh",  userHandler.Refresh)
auth.Post("/logout",   middleware.AuthMiddleware(jwtSecret), userHandler.Logout)
auth.Get("/me",        middleware.AuthMiddleware(jwtSecret), userHandler.Me)
auth.Put("/me",        middleware.AuthMiddleware(jwtSecret), userHandler.UpdateProfile)
```

Read `JWT_SECRET` from env in `main.go`; if empty, `log.Fatal("JWT_SECRET is required")`.

- [ ] TASK-006

---

### TASK-007 · Config — new env vars

Update `config_default.go` (or wherever config is read) to load and expose:

```go
type AppConfig struct {
    Port           string
    JWTSecret      string
    JWTAccessTTL   time.Duration  // default 15m
    JWTRefreshTTL  time.Duration  // default 168h
    AdminEmail     string
    AdminPassword  string
    MySQLDSN       string
    SourcePath     string
}
```

`JWT_SECRET` missing → `log.Fatal`. All others have defaults.

- [ ] TASK-007

---

### TASK-008 · Tests — auth service unit tests

Create `domain/user/service_test.go`. Use mock repository (`MockRepository`). Cover:
- Successful register
- Duplicate email → `409` sentinel
- Password too short → validation error
- Successful login
- Wrong password → `401` sentinel
- Inactive user login → `403` sentinel
- Token refresh with valid token
- Token refresh with expired/revoked token → error

- [ ] TASK-008

---

## Phase 2 — Word Source Tracking

### TASK-009 · Database migration — extend `words` table

Create `database/migrations/004_word_source.sql`:

```sql
ALTER TABLE words
    ADD COLUMN source        ENUM('official','community') NOT NULL DEFAULT 'official' AFTER data,
    ADD COLUMN status        ENUM('active','pending','rejected') NOT NULL DEFAULT 'active' AFTER source,
    ADD COLUMN contributor_id CHAR(36) NULL AFTER status,
    ADD COLUMN approved_by   CHAR(36) NULL AFTER contributor_id,
    ADD COLUMN approved_at   DATETIME NULL AFTER approved_by,
    ADD COLUMN created_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP AFTER approved_at,
    ADD COLUMN updated_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP AFTER created_at,
    ADD CONSTRAINT fk_words_contributor FOREIGN KEY (contributor_id) REFERENCES users(id) ON DELETE SET NULL,
    ADD CONSTRAINT fk_words_approver    FOREIGN KEY (approved_by)    REFERENCES users(id) ON DELETE SET NULL;

-- Backfill: existing rows are official + active
UPDATE words SET source = 'official', status = 'active' WHERE source IS NULL OR source = '';
```

- [ ] TASK-009

---

### TASK-010 · Update `Word` model

In `domain/dictionary/model.go`, add fields to `Word`:

```go
type Word struct {
    // existing fields...
    Source        string     `json:"source,omitempty"`
    ContributorID *string    `json:"contributor_id,omitempty"`
    ApprovedBy    *string    `json:"approved_by,omitempty"`
    ApprovedAt    *time.Time `json:"approved_at,omitempty"`
    CreatedAt     *time.Time `json:"created_at,omitempty"`
    UpdatedAt     *time.Time `json:"updated_at,omitempty"`
}
```

`status` is intentionally NOT in the public `Word` JSON response. It is an internal field — keep it as a separate struct `WordRecord` used inside repository if needed.

- [ ] TASK-010

---

### TASK-011 · Update MySQL repository — filter `status: active`

In `domain/dictionary/repository_mysql.go`, update all read queries to add `AND w.status = 'active'`:

- `GetAlphabets()` — `JOIN words w ON l.letter = w.letter AND w.status = 'active'`
- `GetWordsByAlphabet()` — add `AND w.status = 'active'`
- `GetWord()` — add `AND w.status = 'active'`
- `Search()` — add `AND w.status = 'active'` in both exact and levenshtein queries

Public API now silently ignores pending/rejected words.

- [ ] TASK-011

---

### TASK-012 · Add `source` to public word detail response

Update `GetWord()` in `repository_mysql.go` to also scan the `source` field alongside `data`. Unmarshal `data` JSON into `Word`, then set `Word.Source = source`.

The embed/FS/JSON repositories return words without `source` field (acceptable — they are read-only non-MySQL implementations). No changes needed there.

- [ ] TASK-012

---

## Phase 3 — Contribution Workflow

### TASK-013 · Database migration — `contributions` table

Create `database/migrations/005_contributions.sql`:

```sql
CREATE TABLE IF NOT EXISTS contributions (
    id             CHAR(36)     NOT NULL PRIMARY KEY,
    word_id        CHAR(36)     NOT NULL,
    contributor_id CHAR(36)     NOT NULL,
    reviewer_id    CHAR(36)     NULL,
    action         ENUM('submitted','approved','rejected','revised') NOT NULL,
    notes          TEXT         NULL,
    created_at     DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (contributor_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (reviewer_id)    REFERENCES users(id) ON DELETE SET NULL
);
```

- [ ] TASK-013

---

### TASK-014 · Contribution domain — model, repository, service

Create `domain/contribution/` package.

**`domain/contribution/model.go`**
```go
type Contribution struct {
    ID            string     `json:"id"`
    WordID        string     `json:"word_id"`
    Word          *Word      `json:"word,omitempty"`  // embedded for detail view
    ContributorID string     `json:"contributor_id"`
    ReviewerID    *string    `json:"reviewer_id,omitempty"`
    Action        string     `json:"action"`
    Notes         *string    `json:"notes,omitempty"`
    CreatedAt     time.Time  `json:"created_at"`
}

type SubmitRequest struct {
    Word        string          `json:"word"      validate:"required"`
    Syllable    string          `json:"syllables"`
    Alphabet    string          `json:"alphabet"  validate:"required,len=1"`
    Meanings    []WordMeaning   `json:"meanings"  validate:"required,min=1"`
    Derivatives []WordDerivative `json:"derivatives"`
}

type ReviewRequest struct {
    Notes string `json:"notes"`
}
```

**`domain/contribution/repository.go`** — interface:
```go
type Repository interface {
    CreateWord(w dictionary.Word, contributorID string) (string, error) // returns word ID
    UpdateWord(wordID string, w dictionary.Word) error
    SetWordStatus(wordID, status, reviewerID string, approvedAt *time.Time) error
    GetWordByID(wordID string) (dictionary.Word, string, error) // word + status
    LogAction(c Contribution) error
    ListByContributor(contributorID string, page, limit int) ([]Contribution, int, error)
    ListAll(status string, page, limit int) ([]Contribution, int, error)
    GetByID(id string) (Contribution, error)
    GetLatestByWordID(wordID string) (Contribution, error)
}
```

**`domain/contribution/repository_mysql.go`** — implement. `CreateWord` inserts into `words` table with `source='community'`, `status='pending'`.

**`domain/contribution/service.go`**
```go
type Service interface {
    Submit(contributorID string, req SubmitRequest) (Contribution, error)
    Edit(contributorID, wordID string, req SubmitRequest) (Contribution, error)
    Delete(contributorID, wordID string) error
    Mine(contributorID string, page, limit int) ([]Contribution, int, error)
    GetByID(contributorID, id string) (Contribution, error)
    // admin
    List(status string, page, limit int) ([]Contribution, int, error)
    Approve(reviewerID, id string) error
    Reject(reviewerID, id string, notes string) error
}
```

Service rules:
- `Submit`: validate `SubmitRequest`, serialize `Meanings`+`Derivatives` to JSON for `data` column, insert word + log `submitted` action
- `Edit`: caller must own the word AND word status must be `pending` or `rejected`; re-logs `revised` action, resets status to `pending`
- `Delete`: caller must own word AND status must be `pending`; hard delete word row + contribution rows
- `Approve`: set `status='active'`, set `approved_by`, `approved_at=NOW()`, log `approved`
- `Reject`: set `status='rejected'`, log `rejected` with notes

- [ ] TASK-014

---

### TASK-015 · Contribution HTTP handler + routes

Create `domain/contribution/handler.go`:

```
POST   /api/v1/contributions           → handler.Submit       (AuthMiddleware, role=user|admin)
GET    /api/v1/contributions/mine      → handler.Mine         (AuthMiddleware)
GET    /api/v1/contributions/:id       → handler.GetByID      (AuthMiddleware)
PUT    /api/v1/contributions/:id       → handler.Edit         (AuthMiddleware)
DELETE /api/v1/contributions/:id       → handler.Delete       (AuthMiddleware)

GET    /api/v1/admin/contributions                    → admin.ListContributions  (AuthMiddleware + RoleMiddleware("admin"))
PATCH  /api/v1/admin/contributions/:id/approve        → admin.Approve            (AuthMiddleware + RoleMiddleware("admin"))
PATCH  /api/v1/admin/contributions/:id/reject         → admin.Reject             (AuthMiddleware + RoleMiddleware("admin"))
```

Response for list endpoints:
```json
{
  "data": [...],
  "meta": { "page": 1, "limit": 20, "total": 45, "total_pages": 3 }
}
```

- [ ] TASK-015

---

### TASK-016 · Admin word management routes

Create `domain/admin/` or extend contribution handler. Admin-only word endpoints:

```
GET    /api/v1/admin/words          → list all words (all statuses), paginated, filterable by ?status=&source=
POST   /api/v1/admin/words          → add official word (source=official, status=active, no contribution record)
PUT    /api/v1/admin/words/:id      → edit any word
DELETE /api/v1/admin/words/:id      → soft-delete: set status='rejected' (or add status='deleted' variant)
```

All routes require `AuthMiddleware` + `RoleMiddleware("admin")`.

`POST /api/v1/admin/words` reuses `SubmitRequest` body but inserts directly as `source='official'`, `status='active'`, `contributor_id=NULL`.

- [ ] TASK-016

---

### TASK-017 · Tests — contribution service unit tests

Create `domain/contribution/service_test.go`. Mock repository. Cover:
- Successful submit
- Edit own pending word
- Edit another user's word → error
- Edit active word → error
- Delete pending word
- Delete active word → error
- Admin approve → status changes
- Admin reject → status changes + notes saved

- [ ] TASK-017

---

## Phase 4 — Community Features

### TASK-018 · Database migrations — votes, comments, bookmarks, word-of-the-day

Create `database/migrations/006_community.sql`:

```sql
CREATE TABLE IF NOT EXISTS word_votes (
    id         CHAR(36)  NOT NULL PRIMARY KEY,
    word_id    CHAR(36)  NOT NULL,
    user_id    CHAR(36)  NOT NULL,
    vote       TINYINT   NOT NULL,  -- 1 = up, -1 = down
    created_at DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uq_vote (word_id, user_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS word_comments (
    id         CHAR(36)  NOT NULL PRIMARY KEY,
    word_id    CHAR(36)  NOT NULL,
    user_id    CHAR(36)  NOT NULL,
    parent_id  CHAR(36)  NULL,
    body       TEXT      NOT NULL,
    created_at DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id)   REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES word_comments(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS word_bookmarks (
    id         CHAR(36)  NOT NULL PRIMARY KEY,
    word       VARCHAR(255) NOT NULL,
    user_id    CHAR(36)  NOT NULL,
    created_at DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uq_bookmark (word, user_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS word_of_the_day (
    id         CHAR(36)  NOT NULL PRIMARY KEY,
    word       VARCHAR(255) NOT NULL,
    date       DATE      NOT NULL UNIQUE,
    set_by     CHAR(36)  NULL,
    created_at DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

- [ ] TASK-018

---

### TASK-019 · Voting feature

Create `domain/community/vote.go` (or a shared `domain/community/` package).

**Repository interface methods:**
```go
Vote(wordID, userID string, vote int) error  // upsert; if same vote exists → delete (toggle)
GetVotes(wordID string) (up, down int, err error)
GetUserVote(wordID, userID string) (int, error)
```

**Service:**
```go
Cast(wordID, userID string, vote int) error  // vote must be 1 or -1
GetVoteSummary(wordID string) (up, down int, err error)
```

**Handler routes:**
```
POST /api/v1/entries/:word/votes   body: {"vote": 1}   (AuthMiddleware)
```

Upsert logic: if user already voted same direction → delete (toggle off). If different direction → update.

Update `GetWord()` public endpoint to include `"votes": {"up": N, "down": N}` in response. Votes fetched alongside word detail from MySQL only; non-MySQL repos return `"votes": null`.

- [ ] TASK-019

---

### TASK-020 · Comments feature

Add to `domain/community/` package.

**Model:**
```go
type Comment struct {
    ID        string    `json:"id"`
    WordID    string    `json:"word_id"`
    UserID    string    `json:"user_id"`
    UserName  string    `json:"user_name"`
    ParentID  *string   `json:"parent_id,omitempty"`
    Body      string    `json:"body"`
    Replies   []Comment `json:"replies,omitempty"`
    CreatedAt time.Time `json:"created_at"`
}
```

**Repository methods:**
```go
CreateComment(c Comment) error
DeleteComment(id, userID string, isAdmin bool) error  // isAdmin bypasses ownership check
ListComments(wordID string) ([]Comment, error)         // returns top-level + nested replies
```

**Handler routes:**
```
GET    /api/v1/entries/:word/comments          (public)
POST   /api/v1/entries/:word/comments          (AuthMiddleware) body: {"body":"...", "parent_id":"..."}
DELETE /api/v1/entries/:word/comments/:id      (AuthMiddleware) — owner or admin
DELETE /api/v1/admin/entries/:word/comments/:id (AuthMiddleware + RoleMiddleware("admin"))
```

`ListComments` query: fetch all comments for word, build tree in Go (not nested SQL). Top-level where `parent_id IS NULL`; replies where `parent_id = topLevel.id`.

- [ ] TASK-020

---

### TASK-021 · Bookmarks feature

Add to `domain/community/` package.

**Repository methods:**
```go
AddBookmark(word, userID string) error
RemoveBookmark(word, userID string) error
ListBookmarks(userID string, page, limit int) ([]string, int, error)
```

**Handler routes:**
```
GET    /api/v1/me/bookmarks        (AuthMiddleware) — paginated
POST   /api/v1/me/bookmarks/:word  (AuthMiddleware)
DELETE /api/v1/me/bookmarks/:word  (AuthMiddleware)
```

`AddBookmark`: ignore duplicate (upsert or check-then-insert). Return `200` on toggle-off.

- [ ] TASK-021

---

### TASK-022 · Word of the Day feature

Add to `domain/community/` package.

**Repository methods:**
```go
GetWordOfTheDay(date string) (string, error)  // date = "YYYY-MM-DD"
SetWordOfTheDay(word, setBy, date string) error
GetRandomActiveWord() (string, error)
```

**Service:**
- `Get()`: try DB for today's date. If miss, call `GetRandomActiveWord()`, persist it, return it
- `Set(word, adminID, date string)`: admin manually schedules a word for a given date (typically tomorrow)

**Handler routes:**
```
GET /api/v1/word-of-the-day              (public, cache-control: max-age=3600)
PUT /api/v1/admin/word-of-the-day        (AuthMiddleware + RoleMiddleware("admin")) body: {"word":"...", "date":"YYYY-MM-DD"}
```

Response: full word detail (same shape as `GET /api/v1/entries/:word`).

`GetRandomActiveWord()` SQL: `SELECT word FROM words WHERE status='active' ORDER BY RAND() LIMIT 1`

- [ ] TASK-022

---

## Phase 5 — Admin Dashboard & User Management

### TASK-023 · Admin user management routes

Wire into `domain/user/handler.go` (add admin handler group) or new `domain/admin/` package:

```
GET   /api/v1/admin/users               → list paginated users
PATCH /api/v1/admin/users/:id/deactivate → set is_active=0
PATCH /api/v1/admin/users/:id/activate   → set is_active=1
PATCH /api/v1/admin/users/:id/promote    → set role='admin'
```

All require `AuthMiddleware` + `RoleMiddleware("admin")`.

`ListUsers` supports `?page=&limit=&role=&active=` filters.

Admin cannot deactivate themselves. Add guard: if `userID == c.Locals("userID")` return `400 cannot modify own account`.

- [ ] TASK-023

---

### TASK-024 · Admin stats endpoint

Create `GET /api/v1/admin/stats`. Runs 5–6 cheap SQL COUNT queries (or a single query with CASE). Response:

```json
{
  "words": {
    "official": 1200,
    "community": 45,
    "pending": 12,
    "rejected": 8
  },
  "users": {
    "total": 320,
    "active": 310,
    "inactive": 10
  },
  "contributions": {
    "pending": 12,
    "this_week": 5,
    "this_month": 18
  },
  "top_contributors": [
    { "user_id": "...", "name": "...", "approved_count": 10 }
  ]
}
```

Top contributors query:
```sql
SELECT u.id, u.name, COUNT(c.id) as approved_count
FROM contributions c
JOIN users u ON c.contributor_id = u.id
WHERE c.action = 'approved'
GROUP BY u.id, u.name
ORDER BY approved_count DESC
LIMIT 5;
```

- [ ] TASK-024

---

## Phase 6 — Hardening & Cross-Cutting

### TASK-025 · Pagination helper

Create `pkg/pagination/pagination.go`:

```go
type Meta struct {
    Page       int `json:"page"`
    Limit      int `json:"limit"`
    Total      int `json:"total"`
    TotalPages int `json:"total_pages"`
}

func Parse(c *fiber.Ctx) (page, limit int) {
    page = max(1, c.QueryInt("page", 1))
    limit = min(100, max(1, c.QueryInt("limit", 20)))
    return
}

func NewMeta(page, limit, total int) Meta {
    pages := (total + limit - 1) / limit
    return Meta{Page: page, Limit: limit, Total: total, TotalPages: pages}
}
```

Use in all list handlers.

- [ ] TASK-025

---

### TASK-026 · Unified error response

Ensure all handlers return consistent error JSON:

```json
{ "error": "message here" }
```

Create `pkg/response/response.go` with helpers:
```go
func Error(c *fiber.Ctx, status int, msg string) error {
    return c.Status(status).JSON(fiber.Map{"error": msg})
}

func OK(c *fiber.Ctx, data interface{}) error {
    return c.JSON(fiber.Map{"data": data})
}

func Paginated(c *fiber.Ctx, data interface{}, meta pagination.Meta) error {
    return c.JSON(fiber.Map{"data": data, "meta": meta})
}
```

- [ ] TASK-026

---

### TASK-027 · Guard non-MySQL write operations

In embed/FS/JSON repositories, add stub methods for all new write operations that return a sentinel error `ErrNotSupported`. In `main.go`, when auth middleware is hit and DB is not MySQL, the handler layer returns `503 Service Unavailable — requires MySQL backend`.

Actually: simpler approach — only register auth/contribution/admin routes when `MYSQL_DSN` is set. Skip route registration entirely otherwise. Document this in README.

- [ ] TASK-027

---

### TASK-028 · Update README + env var docs

Update `README.md` (or create if missing) to document:
- New env vars: `JWT_SECRET`, `JWT_ACCESS_TTL`, `JWT_REFRESH_TTL`, `ADMIN_EMAIL`, `ADMIN_PASSWORD`
- Auth routes require MySQL (`MYSQL_DSN` must be set)
- Quick-start with Docker Compose (MySQL variant already exists in `docker-compose.mysql.yml`)
- Migration order: run `001` → `002` → `003` → `004` → `005` → `006`

- [ ] TASK-028

---

### TASK-029 · Integration tests — auth + contribution flow

Create `integration/` or `e2e/` test file (build tag `integration`). Uses a real test DB (separate DSN from env). Test full flows:

1. Register → login → get token
2. Submit word → verify pending status
3. Admin login → approve → verify active
4. Public endpoint → verify approved word appears
5. Reject flow → verify word not in public results

Run with: `go test -tags=integration -v ./integration/...`

- [ ] TASK-029

---

### TASK-030 · Docker Compose updates

Update `docker-compose.mysql.yml` and/or `docker-compose.prod.yml`:
- Add `JWT_SECRET` env var (required)
- Add `ADMIN_EMAIL` + `ADMIN_PASSWORD` env vars
- Update `docker-entrypoint-initdb.d` to include new migration files (004, 005, 006)
- Update `.env.mysql.example` with new vars + comments

- [ ] TASK-030

---

## Summary

| Phase | Tasks | Description |
|---|---|---|
| 1 — Auth | TASK-001 – 008 | Dependencies, DB tables, user domain, JWT middleware, auth routes |
| 2 — Source Tracking | TASK-009 – 012 | Migrate words table, filter active, expose source field |
| 3 — Contributions | TASK-013 – 017 | Contribution domain, submit/review flow, admin word mgmt |
| 4 — Community | TASK-018 – 022 | Votes, comments, bookmarks, word of the day |
| 5 — Admin Dashboard | TASK-023 – 024 | User management, stats endpoint |
| 6 — Hardening | TASK-025 – 030 | Pagination helper, error response, non-MySQL guards, docs, Docker |

**Recommended order:** follow phase order. Phase 2 can start in parallel with Phase 1 (no dependency on auth code). Phase 3 depends on both Phase 1 and 2.

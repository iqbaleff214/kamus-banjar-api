# Changelog

## [Unreleased]

## [2.0.0] - 2026-04-18

> **Breaking change:** MySQL is now the primary data backend. All endpoints below the dictionary read-only layer (`/alphabets`, `/entries`) require `MYSQL_DSN` and `JWT_SECRET` to be set. The embedded/FS dictionary still works without MySQL — those endpoints remain unchanged.

### Added

#### Auth
- User registration, login, logout, and token refresh (`POST /auth/register`, `/auth/login`, `/auth/logout`, `/auth/refresh`)
- JWT-based authentication with access token (default 15 min) and refresh token (default 7 days) rotation
- Refresh tokens stored as SHA-256 hashes — safe against DB breach replay
- `GET /auth/me` and `PUT /auth/me` for profile retrieval and update (including password change)
- `middleware.Auth` (validates Bearer JWT) and `middleware.Role` (RBAC guard)

#### Word Source Tracking
- Words now carry `source` (`official` | `community`) and `status` (`active` | `pending` | `rejected`) fields
- Contributor and reviewer audit columns (`contributor_id`, `approved_by`, `approved_at`) on the words table
- Public endpoints automatically filter to `status = active` only

#### Contribution Workflow
- Users can submit new word contributions (`POST /contributions`)
- Submitted words enter `pending` status; contributors can edit or delete while pending
- `GET /contributions/mine` for paginated personal history; `GET /contributions/:id` for detail
- Admin contribution review: `PATCH /admin/contributions/:id/approve` and `/reject`
- Admin word management: `GET|POST /admin/words`, `PUT|DELETE /admin/words/:id`

#### Community Features
- **Voting:** `POST /entries/:word/votes` — upvote (+1) or downvote (−1) with toggle support; vote counts included in `GET /entries/:word`
- **Comments:** threaded comments on dictionary entries — `GET|POST /entries/:word/comments`, `DELETE /entries/:word/comments/:id`
- **Bookmarks:** `GET|POST|DELETE /me/bookmarks` and `/me/bookmarks/:word`
- **Word of the Day:** `GET /word-of-the-day` (public, cached 1 h); `PUT /admin/word-of-the-day` for admin scheduling

#### Admin Dashboard & User Management
- `GET /admin/stats` — platform snapshot: word counts by source/status, user counts, contribution activity, top 5 contributors
- `GET /admin/users` with pagination and role/active filters
- `PATCH /admin/users/:id/deactivate`, `/activate`, `/promote`
- Self-modification guard: admins cannot deactivate their own account

#### Developer Experience
- `pkg/pagination` — shared `Parse` and `NewMeta` helpers used across all list endpoints
- `pkg/response` — unified `OK`, `Created`, `NoContent`, `Paginated` helpers; all responses follow `{code, status, message, data}` envelope
- Integration test suite (`integration/` package, build tag `integration`) covering register → login → submit → approve flows
- `docker-compose.yml` (production, Traefik + MySQL) and `docker-compose.dev.yml` (development, direct ports)
- `.env.example` and `.env.dev.example` templates

### Changed
- Project restructured into `internal/` clean-architecture packages: `config`, `dictionary`, `user`, `contribution`, `community`, `admin`, `middleware`, `server`, `seeder`
- `GET /entries/:word` now includes a `votes` field (`{up, down}`) when community features are enabled
- All paginated list responses include a `meta` object (`{page, limit, total, total_pages}`)
- `docker-compose.yml` promoted to the canonical production compose file (Traefik + MySQL); old `docker-compose.prod.yml` and `docker-compose.mysql.yml` removed

### Fixed
- Conditional route registration: auth, community write, contribution, and admin routes are only registered when `MYSQL_DSN` is set — the server starts cleanly in dictionary-only mode otherwise

### Migration
Run migrations in order against your MySQL database before starting v2.0.0:
```
001_schema.sql
002_seed.sql  (seeds)
003_refresh_tokens.sql
004_word_source.sql
005_contributions.sql
006_community.sql
```

---

## [1.1.0] - 2025-03-06

### Added
- Added Dockerfile for containerization
- Added `.vscode` configuration for development setup
- Added endpoint for word search using Levenshtein algorithm
- Improved test cases for word search service

### Changed
- Optimized cache structure from slice/array to hash table for better search performance

### Fixed
- Fixed capitalization issue in vocabulary data (converted uppercase words to lowercase)

# Product Requirements Document
## Kamus Banjar API — Community Edition

**Version:** 1.0  
**Date:** 2026-04-17  
**Author:** M. Iqbal Effendi  
**Status:** Draft

---

## 1. Overview

Kamus Banjar API is a REST API for the Banjar language dictionary. Currently it serves read-only dictionary data from static JSON files or MySQL. This PRD defines the roadmap to evolve it into a **community-driven dictionary platform** with user accounts, contribution workflows, and moderation tools — while preserving the existing read-only public API surface.

---

## 2. Goals

- Add user authentication and role-based access control (RBAC)
- Allow community members to contribute new words and definitions
- Give admins tools to review, approve, reject, and manage contributions
- Track the source/origin of every word entry
- Maintain API backward compatibility for all existing public endpoints

---

## 3. Non-Goals

- Mobile app (API-first; mobile is consumer's concern)
- Real-time features (WebSocket, push notifications)
- AI-generated definitions
- Paid tiers or monetization

---

## 4. User Roles

| Role | Description |
|---|---|
| **Guest** | Unauthenticated. Read-only access to public endpoints. |
| **User** | Registered community member. Can submit word contributions, track own submissions, vote on entries. |
| **Admin** | Dictionary curator. Can approve/reject contributions, add official words, manage users. |

Admins are created by seeding or by another admin promoting a user — no self-registration as admin.

---

## 5. Data Model Changes

### 5.1 `Word` — New Fields

| Field | Type | Values | Notes |
|---|---|---|---|
| `source` | enum | `official`, `community` | Set by system, not user input |
| `status` | enum | `active`, `pending`, `rejected` | `active` = visible in public API |
| `contributor_id` | UUID / nullable | — | FK to user; null for seeded official data |
| `approved_by` | UUID / nullable | — | FK to admin user |
| `approved_at` | timestamp / nullable | — | Set on approval |
| `created_at` | timestamp | — | |
| `updated_at` | timestamp | — | |

**Source rules:**
- Admin adds word → `source: official`, `status: active` immediately
- User submits word → `source: community`, `status: pending` until admin approves
- On approval → `status: active`; on rejection → `status: rejected`

Public read endpoints (`/api/v1/entries`, `/api/v1/alphabets`) only return `status: active` words.

### 5.2 `User` — New Model

| Field | Type | Notes |
|---|---|---|
| `id` | UUID | PK |
| `name` | string | Display name |
| `email` | string | Unique, used for login |
| `password_hash` | string | bcrypt |
| `role` | enum | `admin`, `user` |
| `is_active` | bool | Admin can deactivate |
| `created_at` | timestamp | |
| `updated_at` | timestamp | |

### 5.3 `Contribution` — Audit/Tracking Model

Tracks every submission with full history (submit → review → decision).

| Field | Type | Notes |
|---|---|---|
| `id` | UUID | PK |
| `word_id` | UUID | FK to word (created on submit) |
| `contributor_id` | UUID | FK to user |
| `reviewer_id` | UUID / nullable | FK to admin |
| `action` | enum | `submitted`, `approved`, `rejected`, `revised` |
| `notes` | text / nullable | Reviewer feedback on rejection |
| `created_at` | timestamp | |

---

## 6. Authentication

- **Method:** JWT (access token + refresh token)
- **Access token TTL:** 15 minutes
- **Refresh token TTL:** 7 days, stored server-side (DB or Redis) for revocation
- **Transport:** `Authorization: Bearer <token>` header

### 6.1 Auth Endpoints

```
POST /api/v1/auth/register     # Create user account
POST /api/v1/auth/login        # Get access + refresh token
POST /api/v1/auth/refresh      # Rotate tokens
POST /api/v1/auth/logout       # Revoke refresh token
GET  /api/v1/auth/me           # Current user profile
PUT  /api/v1/auth/me           # Update profile (name, password)
```

---

## 7. Feature Requirements

### 7.1 User Registration & Login

**FR-AUTH-01:** Users can register with name, email, password.  
**FR-AUTH-02:** Email must be unique. Duplicate registration returns `409 Conflict`.  
**FR-AUTH-03:** Password min 8 chars; stored as bcrypt hash (cost ≥ 12).  
**FR-AUTH-04:** Login returns `access_token` + `refresh_token`.  
**FR-AUTH-05:** Refresh endpoint rotates both tokens (refresh token rotation).  
**FR-AUTH-06:** Logout invalidates the refresh token server-side.

### 7.2 Word Contribution (User)

**FR-CONTRIB-01:** Authenticated users can submit a new word via `POST /api/v1/contributions`.  
**FR-CONTRIB-02:** Submission creates a `Word` with `status: pending`, `source: community`.  
**FR-CONTRIB-03:** Pending words are NOT visible in public read endpoints.  
**FR-CONTRIB-04:** Users can view their own submissions (all statuses) via `GET /api/v1/contributions/mine`.  
**FR-CONTRIB-05:** Users can edit a `pending` or `rejected` submission and resubmit. Status resets to `pending`.  
**FR-CONTRIB-06:** Users cannot delete an `active` (approved) word.

### 7.3 Admin — Word Management

**FR-ADMIN-01:** Admins can add words directly → `status: active`, `source: official`.  
**FR-ADMIN-02:** Admins can edit any word (official or community).  
**FR-ADMIN-03:** Admins can soft-delete (deactivate) any word.  
**FR-ADMIN-04:** Admins can view all words including `pending` and `rejected`.

### 7.4 Admin — Contribution Review

**FR-REVIEW-01:** Admins can list all pending contributions via `GET /api/v1/admin/contributions?status=pending`.  
**FR-REVIEW-02:** Admin approves via `PATCH /api/v1/admin/contributions/:id/approve` → word `status` → `active`.  
**FR-REVIEW-03:** Admin rejects via `PATCH /api/v1/admin/contributions/:id/reject` with optional `notes`.  
**FR-REVIEW-04:** Contributor can see rejection notes on their submission.  
**FR-REVIEW-05:** Audit log entry created on every review action.

### 7.5 Admin — User Management

**FR-USER-01:** Admins can list users (`GET /api/v1/admin/users`).  
**FR-USER-02:** Admins can deactivate/reactivate users.  
**FR-USER-03:** Admins can promote a user to admin role.  
**FR-USER-04:** Deactivated users cannot log in.

### 7.6 Word Voting / Rating (Community Feature)

**FR-VOTE-01:** Authenticated users can upvote/downvote `active` word entries.  
**FR-VOTE-02:** One vote per user per word; repeat call toggles the vote.  
**FR-VOTE-03:** Vote counts exposed on word detail endpoint: `votes: { up: N, down: N }`.  
**FR-VOTE-04:** Admins can see vote counts to gauge community consensus before approving contributions.

### 7.7 Word Comments / Discussion (Community Feature)

**FR-COMMENT-01:** Authenticated users can comment on `active` word entries.  
**FR-COMMENT-02:** Comments are public (visible to guests too).  
**FR-COMMENT-03:** Commenter can delete their own comment. Admin can delete any.  
**FR-COMMENT-04:** Comments support basic threading (one level: reply to comment).

### 7.8 Bookmarks / Favorites (User Feature)

**FR-BOOKMARK-01:** Authenticated users can bookmark words.  
**FR-BOOKMARK-02:** `GET /api/v1/me/bookmarks` returns paginated bookmark list.  
**FR-BOOKMARK-03:** `POST /api/v1/me/bookmarks/:word` adds; `DELETE` removes.

### 7.9 Word of the Day

**FR-WOTD-01:** System selects one `active` word per day (random or curated by admin).  
**FR-WOTD-02:** Exposed at `GET /api/v1/word-of-the-day` — cached, no auth required.  
**FR-WOTD-03:** Admin can manually set tomorrow's word via `PUT /api/v1/admin/word-of-the-day`.

### 7.10 Statistics / Dashboard (Admin)

**FR-STATS-01:** `GET /api/v1/admin/stats` returns:
- Total words by source (`official` / `community`)
- Total pending contributions
- Total users (active / inactive)
- Top contributors (by approved submissions)
- Words added this week/month

---

## 8. API Endpoints Summary

### Public (no auth)

```
GET  /api/v1/                           # API info (unchanged)
GET  /api/v1/alphabets                  # All letters (unchanged)
GET  /api/v1/alphabets/:letter          # Words by letter (unchanged)
GET  /api/v1/entries/:word              # Word detail (unchanged, now includes source/votes)
GET  /api/v1/entries?search=keyword     # Fuzzy search (unchanged)
GET  /api/v1/word-of-the-day            # Word of the day
```

### Auth

```
POST /api/v1/auth/register
POST /api/v1/auth/login
POST /api/v1/auth/refresh
POST /api/v1/auth/logout
GET  /api/v1/auth/me
PUT  /api/v1/auth/me
```

### Authenticated User

```
POST   /api/v1/contributions                    # Submit new word
GET    /api/v1/contributions/mine               # Own submissions
GET    /api/v1/contributions/:id                # Single submission detail
PUT    /api/v1/contributions/:id                # Edit pending/rejected submission
DELETE /api/v1/contributions/:id                # Delete own pending submission

POST   /api/v1/entries/:word/votes              # Upvote/downvote
GET    /api/v1/entries/:word/comments           # List comments
POST   /api/v1/entries/:word/comments           # Add comment
DELETE /api/v1/entries/:word/comments/:id       # Delete own comment

GET    /api/v1/me/bookmarks                     # List bookmarks
POST   /api/v1/me/bookmarks/:word               # Add bookmark
DELETE /api/v1/me/bookmarks/:word               # Remove bookmark
```

### Admin Only

```
GET    /api/v1/admin/words                              # All words (all statuses)
POST   /api/v1/admin/words                              # Add official word
PUT    /api/v1/admin/words/:id                          # Edit any word
DELETE /api/v1/admin/words/:id                          # Soft-delete word

GET    /api/v1/admin/contributions                      # All contributions (filterable by status)
PATCH  /api/v1/admin/contributions/:id/approve          # Approve contribution
PATCH  /api/v1/admin/contributions/:id/reject           # Reject contribution

GET    /api/v1/admin/users                              # List users
PATCH  /api/v1/admin/users/:id/deactivate               # Deactivate user
PATCH  /api/v1/admin/users/:id/activate                 # Activate user
PATCH  /api/v1/admin/users/:id/promote                  # Promote to admin

GET    /api/v1/admin/stats                              # Dashboard stats
PUT    /api/v1/admin/word-of-the-day                    # Set tomorrow's word

DELETE /api/v1/admin/entries/:word/comments/:id         # Delete any comment
```

---

## 9. Technical Considerations

### 9.1 Repository Layer

Current architecture has 4 repository implementations (embed, fs, mysql, json-remote). For all new features (users, contributions, votes, comments, bookmarks):
- **MySQL repository** is the only viable implementation (stateful data needs persistence)
- Embed/FS/JSON-remote repositories → return `ErrNotSupported` for write operations
- Consider requiring `MYSQL_DSN` for any authenticated endpoint; graceful 503 otherwise

### 9.2 Database Schema

New tables required:
- `users`
- `contributions` (audit trail)
- `word_votes`
- `word_comments`
- `word_bookmarks`
- `word_of_the_day`

Existing `words` table in MySQL repo needs migration:
- Add columns: `source`, `status`, `contributor_id`, `approved_by`, `approved_at`, `created_at`, `updated_at`
- Backfill: all existing rows → `source: official`, `status: active`

### 9.3 Middleware

- `AuthMiddleware` — validates JWT, injects user into context
- `RoleMiddleware(role)` — checks role, returns `403` if insufficient
- Rate limiting on auth endpoints (prevent brute force): 10 req/min per IP on `/auth/login` and `/auth/register`

### 9.4 Password & Security

- bcrypt with cost 12 minimum
- No password in any log or response
- JWT secret via env var `JWT_SECRET` (required; app exits if unset)
- Refresh tokens stored hashed in DB

### 9.5 Environment Variables (New)

| Variable | Default | Description |
|---|---|---|
| `JWT_SECRET` | — | Required. JWT signing secret (min 32 chars) |
| `JWT_ACCESS_TTL` | `15m` | Access token TTL |
| `JWT_REFRESH_TTL` | `168h` | Refresh token TTL (7 days) |
| `ADMIN_EMAIL` | — | Seed admin email on first startup |
| `ADMIN_PASSWORD` | — | Seed admin password on first startup |

### 9.6 Pagination

All list endpoints that can return large sets must support:
- `?page=1&limit=20` (default limit 20, max 100)
- Response includes `meta: { page, limit, total, total_pages }`

### 9.7 Backward Compatibility

- All existing `GET /api/v1/*` endpoints remain unchanged in structure
- New fields (`source`, `votes`) added to word response are additive — no breaking change
- `status` field NOT exposed in public endpoints (only `active` words returned)

---

## 10. Migration Plan

### Phase 1 — Auth Foundation
1. Add `users` table + migration
2. Implement register/login/refresh/logout/me endpoints
3. Seed initial admin from env vars
4. Add `AuthMiddleware` + `RoleMiddleware`

### Phase 2 — Word Source Tracking
1. Migrate `words` table: add `source`, `status`, `contributor_id`, `approved_by`, `approved_at`
2. Backfill existing words → `official` / `active`
3. Update public read endpoints to filter `status: active`

### Phase 3 — Contribution Workflow
1. Add `contributions` table
2. User submit endpoint
3. Admin review endpoints (approve/reject)
4. Contributor can view own submissions + rejection notes

### Phase 4 — Community Features
1. Voting
2. Comments
3. Bookmarks
4. Word of the Day

### Phase 5 — Admin Dashboard
1. Stats endpoint
2. User management endpoints

---

## 11. Success Metrics

| Metric | Target |
|---|---|
| Registered users (6 months post-launch) | 500+ |
| Community word submissions | 200+ |
| Approval rate | > 60% of submissions |
| Public API uptime | 99.5% |
| P95 response time | < 100ms (read endpoints) |

---

## 12. Open Questions

1. Should community-contributed words be searchable while `pending` for the submitter themselves, or completely hidden until approved?
2. Should rejected submissions be permanently deleted after N days, or kept indefinitely for audit?
3. Is email verification required on registration, or trust-on-register?
4. Should vote counts influence search ranking in the fuzzy search results?
5. Multi-language support for definitions (currently Banjar ↔ Indonesian only)?

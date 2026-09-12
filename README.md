# EduCore

A Go backend for an education management system — currently a REST API covering authentication, role-based access, and student records on PostgreSQL. Built as a from-scratch exercise in Go's standard `net/http`, without a framework, to understand what a framework like Gin or Echo actually does under the hood.

**Status: early-stage / actively developed.** This is not yet a full platform — see [Roadmap](#roadmap) for what's not built yet.

## What it does today

- **Auth** — registration with bcrypt-hashed passwords, login that issues a signed JWT
- **Role-based access control** — middleware that validates a bearer JWT and can gate a route to a specific role (e.g. `admin`)
- **Student records** — list students from PostgreSQL over a JSON API

## Tech stack

| Layer | Choice |
|---|---|
| Language | Go 1.27 |
| HTTP | standard library `net/http` (no framework) |
| Database | PostgreSQL via [`pgx`](https://github.com/jackc/pgx) |
| Auth | [`golang-jwt`](https://github.com/golang-jwt/jwt) + `golang.org/x/crypto/bcrypt` |
| Config | `.env` via [`godotenv`](https://github.com/joho/godotenv) |

## Project structure

```
EduCore/
├── main.go        # HTTP routes + handlers (register, login, students, profile, admin)
├── auth.go        # password hashing (bcrypt)
├── token.go       # JWT creation
├── middleware.go  # auth + role-gating middleware
├── db.go          # PostgreSQL connection (pgx)
└── student/       # student domain types (in progress)
```

## Getting started

**Requirements:** Go 1.27+, a running PostgreSQL instance.

1. Clone the repo and install dependencies:
   ```bash
   git clone https://github.com/AHmed-Alikar/EduCore.git
   cd EduCore
   go mod download
   ```
2. Copy `.env.example` to `.env` and fill in real values — a database URL and a JWT signing secret (`openssl rand -base64 32` generates a reasonable one):
   ```bash
   cp .env.example .env
   ```
3. Create the tables the API expects:
   ```sql
   CREATE TABLE students (
     id INT PRIMARY KEY,
     name TEXT NOT NULL,
     age INT NOT NULL
   );

   CREATE TABLE users (
     id SERIAL PRIMARY KEY,
     name TEXT NOT NULL,
     email TEXT UNIQUE NOT NULL,
     password_hash TEXT NOT NULL,
     role TEXT NOT NULL DEFAULT 'student'
   );
   ```
4. Run the server:
   ```bash
   go run .
   ```
   The API listens on `http://localhost:8080`.

## API

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| POST | `/register` | — | Create a user (name, email, password) |
| POST | `/login` | — | Authenticate, returns a JWT |
| GET | `/students` | — | List all students |
| GET | `/profile` | Bearer token | Any authenticated user |
| GET | `/admin` | Bearer token, `admin` role | Role-gated example route |

## Testing

Password hashing/verification is covered in `auth_test.go`. Because
`token.go`'s `init()` requires `JWT_SECRET` to be set (see
[Known limitations](#known-limitations)), tests need it too:

```bash
JWT_SECRET=any-value-for-tests go test ./...
```

## Roadmap

Ordered by what's actually planned next, not by ambition:

- [ ] Wire up full student CRUD to HTTP (create/update/delete currently exist as internal functions only, not exposed as routes)
- [ ] Teachers, courses, attendance, and grades — the rest of the domain model
- [ ] Input validation on `/register` and `/login`
- [ ] Automated tests for HTTP handlers and the database layer (currently only password hashing is covered)
- [ ] React frontend
- [ ] CI (lint + test on PR)

## Known limitations

- `createStudent`, `updateStudent`, and `deleteStudent` exist in `main.go` but aren't wired to HTTP routes yet — today only reads are exposed.
- Test coverage is limited to password hashing (`auth_test.go`) — the HTTP handlers and database layer aren't tested yet.
- No input validation on `/register` or `/login` beyond what the database schema enforces (e.g. `email` uniqueness).

## License

Not yet licensed — treat as source-available for now.

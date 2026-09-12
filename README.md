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
2. Create a `.env` file in the project root:
   ```
   DATABASE_URL=postgres://user:password@localhost:5432/educore?sslmode=disable
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

## Roadmap

Ordered by what's actually planned next, not by ambition:

- [ ] Wire up full student CRUD to HTTP (create/update/delete currently exist as internal functions only, not exposed as routes)
- [ ] Teachers, courses, attendance, and grades — the rest of the domain model
- [ ] Move the JWT signing secret out of source code and into environment configuration
- [ ] Input validation on `/register` and `/login`
- [ ] Automated tests (currently none)
- [ ] React frontend
- [ ] CI (lint + test on PR)

## Known limitations

- The JWT secret is currently hardcoded in `token.go`. It needs to move to an environment variable before this is ever run anywhere but locally.
- `createStudent`, `updateStudent`, and `deleteStudent` exist in `main.go` but aren't wired to HTTP routes yet — today only reads are exposed.
- No automated tests yet.

## License

Not yet licensed — treat as source-available for now.

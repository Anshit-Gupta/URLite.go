# URLite.go

A URL shortener written in Go. Send it a long URL, get back a short one; visit the short one and you're redirected to the original.

This is the Go version of my earlier [URLite](https://github.com/) project (Next.js + Prisma). I already knew how a URL shortener works, so rebuilding it was a way to focus purely on what's different in Go: writing the HTTP server side by hand with `net/http`, handling JSON without a framework, talking to Postgres directly, and testing handlers with `httptest`. It's deliberately scoped to the core only. No auth, no analytics.

## API

**Create a short URL**

```
POST /createShortURL
Content-Type: application/json

{"LongURL": "https://example.com/some/very/long/path"}
```

Response:

```json
{"ShortURL": "http://localhost:8080/aZ3x9F"}
```

**Use a short URL**

```
GET /aZ3x9F
```

Responds with a `302 Found` redirecting to the original URL. If the code doesn't exist, you get a `404`.

## How it works

- A short code is 6 random characters picked from `a-z`, `A-Z`, `0-9`.
- Codes and their long URLs are stored in a Postgres table, with the code as the primary key, so every lookup on redirect is a direct primary-key read.
- The app talks to Postgres through a `pgxpool` connection pool. Since `net/http` runs every request in its own goroutine, the pool is what lets many requests hit the database at once safely.
- The very first version stored everything in an in-memory map guarded by a `sync.Mutex`. Moving to Postgres removed the need for the mutex entirely, since the database now owns concurrent access to the data.

## Project structure

```text
URLite.go/
├── main.go        # handlers, code generation, server setup
├── main_test.go   # tests
├── go.mod
├── go.sum
└── .env           # DATABASE_URL (not committed)
```

## Setup

You'll need Go and a Postgres database. I used [Neon](https://neon.tech)'s free tier, but any Postgres works.

**1. Create the table**

```sql
CREATE TABLE urls (
    code       TEXT PRIMARY KEY,
    long_url   TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

**2. Add your connection string**

Create a `.env` file in the project root:

```text
DATABASE_URL=postgres://user:password@host/dbname?sslmode=require
```

Make sure `.env` is in your `.gitignore`. It contains your database password.

**3. Run it**

```bash
go run main.go
```

You should see a "database connection successful" message, then the server starts on port `8080`.

## Trying it out

```bash
curl -X POST localhost:8080/createShortURL -d '{"LongURL":"https://example.com"}'
```

Then open the `ShortURL` you get back in a browser.

On Windows PowerShell, `curl` is an alias for a different command and the flags won't work. Use `curl.exe` instead, or PowerShell's native version:

```powershell
Invoke-RestMethod -Uri "http://localhost:8080/createShortURL" -Method POST -Body '{"LongURL":"https://example.com"}' -ContentType "application/json"
```

## Tests

```bash
go test -v
```

There are three tests:

- `TestGenerateShortCode` checks that a code is 6 characters long and only uses characters from the allowed charset.
- `TestCreateShortURL` sends a real POST through the handler and checks the response, then queries the database to confirm the right row was stored.
- `TestRedirect` inserts a known row, calls the redirect handler, and checks for a `302` and the correct `Location` header.

The handler tests use `net/http/httptest`, so no server needs to be running. They do run against the real database from your `.env`, so `DATABASE_URL` must be set. Each test inserts its own data and deletes it afterward, so nothing is left behind in the table.


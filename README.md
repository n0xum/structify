# structify

CLI tool to convert Go domain models to PostgreSQL database schemas and database/sql CRUD code.

## Installation

```bash
go install github.com/ak/structify/cmd/structify@latest
```

Or build from source:

```bash
git clone https://github.com/ak/structify.git
cd structify
go build ./cmd/structify
```

## Setup

1. Create Go struct files with your domain models
2. Add struct tags to define database behavior
3. Run structify with the appropriate flag

## Struct Tags

Add `db` tags to your struct fields to control database behavior:

- `db:"pk"` - Mark field as primary key
- `db:"unique"` - Add UNIQUE constraint
- `db:"-"` - Exclude field from database
- `db:"table:name"` - Override table name

## Usage

Generate PostgreSQL schema:

```bash
structify --to-sql ./models/user.go
```

Generate database/sql CRUD code:

```bash
structify --to-db-sql ./models/user.go
```

Parse and validate without generating output:

```bash
structify ./models/user.go
```

Write output to file:

```bash
structify --to-sql ./models/user.go -o schema.sql
structify --to-db-sql ./models/user.go -o user_repo.go
```

Process multiple files:

```bash
structify --to-sql ./models/*.go
```

## Type Mapping

| Go Type | PostgreSQL Type |
|---------|----------------|
| int64 | BIGINT |
| string | VARCHAR(255) |
| bool | BOOLEAN |
| float64 | DOUBLE PRECISION |
| time.Time | TIMESTAMP |
| []byte | BYTEA |

## Example

Input file `models/user.go`:

```go
package models

type User struct {
    ID       int64     `db:"pk"`
    Username string    `db:"unique"`
    Email    string
    Active   bool
    Created  int64
}
```

Run:

```bash
structify --to-sql ./models/user.go
```

Output:

```sql
CREATE TABLE "user" (
    "id" BIGINT PRIMARY KEY,
    "username" VARCHAR(255) UNIQUE,
    "email" VARCHAR(255),
    "active" BOOLEAN NOT NULL,
    "created" BIGINT NOT NULL
);
```

Run:

```bash
structify --to-db-sql ./models/user.go
```

Output:

```go
package main

import (
    "database/sql"
    "context"
)

type User struct {
    ID int64
    Username string
    Email string
    Active bool
    Created int64
}

func CreateUser(ctx context.Context, db *sql.DB, item *User) (*User, error) {
    query := `INSERT INTO user (username, email, active, created) VALUES ($1, $2, $3, $4) RETURNING id`
    var id int64
    err := db.QueryRowContext(ctx, query, item.Username, item.Email, item.Active, item.Created).Scan(&id)
    if err != nil {
        return nil, err
    }
    return GetUserByID(ctx, db, id)
}

func GetUserByID(ctx context.Context, db *sql.DB, id int64) (*User, error) {
    query := `SELECT id, username, email, active, created FROM user WHERE id = $1`
    var item User
    err := db.QueryRowContext(ctx, query, id).Scan(&item.ID, &item.Username, &item.Email, &item.Active, &item.Created)
    if err != nil {
        return nil, err
    }
    return &item, nil
}

func UpdateUser(ctx context.Context, db *sql.DB, item *User) error {
    query := `UPDATE user SET username = $1, email = $2, active = $3, created = $4 WHERE id = $5`
    _, err := db.ExecContext(ctx, query, item.Username, item.Email, item.Active, item.Created, item.ID)
    return err
}

func DeleteUser(ctx context.Context, db *sql.DB, id int64) error {
    query := `DELETE FROM user WHERE id = $1`
    _, err := db.ExecContext(ctx, query, id)
    return err
}

func ListUser(ctx context.Context, db *sql.DB) ([]*User, error) {
    query := `SELECT id, username, email, active, created FROM user ORDER BY id`
    rows, err := db.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []*User
    for rows.Next() {
        var item User
        if err := rows.Scan(&item.ID, &item.Username, &item.Email, &item.Active, &item.Created); err != nil {
            return nil, err
        }
        items = append(items, &item)
    }
    return items, nil
}
```

## Flags

- `--to-sql, --to-schema` - Generate PostgreSQL CREATE TABLE statements
- `--to-db-sql, --to-dbcode` - Generate database/sql CRUD code
- `--from-json` - Convert JSON to Go struct (not yet implemented)
- `--json-file, -f` - JSON input file for --from-json
- `--output, -o` - Output file (default: stdout)
- `--version, -v` - Show version
- `--help` - Show help message

## Development

Run tests:

```bash
go test ./...
```

Build:

```bash
go build ./cmd/structify
```

## License

MIT

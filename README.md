# structify

Structify reads Go structs and generates two things from them: a PostgreSQL schema and ready-to-use `database/sql` CRUD code. You annotate your struct fields with `db` tags to control how the output looks — primary keys, foreign keys, indexes, constraints, and more. No ORM, no reflection at runtime, just generated SQL and Go code you can read and own.

Try it in the browser at [structify.alexander-kruska.dev](https://structify.alexander-kruska.dev).

## What it generates

**SQL mode** (`--to-sql`) produces `CREATE TABLE` statements with all constraints inline:

- Primary keys, including composite primary keys across multiple fields
- `UNIQUE` constraints, individually or grouped by name
- `CHECK` constraints and `DEFAULT` values
- Enum checks using `IN (...)` clauses
- Named and unnamed indexes, including unique indexes
- Foreign keys with optional `ON DELETE` / `ON UPDATE` cascade rules, including composite foreign keys

**Code mode** (`--to-db-sql`) produces Go functions for `Create`, `GetByID`, `Update`, `Delete`, and `List`, using only the standard `database/sql` package. If foreign keys are defined, it also generates `JOIN` queries between related tables.

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

## Struct Tags

| Tag | Effect |
|-----|--------|
| `db:"pk"` | Primary key (use on multiple fields for a composite key) |
| `db:"unique"` | UNIQUE constraint |
| `db:"unique:group_name"` | Composite unique constraint grouped by name |
| `db:"-"` | Exclude field from all output |
| `db:"table:name"` | Override the generated table name (on the struct itself) |
| `db:"check:expr"` | CHECK constraint with the given expression |
| `db:"default:val"` | DEFAULT value |
| `db:"enum:a,b,c"` | CHECK constraint using `IN (a, b, c)` |
| `db:"index"` | Auto-named index |
| `db:"index:idx_name"` | Named index (same name on multiple fields creates a composite index) |
| `db:"unique_index"` | Auto-named unique index |
| `db:"unique_index:uq_name"` | Named unique index |
| `db:"fk:table,col"` | Foreign key referencing `table(col)` |
| `db:"fk:table,col,on_delete:CASCADE"` | Foreign key with cascade option |
| `db:"fk:name,table,col"` | Composite foreign key (same name groups columns together) |

## Usage

Generate a PostgreSQL schema:

```bash
structify --to-sql ./models/user.go
```

Generate `database/sql` CRUD code:

```bash
structify --to-db-sql ./models/user.go
```

Write output to a file:

```bash
structify --to-sql ./models/user.go -o schema.sql
structify --to-db-sql ./models/user.go -o user_repo.go
```

Process multiple files at once:

```bash
structify --to-sql ./models/*.go
```

## Example

Input (`models/user.go`):

```go
package models

type User struct {
    ID       int64  `db:"pk"`
    Username string `db:"unique"`
    Email    string `db:"unique_index:uq_email"`
    Active   bool   `db:"default:true"`
    Age      int    `db:"check:age >= 0"`
}
```

SQL output:

```sql
CREATE TABLE "user" (
    "id" BIGINT PRIMARY KEY,
    "username" VARCHAR(255) UNIQUE,
    "email" VARCHAR(255) NOT NULL,
    "active" BOOLEAN NOT NULL DEFAULT true,
    "age" INTEGER NOT NULL CHECK (age >= 0)
);

CREATE UNIQUE INDEX "uq_email" ON "user" ("email");
```

## Type Mapping

| Go type | PostgreSQL type |
|---------|----------------|
| `int`, `int32` | `INTEGER` |
| `int64` | `BIGINT` |
| `string` | `VARCHAR(255)` |
| `bool` | `BOOLEAN` |
| `float32` | `REAL` |
| `float64` | `DOUBLE PRECISION` |
| `time.Time` | `TIMESTAMP` |
| `[]byte` | `BYTEA` |

## Flags

| Flag | Description |
|------|-------------|
| `--to-sql`, `--to-schema` | Generate PostgreSQL CREATE TABLE statements |
| `--to-db-sql`, `--to-dbcode` | Generate database/sql CRUD code |
| `--output`, `-o` | Write output to file instead of stdout |
| `--version`, `-v` | Print version |
| `--help` | Show help |

## Development

```bash
go test ./...
go build ./cmd/structify
```

## License

MIT

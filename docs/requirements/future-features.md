# Future Features Requirements

## Current State Analysis

### Implemented Features
- Go struct parsing via go/parser
- Basic db tags (pk, unique, -, table:)
- Type mapping for Go to PostgreSQL
- SQL DDL generation (CREATE TABLE)
- database/sql CRUD code generation
- CQRS architecture
- Repository pattern
- ValidatedEntity pattern
- Read-after-write pattern

### Recently Implemented (Phase 1-6, Feb 2025)
✓ **Category 1: Relationship Support - COMPLETE**
  - ✓ `db:"fk:Table,Column"` tag support for single FKs
  - ✓ `db:"fk:constraint_name,table,column"` for composite FKs
  - ✓ FOREIGN KEY constraints in SQL generation
  - ✓ JOIN query methods (single and multi-JOIN)
  - ✓ Cascade options (CASCADE, SET_NULL, RESTRICT, SET_DEFAULT, NO_ACTION)
  - ✓ Composite primary keys (multiple PK fields)
  - ✓ Composite unique constraints
  - ✓ Composite foreign key relationships

✓ **Category 2: Advanced Constraints - COMPLETE**
  - ✓ `db:"check:expression"` tag support (documented and tested)
  - ✓ `db:"default:value"` tag support (documented and tested)
  - ✓ `db:"enum:value1,value2,value3"` tag support
  - ✓ CHECK constraint generation for enums
  - ✓ `db:"index"` tag for single column index
  - ✓ `db:"index:idx_name"` for named indexes
  - ✓ `db:"unique_index"` tag
  - ✓ Composite index support across multiple fields

### Remaining Features

## Category 3: Schema Management

### 3.1 Migration Files
- Generate timestamped migration files
- Track schema version
- Support up/down migrations
- Add `--migrate` flag

### 3.2 Schema Diff
- Compare two schema versions
- Generate ALTER TABLE statements
- Add `--diff` flag with old/new schema files

## Category 4: Code Generation Enhancements

### 4.1 Transaction Support
- Add `*sql.Tx` parameter option to generated methods
- Generate transaction wrapper methods
- Add context with timeout support

### 4.2 Pagination
- Generate `ListWithLimit` and `ListWithOffset` methods
- Add `db:"paginate"` tag for auto-pagination
- Generate cursor-based pagination helpers

### 4.3 Bulk Operations
- Generate `CreateBatch`, `UpdateBatch` methods
- Use PostgreSQL COPY for bulk inserts
- Add `--batch-size` flag

### 4.4 Soft Delete
- Add `db:"soft_delete"` tag
- Auto-filter deleted rows in queries
- Generate `DeletePermanent` method
- Add `deleted_at` timestamp handling

### 4.5 Timestamp Fields
- Auto-generate `created_at`, `updated_at` if tagged
- Update `updated_at` on save
- Add `db:"auto_timestamp"` tag

### 4.6 NULL Handling
- Support `sql.NullString`, `sql.NullInt64` types
- Generate proper Scan/Value methods
- Add `db:"nullable"` tag

## Category 5: JSON Features

### 5.1 JSON to Struct
- Already has `--from-json` flag (not implemented)
- Parse JSON and generate Go structs
- Inverse of current functionality

### 5.2 JSON Columns
- Support `json.RawMessage` fields
- Generate JSONB columns
- Add JSON path query helpers

## Category 6: Database Support

### 6.1 MySQL Support
- Add MySQL dialect
- Type mapping adjustments
- Backtick identifiers instead of quotes
- Add `--dialect=mysql` flag

### 6.2 SQLite Support
- Add SQLite dialect
- Different type system
- Add `--dialect=sqlite` flag

## Category 7: Query Building

### 7.1 Where Clause Builders
- Generate fluent Where methods
- Type-safe query building
- Optional: Add `--query-builder` flag

### 7.2 Dynamic Query Generation
- Generate methods that accept filter structs
- Support partial field updates

## Category 8: Developer Experience

### 8.1 Validation Code
- Generate struct validation methods
- Add `validate:"required"` tag support
- Integrate with popular validators

### 8.2 Hooks/Callbacks
- Add BeforeSave, AfterSave hooks
- Add BeforeDelete, AfterDelete hooks
- Interface-based hook system

### 8.3 Context Support
- All methods accept context.Context ✓ (already implemented)
- Generate timeout variants
- Add cancellation support

### 8.4 Error Wrapping
- Wrap database errors with custom types
- Add constraint violation error types
- Better error messages

## Category 9: Testing Support

### 9.1 Mock Generation
- Generate mock repositories
- Interface-based mocking
- Add `--generate-mocks` flag

### 9.2 Test Fixtures
- Generate test data builders
- Add `--generate-fixtures` flag
- Support for factory pattern

## Category 10: CLI Enhancements

### 10.1 Watch Mode
- Auto-regenerate on file changes
- Add `--watch` flag

### 10.2 Config File
- Support `.structify.yaml` config
- Default flags in config
- Project-specific settings

### 10.3 Stdin Support
- Read struct definitions from stdin
- Pipe support: `cat models.go | structify --to-sql`

## Implementation Priority Questions

### Question 1: Relationship Support
✅ **COMPLETED** - Foreign key relationships with composite key support fully implemented

### Question 2: Database Support
Do you need support for databases other than PostgreSQL?
- A) Yes - MySQL
- B) Yes - SQLite
- C) No - PostgreSQL only is fine

### Question 3: Migration Management
Is automatic migration file generation important?
- A) Yes - want full migration workflow
- B) No - use separate migration tool
- C) Maybe - just need schema diff

### Question 4: JSON Conversion
Is the JSON to Go struct feature needed?
- A) Yes - frequently convert JSON
- B) No - never use it
- C) Sometimes - occasionally useful

### Question 5: Which Category Priority
Which feature category should be implemented next?
- A) Schema Management (Migrations, Diff)
- B) Code Gen (Bulk, Pagination, Soft Delete)
- C) Database Support (MySQL, SQLite)
- D) Developer Experience (Validation, Hooks, Mocks)
- E) CLI Enhancements (Watch, Config)

### Question 6: Breaking Changes
Are you willing to accept breaking changes to existing generated code?
- A) Yes - for better architecture
- B) No - must maintain backward compatibility
- C) Maybe - with version flag

### Question 7: Transaction Pattern
How should transactions be handled in generated code?
- A) Generate separate `*Tx` methods
- B) Add optional transaction parameter
- C) Generate transaction helper functions
- D) Let users handle manually

## Sources

- [GORM Foreign Key Documentation](https://gorm.io/docs/belongs_to.html)
- [Stop Using Repository Pattern with EF Core](https://levelup.gitconnected.com/stop-using-repository-pattern-with-ef-core-in-2026-heres-why-8f22168aba3e)
- [Golang Gen - Code Generation Tool](https://juejin.cn/post/7267927742797004835)
- [Unit of Work with Dapper](https://binodmahto.medium.com/unit-of-work-with-dapper-sql-server-stop-partial-writes-and-make-transactions-boring-907c4ef73755)
- [GORM Relationships](https://www.cnblogs.com/haima/p/14879268.html)

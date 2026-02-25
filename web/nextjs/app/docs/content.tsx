import type { ReactNode } from "react";


export type DocSection = {
  id: string;
  title: string;
  summary: ReactNode;
  sourceCode: string;
  outputCode: string;
  outputLabel: "Output (SQL)" | "Output (Go)";
  sourceLabel?: "Source (Go)" | "Command";
  sourceLanguage?: "go" | "bash";
  tryLink?: {
    label: string;
    href: string;
  };
};

export type DocCategory = {
  id: string;
  title: string;
  intro: ReactNode;
  sections: DocSection[];
};

export const CATEGORIES: DocCategory[] = [
  {
    id: "sql-schema",
    title: "SQL Schema",
    intro: (
      <>
        Define tables with structs and control SQL output with{" "}
        <code className="rounded bg-zinc-800/60 px-1 py-0.5 text-zinc-200">db:&quot;...&quot;</code> tags.
        These sections explain exactly which tags belong where and what SQL is generated.
      </>
    ),
    sections: [
      {
        id: "sql-struct-basics",
        title: "Struct Basics & Table Naming",
        summary: (
          <>
            One struct maps to one table. Struct names are pluralized and converted to snake_case. Field names become columns in snake_case.
            Use <code className="rounded bg-zinc-800/60 px-1 py-0.5 text-zinc-200">db:&quot;pk&quot;</code> for the primary key.
          </>
        ),
        sourceCode: `package models

type BlogPost struct {
\tID        int64  \`db:"pk"\`
\tTitle     string
\tAuthorID  int64
\tCreatedAt int64
}`,
        outputCode: `CREATE TABLE blog_posts (
  id         BIGINT PRIMARY KEY,
  title      TEXT NOT NULL,
  author_id  BIGINT NOT NULL,
  created_at BIGINT NOT NULL
);`,
        outputLabel: "Output (SQL)",
        tryLink: {
          label: "Try SQL Example",
          href: `/?load=${encodeURIComponent("User")}&mode=sql`,
        },
      },
      {
        id: "sql-tags-constraints",
        title: "Tag Patterns: Constraints, Defaults, Enum",
        summary: (
          <>
            Combine tag instructions comma-separated:{" "}
            <code className="rounded bg-zinc-800/60 px-1 py-0.5 text-zinc-200">db:&quot;check:...,default:...,enum:...&quot;</code>.
            Keep each rule explicit so generated SQL stays predictable.
          </>
        ),
        sourceCode: `package models

type Account struct {
\tID       int64  \`db:"pk"\`
\tEmail    string \`db:"unique,check:length(email) > 5"\`
\tStatus   string \`db:"enum:active,disabled,pending,default:pending"\`
\tAge      int    \`db:"check:age >= 18,default:18"\`
\tVerified bool   \`db:"default:false"\`
}`,
        outputCode: `CREATE TABLE accounts (
  id       BIGINT PRIMARY KEY,
  email    TEXT NOT NULL UNIQUE CHECK (length(email) > 5),
  status   TEXT NOT NULL DEFAULT pending CHECK (status IN ('active', 'disabled', 'pending')),
  age      INTEGER NOT NULL DEFAULT 18 CHECK (age >= 18),
  verified BOOLEAN NOT NULL DEFAULT false
);`,
        outputLabel: "Output (SQL)",
        tryLink: {
          label: "Try Constraints",
          href: `/?load=${encodeURIComponent("Constraints")}&mode=sql`,
        },
      },
      {
        id: "sql-indexes-relations",
        title: "Indexes & Foreign Keys",
        summary: (
          <>
            Use <code className="rounded bg-zinc-800/60 px-1 py-0.5 text-zinc-200">index</code>, <code className="rounded bg-zinc-800/60 px-1 py-0.5 text-zinc-200">unique_index:&lt;name&gt;</code>,
            and <code className="rounded bg-zinc-800/60 px-1 py-0.5 text-zinc-200">fk:&lt;table&gt;,&lt;column&gt;,on_delete:CASCADE</code>.
            Share the same index name across multiple fields for composite indexes.
          </>
        ),
        sourceCode: `package models

type Article struct {
\tID         int64  \`db:"pk"\`
\tSlug       string \`db:"unique_index:uq_articles_slug"\`
\tLanguage   string \`db:"index:idx_articles_title_lang"\`
\tTitle      string \`db:"index:idx_articles_title_lang"\`
\tCategoryID int64  \`db:"fk:categories,id,on_delete:RESTRICT"\`
}`,
        outputCode: `CREATE TABLE articles (
  id          BIGINT PRIMARY KEY,
  slug        TEXT NOT NULL,
  language    TEXT NOT NULL,
  title       TEXT NOT NULL,
  category_id BIGINT NOT NULL,
  CONSTRAINT fk_articles_category_id FOREIGN KEY (category_id)
    REFERENCES categories (id) ON DELETE RESTRICT
);

CREATE UNIQUE INDEX uq_articles_slug ON articles (slug);
CREATE INDEX idx_articles_title_lang ON articles (language, title);`,
        outputLabel: "Output (SQL)",
        tryLink: {
          label: "Try Indexes",
          href: `/?load=${encodeURIComponent("Indexes")}&mode=sql`,
        },
      },
    ],
  },
  {
    id: "repository",
    title: "Repository",
    intro: (
      <>
        Generate implementation code from interfaces with <code className="rounded bg-zinc-800/60 px-1 py-0.5 text-zinc-200">--to-repo</code>.
        Naming conventions and method signatures decide which queries are generated.
      </>
    ),
    sections: [
      {
        id: "repo-prerequisites",
        title: "Pre-requisites Checklist",
        summary: (
          <>
            Before generation: model has primary key, field types are supported, and interface methods use <code className="rounded bg-zinc-800/60 px-1 py-0.5 text-zinc-200">context.Context</code> as first argument.
            Keep signatures deterministic so generated code remains stable.
          </>
        ),
        sourceCode: `package models

import "context"

type User struct {
\tID    int64  \`db:"pk"\`
\tEmail string \`db:"unique"\`
}

type UserRepository interface {
\tFindByID(ctx context.Context, id int64) (*User, error)
\tCreate(ctx context.Context, user *User) error
}`,
        outputCode: `Checklist:
- model contains db:"pk"
- interface name follows <Model>Repository
- methods return (..., error) or error
- parameter order is stable
- find methods use explicit filter names`,
        outputLabel: "Output (Go)",
      },
      {
        id: "repo-interface-conventions",
        title: "Interface Design Rules",
        summary: (
          <>
            Use <code className="rounded bg-zinc-800/60 px-1 py-0.5 text-zinc-200">&lt;Model&gt;Repository</code> as interface name.
            Prefer clear method verbs like <code className="rounded bg-zinc-800/60 px-1 py-0.5 text-zinc-200">FindBy</code>, <code className="rounded bg-zinc-800/60 px-1 py-0.5 text-zinc-200">Create</code>, <code className="rounded bg-zinc-800/60 px-1 py-0.5 text-zinc-200">Update</code>, and <code className="rounded bg-zinc-800/60 px-1 py-0.5 text-zinc-200">DeleteBy</code>.
          </>
        ),
        sourceCode: `type UserRepository interface {
\tFindByID(ctx context.Context, id int64) (*User, error)
\tFindByEmail(ctx context.Context, email string) (*User, error)
\tFindByStatus(ctx context.Context, status string) ([]*User, error)
\tCreate(ctx context.Context, user *User) error
\tUpdate(ctx context.Context, user *User) error
\tDeleteByID(ctx context.Context, id int64) error
}`,
        outputCode: `Supported return styles:
- (*User, error)
- ([]*User, error)
- error

Avoid:
- methods without error return
- unclear names like GetStuff(ctx, x)
- context not passed as first parameter`,
        outputLabel: "Output (Go)",
        tryLink: {
          label: "Try Repo Input",
          href: `/?load=${encodeURIComponent("Indexes")}&mode=code`,
        },
      },
      {
        id: "repo-method-derivation",
        title: "Method Derivation & Query Mapping",
        summary: (
          <>
            Query filters are derived from method names. Example: <code className="rounded bg-zinc-800/60 px-1 py-0.5 text-zinc-200">FindByTitleAndLanguage</code>
            maps to a WHERE clause with both columns.
          </>
        ),
        sourceCode: `type ArticleRepository interface {
\tFindByID(ctx context.Context, id int64) (*Article, error)
\tFindByTitleAndLanguage(ctx context.Context, title string, language string) ([]*Article, error)
\tDeleteByID(ctx context.Context, id int64) error
}`,
        outputCode: `func (r *articleRepository) FindByTitleAndLanguage(ctx context.Context, title string, language string) ([]*models.Article, error) {
\tquery := "SELECT id, slug, title, language FROM articles WHERE title = $1 AND language = $2"
\trows, err := r.db.QueryContext(ctx, query, title, language)
\tif err != nil {
\t\treturn nil, err
\t}
\tdefer rows.Close()

\tvar results []*models.Article
\tfor rows.Next() {
\t\tvar item models.Article
\t\tif err := rows.Scan(&item.ID, &item.Slug, &item.Title, &item.Language); err != nil {
\t\t\treturn nil, err
\t\t}
\t\tresults = append(results, &item)
\t}
\treturn results, rows.Err()
}`,
        outputLabel: "Output (Go)",
      },
      {
        id: "repo-generation-flow",
        title: "End-to-End Generation Flow",
        summary: (
          <>
            Keep models and interfaces in the same domain package, then run generation into your repository package.
            Generated files should stay checked in so behavior is reviewable in PRs.
          </>
        ),
        sourceCode: `# 1) Define model + repository interface
# 2) Generate implementation
structify --input ./internal/domain --output ./internal/repository --to-repo

# 3) Verify generated files
ls ./internal/repository/*.gen.go`,
        sourceLabel: "Command",
        sourceLanguage: "bash",
        outputCode: `Expected output files:
- user_repository.gen.go
- article_repository.gen.go

Each file contains:
- concrete repository struct
- constructor
- generated SQL methods`,
        outputLabel: "Output (Go)",
        tryLink: {
          label: "Try Code Gen",
          href: `/?load=${encodeURIComponent("Indexes")}&mode=code`,
        },
      },
      {
        id: "repo-wiring-services",
        title: "Wiring Generated Repositories",
        summary: (
          <>
            Keep services depending on interfaces, not concrete repository structs. This makes replacement and testing straightforward.
          </>
        ),
        sourceCode: `package service

type UserService struct {
\trepo models.UserRepository
}

func NewUserService(repo models.UserRepository) *UserService {
\treturn &UserService{repo: repo}
}

func Wire(db *sql.DB) *UserService {
\trepo := repository.NewUserRepository(db)
\treturn NewUserService(repo)
}`,
        outputCode: `Benefits:
- service layer remains testable
- generated implementation can be replaced by mocks
- migration to custom repositories stays simple`,
        outputLabel: "Output (Go)",
      },
      {
        id: "repo-errors-transactions",
        title: "Error Handling & Transactions",
        summary: (
          <>
            Generated methods return database errors directly. Wrap errors at service boundaries with operation context.
            For multi-step writes, use explicit transaction orchestration in service layer.
          </>
        ),
        sourceCode: `func (s *OrderService) Place(ctx context.Context, input PlaceOrderInput) error {
\ttx, err := s.db.BeginTx(ctx, nil)
\tif err != nil {
\t\treturn fmt.Errorf("begin tx: %w", err)
\t}
\tdefer tx.Rollback()

\tif err := s.orderRepo.CreateTx(ctx, tx, input.Order); err != nil {
\t\treturn fmt.Errorf("create order: %w", err)
\t}
\tif err := s.stockRepo.ReserveTx(ctx, tx, input.Items); err != nil {
\t\treturn fmt.Errorf("reserve stock: %w", err)
\t}
	return tx.Commit()
}`,
        outputCode: `Guidance:
- do not swallow generated errors
- wrap with %w in service layer
- use transactions for write chains
- keep retry logic outside repository`,
        outputLabel: "Output (Go)",
      },
      {
        id: "repo-testing-generated",
        title: "Testing Generated Repositories",
        summary: (
          <>
            Test both behavior and SQL assumptions. Unit tests can mock DB adapters, integration tests should run against real PostgreSQL.
          </>
        ),
        sourceCode: `# unit tests
go test ./internal/repository/...

# integration tests
go test -tags=integration ./test/integration/...`,
        sourceLabel: "Command",
        sourceLanguage: "bash",
        outputCode: `Suggested cases:
- FindBy returns not found / one row / many rows
- Create handles unique violation
- Update validates affected rows
- DeleteByID is idempotent where needed`,
        outputLabel: "Output (Go)",
      },
      {
        id: "repo-common-pitfalls",
        title: "Common Pitfalls & Fixes",
        summary: (
          <>
            Most generation issues come from ambiguous method names or mismatched field names.
            Keep names explicit and aligned with struct field intent.
          </>
        ),
        sourceCode: `Bad:
Find(ctx context.Context, a string, b string) (*User, error)

Good:
FindByEmailAndStatus(ctx context.Context, email string, status string) (*User, error)

Bad:
GetByUser(ctx context.Context, id int64) (*User, error)

Good:
FindByID(ctx context.Context, id int64) (*User, error)`,
        sourceLabel: "Command",
        sourceLanguage: "bash",
        outputCode: `Fix checklist:
- rename methods to FindBy.../DeleteBy...
- keep parameter names aligned with fields
- avoid overloaded meaning in one method
- regenerate and diff *.gen.go files`,
        outputLabel: "Output (Go)",
      },
    ],
  },
  {
    id: "project-guide",
    title: "Project Guide",
    intro: (
      <>
        Useful project-level references for local development, architecture orientation, and troubleshooting.
      </>
    ),
    sections: [
      {
        id: "project-feature-overview",
        title: "Feature Overview",
        summary: (
          <>
            structify supports two main outputs: SQL schema generation from structs and repository implementation generation from interfaces.
            Both are designed for predictable diffs and reviewable output.
          </>
        ),
        sourceCode: `Inputs:
- Go structs with db tags
- repository interfaces with naming conventions

Outputs:
- SQL CREATE TABLE + INDEX + FK statements
- generated *.gen.go repositories`,
        sourceLabel: "Command",
        sourceLanguage: "bash",
        outputCode: `Use SQL mode for:
- schema planning
- migration baselines

Use repo mode for:
- boilerplate reduction
- consistent query mapping`,
        outputLabel: "Output (Go)",
      },
      {
        id: "project-local-commands",
        title: "Local Commands",
        summary: (
          <>
            Run backend and frontend independently while iterating on generation and docs.
          </>
        ),
        sourceCode: `go run ./cmd/server

go test ./...

go test -tags=integration ./test/integration/...

npm --prefix web/nextjs run dev
npm --prefix web/nextjs run test
npm --prefix web/nextjs run lint`,
        sourceLabel: "Command",
        sourceLanguage: "bash",
        outputCode: `Tip:
- keep backend server running for UI API calls
- run focused frontend tests during docs updates
- run full test suites before PR`,
        outputLabel: "Output (Go)",
      },
      {
        id: "project-structure-reference",
        title: "Recommended Project Structure",
        summary: (
          <>
            Follow the existing repository layout to keep generated output in predictable paths.
          </>
        ),
        sourceCode: `cmd/           # entrypoints
internal/      # core implementation
pkg/           # public packages
api/           # OpenAPI/proto definitions
migrations/    # DB migrations
web/nextjs/    # frontend`,
        sourceLabel: "Command",
        sourceLanguage: "bash",
        outputCode: `Repository generation suggestion:
- interfaces near domain models
- generated repos in internal/repository
- service layer depends on interfaces`,
        outputLabel: "Output (Go)",
      },
      {
        id: "project-troubleshooting",
        title: "Troubleshooting",
        summary: (
          <>
            If output is unexpected, inspect names first: struct fields, tags, interface methods, and parameters.
            Most issues are deterministic and fixed by naming alignment.
          </>
        ),
        sourceCode: `Symptom: method not generated as expected
Check: method verb + FindBy naming + parameter order

Symptom: wrong SQL column
Check: field name and db tags

Symptom: failing test after regeneration
Check: golden output and update expected fixtures intentionally`,
        sourceLabel: "Command",
        sourceLanguage: "bash",
        outputCode: `Quick debug loop:
1) minimal reproduce model/interface
2) regenerate output
3) diff generated file
4) fix naming/tag issue
5) rerun tests`,
        outputLabel: "Output (Go)",
      },
    ],
  },
];

export function getCategory(categoryId: DocCategory["id"]): DocCategory | undefined {
  return CATEGORIES.find((category) => category.id === categoryId);
}

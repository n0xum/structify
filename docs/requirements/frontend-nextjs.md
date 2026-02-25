# Requirements: Next.js Frontend

**Feature:** Web UI for structify
**Status:** Draft
**Date:** 2026-02-13

---

## Decision: Monorepo with /web subfolder

The frontend lives in `web/` inside the existing structify repository.

Reasons:
- Frontend and backend share the same API contract — a single PR can update both sides atomically
- One CI pipeline coordinates the Hetzner backend deploy and the GitHub Pages frontend deploy
- Dependabot, SonarCloud and existing tooling cover both without extra configuration
- A separate repo is only justified when independent teams deploy independently

The Next.js app is exported as a fully static site (`output: 'export'`). It calls the Go HTTP
API running on the Hetzner VPS. GitHub Pages serves the static build output from the `web/out/`
directory via a dedicated GitHub Actions workflow.

```
[Browser]
   |
   |-- static assets --> [GitHub Pages]  (web/out/)
   |
   |-- POST /api/generate/sql   --> [Hetzner VPS]  (cmd/server/)
   |-- POST /api/generate/code  -->     Go HTTP server
```

---

## Overview

A modern, minimal web frontend that exposes the full structify feature set in a browser.
Users paste Go struct source code, select an output mode and receive generated SQL or
repository code immediately — without installing the CLI.

The UI must communicate clearly what the tool does and how to use it within seconds of
first visit. The primary user is a Go developer who already knows what a struct is.

---

## Functional Requirements

### FR-001 — Struct Input Editor
**Priority: HIGH**

A code editor area accepts raw Go struct source code as input.
- Syntax highlighting for Go
- Accepts multiple structs in a single input
- Placeholder text shows a minimal example struct on first load
- Input persists in `localStorage` across page reloads so the user does not lose work
- Maximum input size enforced client-side at 100 KB with a visible warning

### FR-002 — Output Mode Selection
**Priority: HIGH**

Two clearly labeled tabs or toggle buttons select the output mode:
- **SQL Schema** — generates PostgreSQL `CREATE TABLE` statements
- **Repository Code** — generates `database/sql` CRUD functions

The active mode is stored in the URL as a query param (`?mode=sql` / `?mode=code`)
so the page is shareable and bookmarkable.

### FR-003 — Generate Action
**Priority: HIGH**

A "Generate" button submits the struct input and selected mode to the backend.
- Button is disabled and shows a loading indicator while the request is in flight
- On success: output is rendered in the output area
- On error: a human-readable error message appears inline below the input — editor is not cleared
- Keyboard shortcut `Ctrl+Enter` / `Cmd+Enter` triggers generation
- Empty input disables the button

### FR-004 — Output Editor
**Priority: HIGH**

A read-only code editor displays the generation result.
- Syntax highlighting: SQL for schema mode, Go for code mode
- Line numbers visible
- "Copy to clipboard" button in the top-right corner of the output area
- Brief "Copied!" confirmation shown for 2 seconds after copying
- Output area is empty and shows helper text ("Output will appear here") before first generation

### FR-005 — Split View Layout
**Priority: HIGH**

On desktop (>= 768 px): input and output are displayed side by side at equal width
with a fixed header above.
On mobile (< 768 px): input and output stack vertically; output scrolls into view
automatically after successful generation.

### FR-006 — Live Client-Side Validation
**Priority: MEDIUM**

Before any backend call, the frontend validates the input and shows inline warnings:
- Input is empty
- No exported struct found (no uppercase identifier followed by `struct`)
- No fields detected
Warnings appear below the input editor and do not block the Generate button.
They disappear as soon as the issue is resolved.

### FR-007 — Example Loader
**Priority: MEDIUM**

A dropdown in the header offers built-in example structs:
- User (pk, unique, basic types)
- Product (pk, ignored field)
- OrderItem (multiple foreign key fields, float)

Selecting an example populates the input editor and triggers generation automatically.
A reset button (`×`) restores the editor to its previous content.

### FR-008 — Package Name Input
**Priority: MEDIUM**

A small text input below the mode selector lets the user set the Go package name
used in the generated repository code.
- Only visible when "Repository Code" mode is active
- Defaults to `models`
- Validated client-side: only lowercase letters, digits and underscores

### FR-009 — Hero Section / Landing Area
**Priority: MEDIUM**

Above the editor, a compact hero section communicates what structify does:
- One-line headline: "Go structs to PostgreSQL, instantly."
- Two-line description explaining input → output
- Link to the GitHub repository
- Version badge showing the current backend version (fetched from `GET /api/version`)

### FR-010 — Tag Reference Panel
**Priority: LOW**

A collapsible side panel or tooltip overlay shows the supported `db` tag reference:
`pk`, `unique`, `-`, `table:name` with a one-line description each.
Triggered by a "?" button near the input editor label.

### FR-011 — Error Boundary
**Priority: LOW**

If the page crashes client-side (React error boundary), a fallback UI is shown
with a link to the GitHub issues page and the raw error message in a `<details>` element.

---

## Non-Functional Requirements

### NFR-001 — Performance
- Largest Contentful Paint (LCP) < 1.5 s on a standard broadband connection
- Generation round-trip < 500 ms for inputs up to 20 structs (excluding network latency to Hetzner)
- No full-page reloads — all state transitions happen client-side

### NFR-002 — Accessibility
- WCAG 2.1 AA compliance
- All interactive elements reachable and operable by keyboard alone
- Focus indicators visible on all focusable elements
- Error messages associated with their input via `aria-describedby`
- Code editor areas have `role="region"` and accessible labels

### NFR-003 — Responsive Design
- Fully usable on viewports from 375 px (iPhone SE) to 1920 px
- No horizontal overflow at any supported viewport width
- Touch targets minimum 44 × 44 px on mobile

### NFR-004 — No Authentication Required
The tool is public and stateless. No user accounts, sessions or server-side storage.

### NFR-005 — Security
- All communication with the Hetzner API over HTTPS
- The Go API server validates input size (max 500 KB) and content type
- No user input is logged or persisted server-side
- CORS on the Go server restricted to the GitHub Pages domain
- Rate limiting on the Go server: 30 requests per minute per IP

### NFR-006 — Dark Mode
Respects system preference via Tailwind `dark:` classes.
No manual toggle required (can be added later).

---

## Technical Constraints

### TC-001 — Frontend Stack
- Next.js 15 with `output: 'export'` (fully static, no Node.js server at runtime)
- React 19
- TypeScript (strict mode)
- Tailwind CSS v4 — no additional component library
- No Server Actions (incompatible with static export)

### TC-002 — Code Editor
- `@uiw/react-codemirror` with `@codemirror/lang-go` and `@codemirror/lang-sql`
- Dynamically imported to keep initial bundle small
- Theme: `vscodeDark` for both editor and output

### TC-003 — API Client
All backend calls made from a single `lib/api.ts` module using the native `fetch` API.
The base URL is read from `NEXT_PUBLIC_API_URL` environment variable.
For local development this points to `http://localhost:8080`.
For production it points to the Hetzner VPS domain.

```typescript
// lib/api.ts
export async function generateSQL(source: string): Promise<string>
export async function generateCode(source: string, pkg: string): Promise<string>
```

### TC-004 — Go HTTP Server (`cmd/server/`)
New entrypoint exposing three endpoints:

```
POST /api/generate/sql
     Body: { "source": string }
     Response: { "output": string } | { "error": string }

POST /api/generate/code
     Body: { "source": string, "package": string }
     Response: { "output": string } | { "error": string }

GET  /api/version
     Response: { "version": string }
```

Handlers reuse the existing application layer unchanged:
`query.Handler.Parse()` → `command.Handler.GenerateSchema()` / `GenerateCode()`

Uses only `net/http` from the standard library. No new Go dependencies.

### TC-005 — Repository Layout

```
structify/
├── cmd/
│   ├── structify/        # existing CLI
│   └── server/           # new Go HTTP server
│       └── main.go
├── web/                  # Next.js application
│   ├── app/
│   │   ├── layout.tsx
│   │   ├── page.tsx      # main editor page
│   │   └── globals.css
│   ├── components/
│   │   ├── Editor.tsx        # CodeMirror wrapper
│   │   ├── OutputPanel.tsx
│   │   ├── ModeSelector.tsx
│   │   ├── ExampleLoader.tsx
│   │   └── TagReference.tsx
│   ├── lib/
│   │   ├── api.ts            # fetch wrapper
│   │   └── validation.ts     # client-side struct checks
│   ├── next.config.ts
│   ├── package.json
│   └── tsconfig.json
└── .github/
    └── workflows/
        ├── ci.yml            # existing
        ├── deploy-backend.yml  # new: SSH deploy to Hetzner on main push
        └── deploy-frontend.yml # new: build web/ and push to GitHub Pages
```

### TC-006 — GitHub Pages Deployment
A dedicated workflow `deploy-frontend.yml` triggers on push to `main`:
1. `npm ci` in `web/`
2. `npm run build` (produces `web/out/`)
3. `actions/upload-pages-artifact` from `web/out/`
4. `actions/deploy-pages`

Requires `NEXT_PUBLIC_API_URL` set as a GitHub Actions variable (not secret).

### TC-007 — Backend Deployment
A dedicated workflow `deploy-backend.yml` triggers on push to `main`:
1. Cross-compile: `GOOS=linux GOARCH=amd64 go build ./cmd/server/`
2. `rsync` or `scp` the binary to the Hetzner VPS
3. SSH restart of the systemd service

Requires secrets: `HETZNER_HOST`, `HETZNER_USER`, `HETZNER_SSH_KEY`.

---

## Acceptance Criteria

**AC-001 — SQL generation end-to-end**
```
Given a valid Go struct with db tags in the input editor
When the user selects "SQL Schema" and clicks "Generate"
Then the output area displays a valid PostgreSQL CREATE TABLE statement
 And the output contains the correct column types and constraints
```

**AC-002 — Repository code generation end-to-end**
```
Given a valid Go struct in the input editor
When the user selects "Repository Code" and clicks "Generate"
Then the output area displays compilable database/sql CRUD functions
 And the output contains Create, Get, Update, Delete and List functions
 And the output imports "database/sql" and "context"
```

**AC-003 — Error handling**
```
Given a malformed Go struct in the input editor
When the user clicks "Generate"
Then an error message is shown below the input
 And the input editor retains its current content
 And the output area retains any previously generated output
```

**AC-004 — Keyboard shortcut**
```
Given the input editor has content
When the user presses Ctrl+Enter (or Cmd+Enter on macOS)
Then generation is triggered as if the Generate button was clicked
```

**AC-005 — Example loader**
```
Given the page is freshly loaded
When the user selects "User" from the examples dropdown
Then the input editor is populated with the User struct
 And generation is triggered automatically
 And the output area shows the result for the User struct
```

**AC-006 — Input persistence**
```
Given the user has typed a struct into the editor
When the user refreshes the page
Then the struct is still present in the editor
```

**AC-007 — Mode in URL**
```
Given the user selects "Repository Code" mode
When the user copies and opens the URL in a new tab
Then "Repository Code" mode is pre-selected
```

**AC-008 — Mobile layout**
```
Given a viewport width of 375 px
When generation completes
Then the output area scrolls into view automatically
 And both editor and output are fully readable without horizontal scrolling
```

---

## Dependencies

| Dependency | Type | Notes |
|------------|------|-------|
| `cmd/server/` | New Go binary | Wraps existing application layer |
| `internal/application`, `internal/application/query`, `internal/application/command` | Existing | No changes required |
| `@uiw/react-codemirror` | NPM | Code editor |
| `@codemirror/lang-go` | NPM | Go syntax highlighting |
| `@codemirror/lang-sql` | NPM | SQL syntax highlighting |
| `tailwindcss` v4 | NPM | Styling |
| Hetzner VPS | Infrastructure | SSH access, systemd service for Go server |
| GitHub Pages | Hosting | Enabled in repository settings |
| `NEXT_PUBLIC_API_URL` | GitHub Actions variable | Points to Hetzner VPS |
| `HETZNER_HOST`, `HETZNER_USER`, `HETZNER_SSH_KEY` | GitHub Actions secrets | Backend deployment |

---

## Out of Scope (v0.2.0+)

- Shareable permalink to a generation result (requires server-side storage)
- File upload for `.go` files
- Manual dark/light theme toggle
- MySQL / SQLite output mode
- Diff view comparing multiple generation runs

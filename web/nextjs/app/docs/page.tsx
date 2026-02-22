import Link from "next/link";
import React from "react";
import { HighlightedCode } from "@/components/HighlightedCode";

export const metadata = {
  title: "Docs — structify",
  description: "Learn how to use structify to generate database schemas and repository implementations from Go structs and interfaces.",
};

type Section = {
  id: string;
  title: string;
  exampleLabel: string;
  description: React.ReactNode;
  goCode: string;
  sqlCode: string;
};

const SECTIONS: Section[] = [
  {
    id: "basic",
    title: "Basic Mapping",
    exampleLabel: "User",
    description: (
      <>
        Define a table by declaring a Go struct. Fields are mapped to columns using their Go type. Use <code className="text-zinc-200 bg-zinc-800/50 px-1 py-0.5 rounded">db:&quot;pk&quot;</code> for the primary key and <code className="text-zinc-200 bg-zinc-800/50 px-1 py-0.5 rounded">db:&quot;unique&quot;</code> for a UNIQUE constraint. Fields with no tag are plain nullable columns.
      </>
    ),
    goCode: `package models

type User struct {
	ID       int64  \`db:"pk"\`
	Username string \`db:"unique"\`
	Email    string
	Active   bool
	Created  int64
}`,
    sqlCode: `CREATE TABLE users (
  id       BIGINT PRIMARY KEY,
  username TEXT NOT NULL UNIQUE,
  email    TEXT NOT NULL,
  active   BOOLEAN NOT NULL,
  created  BIGINT NOT NULL
);`,
  },
  {
    id: "constraints",
    title: "Constraints",
    exampleLabel: "Constraints",
    description: (
      <>
        Add column-level constraints inline. <code className="text-zinc-200 bg-zinc-800/50 px-1 py-0.5 rounded">db:&quot;check:&lt;expr&gt;&quot;</code> emits a CHECK constraint, <code className="text-zinc-200 bg-zinc-800/50 px-1 py-0.5 rounded">db:&quot;default:&lt;value&gt;&quot;</code> sets a DEFAULT, and <code className="text-zinc-200 bg-zinc-800/50 px-1 py-0.5 rounded">db:&quot;enum:&lt;a&gt;,&lt;b&gt;,...&quot;</code> generates an IN-list CHECK for allowed values.
      </>
    ),
    goCode: `package models

type Person struct {
	ID        int64  \`db:"pk"\`
	Name      string \`db:"check:length(name) > 0"\`
	Age       int    \`db:"check:age >= 18,default:18"\`
	Status    string \`db:"enum:active,inactive,banned"\`
	Active    bool   \`db:"default:true"\`
	CreatedAt int64  \`db:"default:now()"\`
}`,
    sqlCode: `CREATE TABLE persons (
  id         BIGINT PRIMARY KEY,
  name       TEXT NOT NULL CHECK (length(name) > 0),
  age        INTEGER NOT NULL DEFAULT 18 CHECK (age >= 18),
  status     TEXT NOT NULL CHECK (status IN ('active', 'inactive', 'banned')),
  active     BOOLEAN NOT NULL DEFAULT true,
  created_at BIGINT NOT NULL DEFAULT now()
);`,
  },
  {
    id: "indexes",
    title: "Indexes",
    exampleLabel: "Indexes",
    description: (
      <>
        Create indexes with <code className="text-zinc-200 bg-zinc-800/50 px-1 py-0.5 rounded">db:&quot;index&quot;</code> or <code className="text-zinc-200 bg-zinc-800/50 px-1 py-0.5 rounded">db:&quot;unique_index&quot;</code>. Give multiple fields the same named index to produce a composite index. Unnamed indexes get an auto-generated name.
      </>
    ),
    goCode: `package models

type Article struct {
	ID       int64  \`db:"pk"\`
	Slug     string \`db:"unique_index:uq_slug"\`
	Title    string \`db:"index:idx_title_lang"\`
	Language string \`db:"index:idx_title_lang"\`
	Views    int64  \`db:"index"\`
}`,
    sqlCode: `CREATE TABLE articles (
  id       BIGINT PRIMARY KEY,
  slug     TEXT NOT NULL,
  title    TEXT NOT NULL,
  language TEXT NOT NULL,
  views    BIGINT NOT NULL
);

CREATE UNIQUE INDEX uq_slug ON articles (slug);
CREATE INDEX idx_title_lang ON articles (title, language);
CREATE INDEX idx_articles_views ON articles (views);`,
  },
  {
    id: "foreign-keys",
    title: "Foreign Keys",
    exampleLabel: "Foreign Keys",
    description: (
      <>
        Reference another table with <code className="text-zinc-200 bg-zinc-800/50 px-1 py-0.5 rounded">db:&quot;fk:&lt;table&gt;,&lt;col&gt;&quot;</code>. Add <code className="text-zinc-200 bg-zinc-800/50 px-1 py-0.5 rounded">on_delete:CASCADE</code> (or SET NULL / RESTRICT) to control cascading behaviour. The constraint name is generated automatically.
      </>
    ),
    goCode: `package models

type User struct {
	ID    int64  \`db:"pk"\`
	Email string \`db:"unique"\`
}

type Post struct {
	ID      int64  \`db:"pk"\`
	UserID  int64  \`db:"fk:users,id,on_delete:CASCADE"\`
	Title   string
	Content string
}`,
    sqlCode: `CREATE TABLE users (
  id    BIGINT PRIMARY KEY,
  email TEXT NOT NULL UNIQUE
);

CREATE TABLE posts (
  id      BIGINT PRIMARY KEY,
  user_id BIGINT NOT NULL,
  title   TEXT NOT NULL,
  content TEXT NOT NULL,
  CONSTRAINT fk_posts_user_id FOREIGN KEY (user_id)
    REFERENCES users (id) ON DELETE CASCADE
);`,
  },
  {
    id: "repo-generation",
    title: "Repository Generation",
    exampleLabel: "Indexes",
    description: (
      <>
        With the new <code className="text-zinc-200 bg-zinc-800/50 px-1 py-0.5 rounded text-sm tracking-wide">--to-repo</code> command, structify can generate clean, robust Go code that implements a <code className="text-zinc-200 font-medium">Repository Interface</code>. Just provide your model and its corresponding interface, and structify handles the SQL generation and mapping.
      </>
    ),
    goCode: `package models

import "context"

type Article struct {
	ID       int64  \`db:"pk"\`
	Slug     string \`db:"unique"\`
	Title    string
	Language string
}

// Generate implementation with structify --to-repo
type ArticleRepository interface {
	FindByID(ctx context.Context, id int64) (*Article, error)
	FindBySlug(ctx context.Context, slug string) (*Article, error)
	FindByTitleAndLanguage(ctx context.Context, title string, language string) ([]*Article, error)
	Create(ctx context.Context, article *Article) error
}`,
    sqlCode: `// Generated snippet
func (r *articleRepository) FindByTitleAndLanguage(ctx context.Context, title string, language string) ([]*models.Article, error) {
	query := "SELECT id, slug, title, language FROM articles WHERE title = $1 AND language = $2"
	rows, err := r.db.QueryContext(ctx, query, title, language)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*models.Article
	for rows.Next() {
		var item models.Article
		if err := rows.Scan(&item.ID, &item.Slug, &item.Title, &item.Language); err != nil {
			return nil, err
		}
		results = append(results, &item)
	}
	return results, rows.Err()
}`,
  },
];

const NAV_LINKS = SECTIONS.map((s) => ({ id: s.id, title: s.title }));

function CodeBlock({ label, code, isGo }: { label: string; code: string; isGo?: boolean }) {
  return (
    <div className="flex flex-col h-full rounded-2xl overflow-hidden border border-zinc-800 bg-zinc-900/50 backdrop-blur-xl shadow-2xl transition-all duration-500 hover:border-zinc-700">
      <div className="flex items-center px-4 py-3 border-b border-zinc-800 bg-zinc-800/20">
        <div className="flex gap-1.5 mr-4 opacity-50 grayscale">
          <div className="w-3 h-3 rounded-full bg-zinc-500 border border-zinc-400/50" />
          <div className="w-3 h-3 rounded-full bg-zinc-500 border border-zinc-400/50" />
          <div className="w-3 h-3 rounded-full bg-zinc-500 border border-zinc-400/50" />
        </div>
        <span className="text-xs font-medium uppercase tracking-widest text-zinc-400">
          {label}
        </span>
      </div>
      <div className="p-4 overflow-x-auto flex-1 text-sm font-mono leading-relaxed">
        <HighlightedCode code={code} language={isGo ? "go" : "sql"} />
      </div>
    </div>
  );
}

export default function DocsPage() {
  return (
    <div className="min-h-screen bg-[#050505] text-zinc-100 selection:bg-zinc-500/30 font-sans relative">
      {/* Background ambient lights (now grayscale) */}
      <div className="fixed inset-0 overflow-hidden pointer-events-none">
        <div className="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] rounded-full bg-zinc-800/10 blur-[120px]" />
        <div className="absolute bottom-[-10%] right-[-10%] w-[40%] h-[40%] rounded-full bg-zinc-800/10 blur-[120px]" />
      </div>

      <div className="flex max-w-[1400px] mx-auto min-h-screen relative z-10">
        {/* Sidebar */}
        <aside className="hidden lg:flex flex-col sticky top-0 h-screen w-64 border-r border-white/5 py-12 px-8 shrink-0 bg-black/20 backdrop-blur-3xl">
          <Link
            href="/"
            className="group flex items-center gap-2 text-sm font-medium text-zinc-400 hover:text-white transition-colors mb-12"
          >
            <span className="group-hover:-translate-x-1 transition-transform inline-block">&larr;</span> Back to Editor
          </Link>

          <div className="mb-4">
            <div className="text-xs font-bold text-zinc-300 uppercase tracking-[0.2em]">
              Documentation
            </div>
          </div>

          <nav className="flex flex-col gap-1">
            {NAV_LINKS.map((link) => (
              <a
                key={link.id}
                href={`#${link.id}`}
                className="text-sm text-zinc-500 hover:text-zinc-200 transition-all duration-300 py-2 rounded-lg hover:bg-white/5 hover:px-3 focus-visible:outline focus-visible:outline-2 focus-visible:outline-sky-500"
              >
                {link.title}
              </a>
            ))}
          </nav>
        </aside>

        {/* Main Content */}
        <main className="flex-1 min-w-0 px-6 py-12 lg:py-20 lg:px-20">
          <div className="max-w-4xl mx-auto">
            {/* Mobile Nav */}
            <div className="lg:hidden mb-10">
              <Link
                href="/"
                className="inline-flex items-center gap-2 text-sm font-medium text-zinc-400 hover:text-white transition-colors bg-white/5 px-4 py-2 rounded-full border border-white/5 backdrop-blur-md"
              >
                &larr; Back to Editor
              </Link>
            </div>

            <header className="mb-20">
              <h1 className="text-5xl font-extrabold tracking-tight mb-6 bg-clip-text text-transparent bg-gradient-to-b from-white to-white/60">
                Docs & Features
              </h1>
              <p className="text-lg text-zinc-400 leading-relaxed max-w-2xl">
                Explore how structify transforms simple Go structs into powerful PostgreSQL schemas and safe repository implementations.
              </p>
            </header>

            <div className="flex flex-col gap-32">
              {SECTIONS.map((section, idx) => (
                <section key={section.id} id={section.id} className="scroll-mt-32 group">
                  <div className="flex flex-col lg:flex-row lg:items-end justify-between gap-6 mb-8">
                    <div className="max-w-xl">
                      <div className="flex items-center gap-4 mb-4">
                        <span className="flex items-center justify-center w-8 h-8 rounded-full bg-white/5 border border-white/10 text-xs font-mono text-zinc-400 shadow-inner">
                          0{idx + 1}
                        </span>
                        <h2 className="text-2xl font-bold tracking-tight text-zinc-100">{section.title}</h2>
                      </div>
                      <p className="text-base text-zinc-400 leading-relaxed">
                        {section.description}
                      </p>
                    </div>
                    {section.id !== "repo-generation" && (
                      <Link
                        href={`/?load=${encodeURIComponent(section.exampleLabel)}&mode=sql`}
                        className="inline-flex items-center justify-center gap-2 text-sm px-5 py-2.5 rounded-full bg-zinc-800/50 hover:bg-zinc-800 text-zinc-200 transition-all border border-zinc-700/50 hover:border-zinc-600 focus-visible:outline focus-visible:outline-2 focus-visible:outline-zinc-400 font-medium whitespace-nowrap"
                      >
                        Try Example <span aria-hidden="true">&rarr;</span>
                      </Link>
                    )}
                    {section.id === "repo-generation" && (
                      <Link
                        href={`/?load=${encodeURIComponent(section.exampleLabel)}&mode=code`}
                        className="inline-flex items-center justify-center gap-2 text-sm px-5 py-2.5 rounded-full bg-zinc-100 hover:bg-white text-zinc-950 transition-all focus-visible:outline focus-visible:outline-2 focus-visible:outline-zinc-400 font-medium whitespace-nowrap shadow-sm"
                      >
                        Try Code Gen <span aria-hidden="true">&rarr;</span>
                      </Link>
                    )}
                  </div>

                  <div className="grid grid-cols-1 xl:grid-cols-2 gap-6 relative">
                    <div className="absolute inset-0 bg-zinc-800/5 blur-3xl -z-10 rounded-full opacity-0 group-hover:opacity-100 transition-opacity duration-700" />
                    <div className="w-full">
                      <CodeBlock label="Source (Go)" code={section.goCode} isGo={true} />
                    </div>
                    <div className="w-full">
                      <CodeBlock label={section.id === "repo-generation" ? "Output (Go)" : "Output (SQL)"} code={section.sqlCode} />
                    </div>
                  </div>
                </section>
              ))}
            </div>

            <footer className="mt-32 pt-10 border-t border-white/5 flex flex-col md:flex-row justify-between items-center gap-4">
              <p className="text-sm text-zinc-500">
                structify — Go structs to PostgreSQL, instantly.
              </p>
              <a
                href="https://github.com/n0xum/structify"
                target="_blank"
                rel="noreferrer"
                className="text-sm text-zinc-400 hover:text-white transition-colors"
              >
                GitHub Repository
              </a>
            </footer>
          </div >
        </main >
      </div >
    </div >
  );
}

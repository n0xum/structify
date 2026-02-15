import Link from "next/link";

export const metadata = {
  title: "Docs — structify",
  description: "All structify db: tag features with Go examples and SQL output.",
};

type Section = {
  id: string;
  title: string;
  exampleLabel: string;
  description: string;
  goCode: string;
  sqlCode: string;
};

const SECTIONS: Section[] = [
  {
    id: "basic",
    title: "Basic",
    exampleLabel: "User",
    description:
      "Define a table by declaring a Go struct. Fields are mapped to columns using their Go type. Use db:\"pk\" for the primary key and db:\"unique\" for a UNIQUE constraint. Fields with no tag are plain nullable columns.",
    goCode: `package models

type User struct {
\tID       int64  \`db:"pk"\`
\tUsername string \`db:"unique"\`
\tEmail    string
\tActive   bool
\tCreated  int64
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
    description:
      "Add column-level constraints inline. db:\"check:<expr>\" emits a CHECK constraint, db:\"default:<value>\" sets a DEFAULT, and db:\"enum:<a>,<b>,...\" generates an IN-list CHECK for allowed values.",
    goCode: `package models

type Person struct {
\tID        int64  \`db:"pk"\`
\tName      string \`db:"check:length(name) > 0"\`
\tAge       int    \`db:"check:age >= 18,default:18"\`
\tStatus    string \`db:"enum:active,inactive,banned"\`
\tActive    bool   \`db:"default:true"\`
\tCreatedAt int64  \`db:"default:now()"\`
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
    description:
      "Create indexes with db:\"index\" or db:\"unique_index\". Give multiple fields the same named index to produce a composite index. Unnamed indexes get an auto-generated name.",
    goCode: `package models

type Article struct {
\tID       int64  \`db:"pk"\`
\tSlug     string \`db:"unique_index:uq_slug"\`
\tTitle    string \`db:"index:idx_title_lang"\`
\tLanguage string \`db:"index:idx_title_lang"\`
\tViews    int64  \`db:"index"\`
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
    description:
      "Reference another table with db:\"fk:<table>,<col>\". Add on_delete:CASCADE (or SET NULL / RESTRICT) to control cascading behaviour. The constraint name is generated automatically.",
    goCode: `package models

type User struct {
\tID    int64  \`db:"pk"\`
\tEmail string \`db:"unique"\`
}

type Post struct {
\tID      int64  \`db:"pk"\`
\tUserID  int64  \`db:"fk:users,id,on_delete:CASCADE"\`
\tTitle   string
\tContent string
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
    id: "composite-keys",
    title: "Composite Keys",
    exampleLabel: "Composite PK & FK",
    description:
      "Mark multiple fields with db:\"pk\" to form a composite primary key. For composite foreign keys, give them the same constraint name as the first tag segment: db:\"fk:<name>,<table>,<col>\".",
    goCode: `package models

type OrderItem struct {
\tOrderID   int64 \`db:"pk"\`
\tProductID int64 \`db:"pk"\`
\tQuantity  int
\tPrice     float64
}

type OrderItemNote struct {
\tNoteID    int64  \`db:"pk"\`
\tOrderID   int64  \`db:"fk:fk_oi,order_items,order_id,on_delete:CASCADE"\`
\tProductID int64  \`db:"fk:fk_oi,order_items,product_id"\`
\tNote      string
}`,
    sqlCode: `CREATE TABLE order_items (
  order_id   BIGINT NOT NULL,
  product_id BIGINT NOT NULL,
  quantity   INTEGER NOT NULL,
  price      DOUBLE PRECISION NOT NULL,
  PRIMARY KEY (order_id, product_id)
);

CREATE TABLE order_item_notes (
  note_id    BIGINT PRIMARY KEY,
  order_id   BIGINT NOT NULL,
  product_id BIGINT NOT NULL,
  note       TEXT NOT NULL,
  CONSTRAINT fk_oi FOREIGN KEY (order_id, product_id)
    REFERENCES order_items (order_id, product_id) ON DELETE CASCADE
);`,
  },
];

const NAV_LINKS = SECTIONS.map((s) => ({ id: s.id, title: s.title }));

export default function DocsPage() {
  return (
    <div className="flex min-h-screen bg-zinc-950 text-zinc-100">
      {/* Sidebar — desktop only */}
      <aside className="hidden md:flex flex-col sticky top-0 h-screen w-52 border-r border-zinc-800 py-8 px-4 gap-1 shrink-0">
        <Link
          href="/"
          className="text-xs text-zinc-500 hover:text-zinc-300 transition-colors mb-6 font-mono"
        >
          ← structify
        </Link>
        <p className="text-xs font-semibold text-zinc-500 uppercase tracking-widest mb-2">
          Features
        </p>
        {NAV_LINKS.map((link) => (
          <a
            key={link.id}
            href={`#${link.id}`}
            className="text-sm text-zinc-400 hover:text-zinc-100 transition-colors py-1 rounded"
          >
            {link.title}
          </a>
        ))}
      </aside>

      {/* Content */}
      <div className="flex flex-col flex-1 min-w-0 px-4 py-8 md:px-10 max-w-5xl">
        {/* Mobile header */}
        <div className="flex items-center gap-3 md:hidden mb-8">
          <Link
            href="/"
            className="text-sm text-zinc-400 hover:text-zinc-200 transition-colors font-mono"
          >
            ← structify
          </Link>
        </div>

        <h1 className="text-2xl font-bold text-zinc-100 mb-2">Docs</h1>
        <p className="text-zinc-400 text-sm mb-10">
          All <code className="font-mono text-sky-400">db:</code> tag features with Go input and SQL
          output.
        </p>

        <div className="flex flex-col gap-16">
          {SECTIONS.map((section) => (
            <section key={section.id} id={section.id} className="scroll-mt-8">
              <div className="flex items-center justify-between gap-4 mb-3 flex-wrap">
                <h2 className="text-lg font-semibold text-zinc-100">{section.title}</h2>
                <Link
                  href={`/?load=${encodeURIComponent(section.exampleLabel)}`}
                  className="text-sm px-3 py-1.5 rounded-md bg-sky-700 hover:bg-sky-600 text-white transition-colors font-medium focus-visible:outline focus-visible:outline-2 focus-visible:outline-sky-400"
                >
                  Try it →
                </Link>
              </div>
              <p className="text-zinc-400 text-sm mb-5">{section.description}</p>

              <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
                {/* Go block */}
                <div className="flex flex-col gap-1.5">
                  <span className="text-xs font-medium text-zinc-500 uppercase tracking-widest">
                    Go
                  </span>
                  <pre className="bg-zinc-900 border border-zinc-800 rounded-lg p-4 overflow-x-auto">
                    <code className="font-mono text-sm text-zinc-300 whitespace-pre">
                      {section.goCode}
                    </code>
                  </pre>
                </div>

                {/* SQL block */}
                <div className="flex flex-col gap-1.5">
                  <span className="text-xs font-medium text-zinc-500 uppercase tracking-widest">
                    SQL
                  </span>
                  <pre className="bg-zinc-900 border border-zinc-800 rounded-lg p-4 overflow-x-auto">
                    <code className="font-mono text-sm text-zinc-300 whitespace-pre">
                      {section.sqlCode}
                    </code>
                  </pre>
                </div>
              </div>
            </section>
          ))}
        </div>
      </div>
    </div>
  );
}

import Link from "next/link";
import { notFound } from "next/navigation";
import { HighlightedCode } from "@/components/HighlightedCode";
import { getCategory } from "@/app/docs/content";

function CodeBlock({ label, code, language }: { label: string; code: string; language: "go" | "sql" | "bash" }) {
  return (
    <div className="h-full overflow-hidden rounded-2xl border border-zinc-800 bg-zinc-900/70 shadow-2xl">
      <div className="flex items-center gap-3 border-b border-zinc-800 bg-zinc-800/30 px-4 py-3">
        <div className="flex gap-1.5 opacity-60">
          <div className="h-3 w-3 rounded-full bg-zinc-500" />
          <div className="h-3 w-3 rounded-full bg-zinc-500" />
          <div className="h-3 w-3 rounded-full bg-zinc-500" />
        </div>
        <span className="text-xs font-semibold uppercase tracking-widest text-zinc-400">{label}</span>
      </div>
      <div className="overflow-x-auto p-4 text-sm">
        <HighlightedCode code={code} language={language} />
      </div>
    </div>
  );
}

export function DocsCategoryPage({ categoryId }: { categoryId: "sql-schema" | "repository" | "project-guide" }) {
  const category = getCategory(categoryId);

  if (!category) {
    notFound();
  }

  return (
    <div className="mx-auto max-w-6xl">
      <div className="mb-8 lg:hidden">
        <Link
          href="/"
          className="inline-flex items-center gap-2 rounded-full border border-white/10 bg-white/5 px-4 py-2 text-sm font-medium text-zinc-300"
        >
          &larr; Back to Editor
        </Link>
      </div>

      <header className="mb-16 space-y-5">
        <h1 className="text-4xl font-extrabold tracking-tight md:text-5xl">Docs &amp; Features</h1>
        <p className="max-w-4xl text-lg leading-relaxed text-zinc-300">
          The documentation is grouped by <strong>SQL Schema</strong>, <strong>Repository</strong>, and <strong>Project Guide</strong>.
          Start with tags, continue with repository generation workflow, then use the project guide for local development and troubleshooting.
        </p>
      </header>

      <section id={category.id} className="scroll-mt-28">
        <div className="mb-7 rounded-2xl border border-white/10 bg-white/[0.03] p-6">
          <h2 className="text-2xl font-bold tracking-tight text-zinc-100">{category.title}</h2>
          <p className="mt-3 max-w-4xl text-zinc-300">{category.intro}</p>
        </div>

        <div className="space-y-8">
          {category.sections.map((section, index) => {
            const isRepoCode = section.outputLabel === "Output (Go)";
            const sourceLabel = section.sourceLabel ?? "Source (Go)";
            const sourceLanguage = section.sourceLanguage ?? "go";

            return (
              <article
                key={section.id}
                id={section.id}
                className="scroll-mt-28 rounded-2xl border border-white/10 bg-zinc-950/50 p-6 lg:p-8"
              >
                <div className="mb-6 flex flex-wrap items-start justify-between gap-4">
                  <div className="max-w-3xl">
                    <div className="mb-3 flex items-center gap-3">
                      <span className="inline-flex h-8 w-8 items-center justify-center rounded-full border border-white/10 bg-white/5 font-mono text-xs text-zinc-300">
                        {String(index + 1).padStart(2, "0")}
                      </span>
                      <h3 className="text-xl font-semibold text-zinc-100">{section.title}</h3>
                    </div>
                    <p className="leading-relaxed text-zinc-300">{section.summary}</p>
                  </div>
                  {section.tryLink && (
                    <Link
                      href={section.tryLink.href}
                      className="inline-flex items-center gap-2 rounded-full border border-zinc-700 bg-zinc-800/40 px-4 py-2 text-sm font-medium text-zinc-100 transition hover:bg-zinc-700/40"
                    >
                      {section.tryLink.label} <span aria-hidden="true">&rarr;</span>
                    </Link>
                  )}
                </div>

                <div className="grid grid-cols-1 gap-5 xl:grid-cols-2">
                  <CodeBlock label={sourceLabel} code={section.sourceCode} language={sourceLanguage} />
                  <CodeBlock label={section.outputLabel} code={section.outputCode} language={isRepoCode ? "go" : "sql"} />
                </div>
              </article>
            );
          })}
        </div>
      </section>

      <footer className="mt-16 flex flex-col items-start justify-between gap-4 border-t border-white/10 pt-8 text-sm text-zinc-500 md:flex-row md:items-center">
        <p>structify - Go structs to PostgreSQL and repository code.</p>
        <a
          href="https://github.com/n0xum/structify"
          target="_blank"
          rel="noreferrer"
          className="text-zinc-300 transition hover:text-white"
        >
          GitHub Repository
        </a>
      </footer>
    </div>
  );
}

"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { ThemeToggle } from "@/components/ThemeToggle";
import { CATEGORIES } from "@/app/docs/content";

export function DocsSidebar() {
  const pathname = usePathname();

  return (
    <aside className="sticky top-0 hidden h-screen w-80 shrink-0 flex-col border-r border-[var(--color-border)] bg-[var(--color-bg-overlay)] px-8 py-12 backdrop-blur-3xl lg:flex">
      <div className="mb-8 flex items-center justify-between">
        <Link
          href="/"
          className="inline-flex items-center gap-2 text-sm font-medium text-[var(--color-text-secondary)] transition-colors hover:text-[var(--color-text-primary)]"
        >
          <span aria-hidden="true">&larr;</span> Back to Editor
        </Link>
        <ThemeToggle />
      </div>

      <div className="mb-4 text-xs font-bold uppercase tracking-[0.2em] text-[var(--color-text-secondary)]">Documentation</div>
      <nav className="space-y-5 overflow-y-auto pr-2">
        {CATEGORIES.map((category) => {
          const categoryPath = `/docs/${category.id}`;
          const isActive = pathname === categoryPath;

          return (
            <details key={category.id} open className="rounded-xl border border-[var(--color-border)] bg-[var(--color-bg-elevated)] p-3">
              <summary className="list-none">
                <Link
                  href={categoryPath}
                  className={`block rounded-md px-2 py-1.5 text-sm font-semibold transition ${
                    isActive
                      ? "bg-[var(--color-accent-soft)] text-[var(--color-text-primary)]"
                      : "text-[var(--color-text-secondary)] hover:bg-[var(--color-bg-subtle)]"
                  }`}
                >
                  {category.title}
                </Link>
              </summary>
              <div className="mt-2 flex flex-col gap-1 border-l border-[var(--color-border)] pl-3">
                {category.sections.map((section) => (
                  <a
                    key={section.id}
                    href={`${categoryPath}#${section.id}`}
                    className="rounded-md px-2 py-1.5 text-sm text-[var(--color-text-muted)] transition hover:bg-[var(--color-bg-subtle)] hover:text-[var(--color-text-primary)]"
                  >
                    {section.title}
                  </a>
                ))}
              </div>
            </details>
          );
        })}
      </nav>
    </aside>
  );
}

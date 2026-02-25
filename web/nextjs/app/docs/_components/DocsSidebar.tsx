"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { CATEGORIES } from "@/app/docs/content";

export function DocsSidebar() {
  const pathname = usePathname();

  return (
    <aside className="sticky top-0 hidden h-screen w-80 shrink-0 flex-col border-r border-white/5 bg-black/20 px-8 py-12 backdrop-blur-3xl lg:flex">
      <Link
        href="/"
        className="mb-10 inline-flex items-center gap-2 text-sm font-medium text-zinc-400 transition-colors hover:text-white"
      >
        <span aria-hidden="true">&larr;</span> Back to Editor
      </Link>

      <div className="mb-4 text-xs font-bold uppercase tracking-[0.2em] text-zinc-300">Documentation</div>
      <nav className="space-y-5 overflow-y-auto pr-2">
        {CATEGORIES.map((category) => {
          const categoryPath = `/docs/${category.id}`;
          const isActive = pathname === categoryPath;

          return (
            <details key={category.id} open className="rounded-xl border border-white/5 bg-white/[0.02] p-3">
              <summary className="list-none">
                <Link
                  href={categoryPath}
                  className={`block rounded-md px-2 py-1.5 text-sm font-semibold transition ${
                    isActive ? "bg-white/10 text-zinc-100" : "text-zinc-200 hover:bg-white/5"
                  }`}
                >
                  {category.title}
                </Link>
              </summary>
              <div className="mt-2 flex flex-col gap-1 border-l border-white/10 pl-3">
                {category.sections.map((section) => (
                  <a
                    key={section.id}
                    href={`${categoryPath}#${section.id}`}
                    className="rounded-md px-2 py-1.5 text-sm text-zinc-400 transition hover:bg-white/5 hover:text-zinc-100"
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

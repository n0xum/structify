"use client";

import React, { useState } from "react";
import { TAG_GROUPS, METHOD_GROUPS } from "@/lib/db-tags";
import type { Mode } from "./ModeSelector";

type TagReferenceProps = {
  mode: Mode;
};

export function TagReference({ mode }: Readonly<TagReferenceProps>) {
  const [open, setOpen] = useState(false);
  const isSQLMode = mode === "sql";

  return (
    <div className="relative">
      <button
        onClick={() => setOpen((o) => !o)}
        aria-label={isSQLMode ? "Show db tag reference" : "Show method naming reference"}
        aria-expanded={open}
        className="flex h-6 w-6 items-center justify-center rounded-full border border-[var(--color-border-strong)] text-xs text-[var(--color-text-secondary)] transition-colors hover:border-[var(--color-accent)] hover:text-[var(--color-text-primary)] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[var(--color-accent)]"
      >
        ?
      </button>

      {open && (
        <div
          role="tooltip"
          className="absolute left-0 top-8 z-10 w-[min(44rem,calc(100vw-2rem))] rounded-lg border border-[var(--color-border)] bg-[var(--color-bg-elevated)] p-3 shadow-xl md:w-[42rem]"
        >
          {isSQLMode ? (
            <>
              <p className="mb-2 text-xs font-semibold uppercase tracking-widest text-[var(--color-text-secondary)]">
                db tag reference
              </p>
              <table className="w-full table-fixed text-xs">
                <tbody>
                  {TAG_GROUPS.map(({ group, tags }) => (
                    <React.Fragment key={group}>
                      <tr>
                        <td
                          colSpan={2}
                          className="pb-1 pt-2 text-[10px] font-semibold uppercase tracking-widest text-[var(--color-text-muted)]"
                        >
                          {group}
                        </td>
                      </tr>
                      {tags.map(({ tag, description }) => (
                        <tr key={tag} className="border-t border-[var(--color-border)]">
                          <td className="w-[44%] break-words py-1.5 pr-3 font-mono text-[var(--color-accent)] [overflow-wrap:anywhere]">{tag}</td>
                          <td className="py-1.5 leading-relaxed text-[var(--color-text-secondary)]">{description}</td>
                        </tr>
                      ))}
                    </React.Fragment>
                  ))}
                </tbody>
              </table>
            </>
          ) : (
            <>
              <p className="mb-2 text-xs font-semibold uppercase tracking-widest text-[var(--color-text-secondary)]">
                method naming reference
              </p>
              <p className="mb-2 text-xs leading-relaxed text-[var(--color-text-secondary)]">
                Use interface methods with <code className="rounded bg-[var(--color-bg-subtle)] px-1 py-0.5 text-[var(--color-text-primary)]">context.Context</code> as the first argument and return <code className="rounded bg-[var(--color-bg-subtle)] px-1 py-0.5 text-[var(--color-text-primary)]">error</code> as the last value.
              </p>
              <table className="w-full table-fixed text-xs">
                <tbody>
                  {METHOD_GROUPS.map(({ group, methods }) => (
                    <React.Fragment key={group}>
                      <tr>
                        <td
                          colSpan={2}
                          className="pb-1 pt-2 text-[10px] font-semibold uppercase tracking-widest text-[var(--color-text-muted)]"
                        >
                          {group}
                        </td>
                      </tr>
                      {methods.map(({ pattern, description }) => (
                        <tr key={pattern} className="border-t border-[var(--color-border)]">
                          <td className="w-[48%] break-words py-1.5 pr-3 font-mono text-[var(--color-accent)] [overflow-wrap:anywhere]">{pattern}</td>
                          <td className="py-1.5 leading-relaxed text-[var(--color-text-secondary)]">{description}</td>
                        </tr>
                      ))}
                    </React.Fragment>
                  ))}
                </tbody>
              </table>
            </>
          )}
        </div>
      )}
    </div>
  );
}

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
        className="flex items-center justify-center w-6 h-6 rounded-full border border-zinc-600 text-zinc-400 hover:text-zinc-200 hover:border-zinc-400 text-xs transition-colors focus-visible:outline focus-visible:outline-2 focus-visible:outline-zinc-400"
      >
        ?
      </button>

      {open && (
        <div
          role="tooltip"
          className="absolute left-0 top-8 z-10 w-[min(44rem,calc(100vw-2rem))] rounded-lg border border-zinc-700 bg-zinc-900 shadow-xl p-3 md:w-[42rem]"
        >
          {isSQLMode ? (
            <>
              <p className="text-xs font-semibold text-zinc-300 mb-2 uppercase tracking-widest">
                db tag reference
              </p>
              <table className="w-full text-xs table-fixed">
                <tbody>
                  {TAG_GROUPS.map(({ group, tags }) => (
                    <React.Fragment key={group}>
                      <tr>
                        <td
                          colSpan={2}
                          className="pt-2 pb-1 text-zinc-500 uppercase tracking-widest text-[10px] font-semibold"
                        >
                          {group}
                        </td>
                      </tr>
                      {tags.map(({ tag, description }) => (
                        <tr key={tag} className="border-t border-zinc-800">
                          <td className="w-[44%] py-1.5 pr-3 font-mono text-sky-400 break-words [overflow-wrap:anywhere]">{tag}</td>
                          <td className="py-1.5 text-zinc-400 leading-relaxed">{description}</td>
                        </tr>
                      ))}
                    </React.Fragment>
                  ))}
                </tbody>
              </table>
            </>
          ) : (
            <>
              <p className="text-xs font-semibold text-zinc-300 mb-2 uppercase tracking-widest">
                method naming reference
              </p>
              <p className="mb-2 text-xs text-zinc-400 leading-relaxed">
                Use interface methods with <code className="rounded bg-zinc-800 px-1 py-0.5 text-zinc-300">context.Context</code> as the first argument and return <code className="rounded bg-zinc-800 px-1 py-0.5 text-zinc-300">error</code> as the last value.
              </p>
              <table className="w-full text-xs table-fixed">
                <tbody>
                  {METHOD_GROUPS.map(({ group, methods }) => (
                    <React.Fragment key={group}>
                      <tr>
                        <td
                          colSpan={2}
                          className="pt-2 pb-1 text-zinc-500 uppercase tracking-widest text-[10px] font-semibold"
                        >
                          {group}
                        </td>
                      </tr>
                      {methods.map(({ pattern, description }) => (
                        <tr key={pattern} className="border-t border-zinc-800">
                          <td className="w-[48%] py-1.5 pr-3 font-mono text-sky-400 break-words [overflow-wrap:anywhere]">{pattern}</td>
                          <td className="py-1.5 text-zinc-400 leading-relaxed">{description}</td>
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

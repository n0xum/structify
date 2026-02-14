"use client";

import { useState } from "react";

const TAGS = [
  { tag: `db:"pk"`, description: "Mark field as primary key" },
  { tag: `db:"unique"`, description: "Add UNIQUE constraint" },
  { tag: `db:"-"`, description: "Exclude field from schema and generated code" },
  { tag: `db:"table:name"`, description: 'Override the generated table name (on the struct)' },
];

export function TagReference() {
  const [open, setOpen] = useState(false);

  return (
    <div className="relative">
      <button
        onClick={() => setOpen((o) => !o)}
        aria-label="Show db tag reference"
        aria-expanded={open}
        className="flex items-center justify-center w-6 h-6 rounded-full border border-zinc-600 text-zinc-400 hover:text-zinc-200 hover:border-zinc-400 text-xs transition-colors focus-visible:outline focus-visible:outline-2 focus-visible:outline-sky-500"
      >
        ?
      </button>

      {open && (
        <div
          role="tooltip"
          className="absolute left-0 top-8 z-10 w-72 rounded-lg border border-zinc-700 bg-zinc-900 shadow-xl p-3"
        >
          <p className="text-xs font-semibold text-zinc-300 mb-2 uppercase tracking-widest">
            db tag reference
          </p>
          <table className="w-full text-xs">
            <tbody>
              {TAGS.map(({ tag, description }) => (
                <tr key={tag} className="border-t border-zinc-800 first:border-0">
                  <td className="py-1.5 pr-3 font-mono text-sky-400 whitespace-nowrap">{tag}</td>
                  <td className="py-1.5 text-zinc-400">{description}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}

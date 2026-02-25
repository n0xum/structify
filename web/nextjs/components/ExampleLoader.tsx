"use client";

import { EXAMPLES, REPO_EXAMPLES } from "@/lib/examples";
import type { Mode } from "./ModeSelector";

type ExampleLoaderProps = {
  mode: Mode;
  onSelect: (source: string) => void;
};

export function ExampleLoader({ mode, onSelect }: Readonly<ExampleLoaderProps>) {
  const examples = mode === "repo" ? REPO_EXAMPLES : EXAMPLES;

  function handleChange(e: React.ChangeEvent<HTMLSelectElement>) {
    const example = examples.find((ex) => ex.label === e.target.value);
    if (example) onSelect(example.source);
    e.target.value = "";
  }

  return (
    <select
      defaultValue=""
      onChange={handleChange}
      aria-label={mode === "repo" ? "Load example interface" : "Load example struct"}
      className="bg-zinc-800 border border-zinc-700 rounded-md px-3 py-1.5 text-sm text-zinc-300 focus:outline-none focus:ring-2 focus:ring-zinc-400 cursor-pointer"
    >
      <option value="" disabled>
        Load example…
      </option>
      {examples.map((ex) => (
        <option key={ex.label} value={ex.label}>
          {ex.label} — {ex.description}
        </option>
      ))}
    </select>
  );
}

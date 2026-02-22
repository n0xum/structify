"use client";

import { EXAMPLES } from "@/lib/examples";

type ExampleLoaderProps = {
  onSelect: (source: string) => void;
};

export function ExampleLoader({ onSelect }: Readonly<ExampleLoaderProps>) {
  function handleChange(e: React.ChangeEvent<HTMLSelectElement>) {
    const example = EXAMPLES.find((ex) => ex.label === e.target.value);
    if (example) onSelect(example.source);
    e.target.value = "";
  }

  return (
    <select
      defaultValue=""
      onChange={handleChange}
      aria-label="Load example struct"
      className="bg-zinc-800 border border-zinc-700 rounded-md px-3 py-1.5 text-sm text-zinc-300 focus:outline-none focus:ring-2 focus:ring-zinc-400 cursor-pointer"
    >
      <option value="" disabled>
        Load example…
      </option>
      {EXAMPLES.map((ex) => (
        <option key={ex.label} value={ex.label}>
          {ex.label} — {ex.description}
        </option>
      ))}
    </select>
  );
}

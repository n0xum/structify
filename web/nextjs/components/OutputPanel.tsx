"use client";

import { useState } from "react";
import { Editor } from "./Editor";
import type { Mode } from "./ModeSelector";

type OutputPanelProps = {
  output: string;
  mode: Mode;
};

export function OutputPanel({ output, mode }: Readonly<OutputPanelProps>) {
  const [copied, setCopied] = useState(false);

  async function handleCopy() {
    await navigator.clipboard.writeText(output);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  }

  return (
    <div className="flex h-full flex-col gap-2">
      <div className="flex items-center justify-between">
        <span className="text-xs font-medium uppercase tracking-widest text-[var(--color-text-secondary)]">
          Output
        </span>
        {output && (
          <button
            onClick={handleCopy}
            aria-label="Copy output to clipboard"
            className="rounded px-2 py-1 text-xs text-[var(--color-text-secondary)] transition-colors hover:text-[var(--color-text-primary)] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[var(--color-accent)]"
          >
            {copied ? "Copied!" : "Copy"}
          </button>
        )}
      </div>

      {output ? (
        <div className="min-h-0 flex-1">
          <Editor
            value={output}
            language={mode === "sql" ? "sql" : "go"}
            readOnly
            label="Generated output"
            id="output-editor"
          />
        </div>
      ) : (
        <div className="flex flex-1 items-center justify-center rounded-lg border border-dashed border-[var(--color-border)] text-sm text-[var(--color-text-muted)]">
          Output will appear here
        </div>
      )}
    </div>
  );
}

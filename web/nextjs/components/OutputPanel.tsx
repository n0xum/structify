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
    <div className="flex flex-col h-full gap-2">
      <div className="flex items-center justify-between">
        <span className="text-xs font-medium text-zinc-400 uppercase tracking-widest">
          Output
        </span>
        {output && (
          <button
            onClick={handleCopy}
            aria-label="Copy output to clipboard"
            className="text-xs text-zinc-400 hover:text-zinc-200 transition-colors px-2 py-1 rounded focus-visible:outline focus-visible:outline-2 focus-visible:outline-sky-500"
          >
            {copied ? "Copied!" : "Copy"}
          </button>
        )}
      </div>

      {output ? (
        <div className="flex-1 min-h-0">
          <Editor
            value={output}
            language={mode === "sql" ? "sql" : "go"}
            readOnly
            label="Generated output"
            id="output-editor"
          />
        </div>
      ) : (
        <div className="flex-1 flex items-center justify-center rounded-lg border border-dashed border-zinc-700 text-sm text-zinc-500">
          Output will appear here
        </div>
      )}
    </div>
  );
}

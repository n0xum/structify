"use client";

export type Mode = "sql" | "code";

type ModeSelectorProps = {
  mode: Mode;
  onChange: (mode: Mode) => void;
  packageName: string;
  onPackageChange: (name: string) => void;
  packageError: string | null;
};

export function ModeSelector({
  mode,
  onChange,
  packageName,
  onPackageChange,
  packageError,
}: Readonly<ModeSelectorProps>) {
  return (
    <div className="flex flex-col gap-3">
      <div
        className="flex rounded-lg border border-zinc-700 overflow-hidden"
        role="tablist"
        aria-label="Output mode"
      >
        <button
          role="tab"
          aria-selected={mode === "sql"}
          onClick={() => onChange("sql")}
          className={`flex-1 py-2 px-4 text-sm font-medium transition-colors focus-visible:outline focus-visible:outline-2 focus-visible:outline-zinc-300 ${mode === "sql"
            ? "bg-zinc-100 text-zinc-900 shadow-sm"
            : "bg-zinc-900 text-zinc-400 hover:text-zinc-200"
            }`}
        >
          SQL Schema
        </button>
        <button
          role="tab"
          aria-selected={mode === "code"}
          onClick={() => onChange("code")}
          className={`flex-1 py-2 px-4 text-sm font-medium transition-colors focus-visible:outline focus-visible:outline-2 focus-visible:outline-zinc-300 ${mode === "code"
            ? "bg-zinc-100 text-zinc-900 shadow-sm"
            : "bg-zinc-900 text-zinc-400 hover:text-zinc-200"
            }`}
        >
          Repository Implementation
        </button>
      </div>

      {mode === "code" && (
        <div className="flex flex-col gap-1">
          <label htmlFor="pkg-input" className="text-xs text-zinc-400">
            Package name
          </label>
          <input
            id="pkg-input"
            type="text"
            value={packageName}
            onChange={(e) => onPackageChange(e.target.value)}
            placeholder="models"
            aria-describedby={packageError ? "pkg-error" : undefined}
            className="bg-zinc-900 border border-zinc-700/50 rounded-md px-3 py-1.5 text-sm text-zinc-100 placeholder-zinc-500 focus:outline-none focus:ring-2 focus:ring-zinc-400"
          />
          {packageError && (
            <p id="pkg-error" className="text-xs text-red-400">
              {packageError}
            </p>
          )}
        </div>
      )}
    </div>
  );
}

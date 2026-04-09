"use client";

export type Mode = "sql" | "repo";

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
        className="flex overflow-hidden rounded-lg border border-[var(--color-border)]"
        role="tablist"
        aria-label="Output mode"
      >
        <button
          role="tab"
          aria-selected={mode === "sql"}
          onClick={() => onChange("sql")}
          className={`flex-1 px-4 py-2 text-sm font-medium transition-colors focus-visible:outline focus-visible:outline-2 focus-visible:outline-[var(--color-accent)] ${mode === "sql"
            ? "bg-[var(--color-text-primary)] text-[var(--color-text-inverse)] shadow-sm"
            : "bg-[var(--color-bg-elevated)] text-[var(--color-text-secondary)] hover:text-[var(--color-text-primary)]"
            }`}
        >
          SQL Schema
        </button>
        <button
          role="tab"
          aria-selected={mode === "repo"}
          onClick={() => onChange("repo")}
          className={`flex-1 px-4 py-2 text-sm font-medium transition-colors focus-visible:outline focus-visible:outline-2 focus-visible:outline-[var(--color-accent)] ${mode === "repo"
            ? "bg-[var(--color-text-primary)] text-[var(--color-text-inverse)] shadow-sm"
            : "bg-[var(--color-bg-elevated)] text-[var(--color-text-secondary)] hover:text-[var(--color-text-primary)]"
            }`}
        >
          Interface Repository
        </button>
      </div>

      {mode === "repo" && (
        <div className="flex flex-col gap-1">
          <label htmlFor="pkg-input" className="text-xs text-[var(--color-text-secondary)]">
            Package name
          </label>
          <input
            id="pkg-input"
            type="text"
            value={packageName}
            onChange={(e) => onPackageChange(e.target.value)}
            placeholder="models"
            aria-describedby={packageError ? "pkg-error" : undefined}
            className="rounded-md border border-[var(--color-border)] bg-[var(--color-bg-elevated)] px-3 py-1.5 text-sm text-[var(--color-text-primary)] placeholder:text-[var(--color-text-muted)] focus:outline-none focus:ring-2 focus:ring-[var(--color-accent)]"
          />
          {packageError && (
            <p id="pkg-error" className="text-xs text-[var(--color-danger)]">
              {packageError}
            </p>
          )}
        </div>
      )}
    </div>
  );
}

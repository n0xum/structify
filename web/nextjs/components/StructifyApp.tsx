"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import Link from "next/link";
import { Editor } from "./Editor";
import { OutputPanel } from "./OutputPanel";
import { ModeSelector, type Mode } from "./ModeSelector";
import { ExampleLoader } from "./ExampleLoader";
import { TagReference } from "./TagReference";
import { ErrorBoundary } from "./ErrorBoundary";
import { ThemeToggle } from "./ThemeToggle";
import { generateSQL, generateRepository } from "@/lib/api";
import { validateInput, validatePackageName } from "@/lib/validation";
import { EXAMPLES } from "@/lib/examples";

const STORAGE_KEY = "structify_input";
const DEFAULT_PLACEHOLDER = `package models

import "context"

type User struct {
	ID       int64  \`db:"pk"\`
	Username string \`db:"unique"\`
	Email    string
	Active   bool
	Created  int64
}

// UserRepository defines the methods to access User data.
// Structify will generate the implementation for this interface.
type UserRepository interface {
	FindByID(ctx context.Context, id int64) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
}`;

export function StructifyApp() {
  const version = "0.1.0";
  const router = useRouter();
  const searchParams = useSearchParams();

  const [input, setInput] = useState("");
  const [output, setOutput] = useState("");
  const [mode, setMode] = useState<Mode>((searchParams.get("mode") as Mode) ?? "sql");
  const [packageName, setPackageName] = useState("models");
  const [packageError, setPackageError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [warnings, setWarnings] = useState<string[]>([]);
  const outputRef = useRef<HTMLDivElement>(null);
  const pendingGenerate = useRef(false);

  useEffect(() => {
    try {
      const saved = localStorage.getItem(STORAGE_KEY);
      if (saved) setInput(saved);
    } catch {
      // localStorage unavailable
    }
  }, []);

  useEffect(() => {
    const loadLabel = searchParams.get("load");
    if (!loadLabel) return;

    const example = EXAMPLES.find((e) => e.label === loadLabel);
    if (!example) return;

    setInput(example.source);
    setOutput("");
    setError(null);
    pendingGenerate.current = true;

    const params = new URLSearchParams(searchParams.toString());
    params.delete("load");
    const newSearch = params.toString();
    router.replace(newSearch ? `?${newSearch}` : "/", { scroll: false });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    try {
      localStorage.setItem(STORAGE_KEY, input);
    } catch {
      // storage full or unavailable
    }
    setWarnings(validateInput(input));
  }, [input]);

  useEffect(() => {
    if (searchParams.get("mode") === mode) return;
    const params = new URLSearchParams(searchParams.toString());
    params.set("mode", mode);
    router.replace(`?${params.toString()}`, { scroll: false });
  }, [mode, router, searchParams]);

  const handleGenerate = useCallback(async (sourceOverride?: string) => {
    const source = sourceOverride ?? input;
    if (!source.trim() || loading) return;

    if (mode === "repo") {
      const pkgErr = validatePackageName(packageName);
      setPackageError(pkgErr);
      if (pkgErr) return;
    }

    setLoading(true);
    setError(null);

    try {
      const result =
        mode === "sql"
          ? await generateSQL(source)
          : await generateRepository(source, packageName);
      setOutput(result);

      if (globalThis.innerWidth < 768) {
        setTimeout(() => outputRef.current?.scrollIntoView({ behavior: "smooth" }), 100);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "An unexpected error occurred.");
    } finally {
      setLoading(false);
    }
  }, [input, mode, packageName, loading]);

  useEffect(() => {
    if (pendingGenerate.current && input) {
      pendingGenerate.current = false;
      handleGenerate(input);
    }
  }, [input, handleGenerate]);

  useEffect(() => {
    function handleKeyDown(e: KeyboardEvent) {
      if ((e.ctrlKey || e.metaKey) && e.key === "Enter") {
        e.preventDefault();
        handleGenerate();
      }
    }
    globalThis.addEventListener("keydown", handleKeyDown);
    return () => globalThis.removeEventListener("keydown", handleKeyDown);
  }, [handleGenerate]);

  function handleModeChange(newMode: Mode) {
    setMode(newMode);
    setOutput("");
    setError(null);
  }

  function handleExampleSelect(source: string) {
    setInput(source);
    setOutput("");
    setError(null);
    pendingGenerate.current = true;
  }

  function handlePackageChange(name: string) {
    setPackageName(name);
    setPackageError(validatePackageName(name));
  }

  const canGenerate = input.trim().length > 0 && !loading;

  return (
    <ErrorBoundary>
      <div className="flex min-h-screen flex-col bg-[var(--color-bg)] text-[var(--color-text-primary)]">
        <header className="flex items-center justify-between gap-4 border-b border-[var(--color-border)] px-4 py-3">
          <div className="flex items-center gap-3">
            <span className="font-mono text-lg font-bold text-[var(--color-text-primary)]">structify</span>
            <span className="rounded-full border border-[var(--color-border)] bg-[var(--color-bg-subtle)] px-2 py-0.5 text-xs text-[var(--color-text-muted)]">
              v{version}
            </span>
            <p className="hidden text-sm text-[var(--color-text-secondary)] sm:block">
              {mode === "sql"
                ? "Go structs to PostgreSQL, instantly."
                : "Go interfaces to repository implementations, instantly."}
            </p>
          </div>

          <div className="flex items-center gap-3">
            <ThemeToggle />
            <ExampleLoader mode={mode} onSelect={handleExampleSelect} />
            <Link
              href="/docs"
              className="rounded text-sm text-[var(--color-text-secondary)] transition-colors hover:text-[var(--color-text-primary)] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[var(--color-accent)]"
            >
              Docs
            </Link>
            <a
              href="https://github.com/n0xum/structify"
              target="_blank"
              rel="noopener noreferrer"
              aria-label="GitHub repository"
              className="rounded text-sm text-[var(--color-text-secondary)] transition-colors hover:text-[var(--color-text-primary)] focus-visible:outline focus-visible:outline-2 focus-visible:outline-[var(--color-accent)]"
            >
              GitHub
            </a>
          </div>
        </header>

        <main className="flex flex-1 flex-col gap-4 p-4 md:flex-row">
          <div className="flex min-h-[400px] flex-1 flex-col gap-3 md:min-h-0">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <span className="text-xs font-medium uppercase tracking-widest text-[var(--color-text-secondary)]">
                  {mode === "sql" ? "Go Struct" : "Go Struct + Interface"}
                </span>
                <TagReference mode={mode} />
              </div>
              <span className="text-xs text-[var(--color-text-muted)]">
                {(input.length / 1024).toFixed(1)} KB
              </span>
            </div>

            <div className="min-h-0 flex-1">
              <Editor
                value={input}
                onChange={setInput}
                language="go"
                placeholder={DEFAULT_PLACEHOLDER}
                label="Go struct input editor"
                id="input-editor"
              />
            </div>

            {warnings.length > 0 && (
              <ul className="flex flex-col gap-1" aria-describedby="input-editor">
                {warnings.map((w) => (
                  <li key={w} className="text-xs text-[var(--color-text-secondary)]">
                    {w}
                  </li>
                ))}
              </ul>
            )}

            {error && (
              <p
                className="rounded-md border border-[var(--color-border)] bg-[var(--color-bg-elevated)] px-3 py-2 text-xs text-[var(--color-danger)]"
                role="alert"
                aria-live="polite"
              >
                {error}
              </p>
            )}

            <div className="flex flex-col gap-3">
              <ModeSelector
                mode={mode}
                onChange={handleModeChange}
                packageName={packageName}
                onPackageChange={handlePackageChange}
                packageError={packageError}
              />
              <button
                onClick={() => handleGenerate()}
                disabled={!canGenerate}
                aria-label="Generate output (Ctrl+Enter)"
                className="w-full rounded-lg bg-[var(--color-text-primary)] py-2.5 text-sm font-semibold text-[var(--color-text-inverse)] shadow-sm transition-colors hover:opacity-90 disabled:cursor-not-allowed disabled:opacity-40 focus-visible:outline focus-visible:outline-2 focus-visible:outline-[var(--color-accent)]"
              >
                {loading ? (
                  <span className="flex items-center justify-center gap-2">
                    <svg className="h-4 w-4 animate-spin" viewBox="0 0 24 24" fill="none">
                      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z" />
                    </svg>
                    Generating...
                  </span>
                ) : (
                  "Generate"
                )}
              </button>
              <p className="text-center text-xs text-[var(--color-text-muted)]">
                or press{" "}
                <kbd className="rounded bg-[var(--color-bg-subtle)] px-1 font-mono text-[var(--color-text-secondary)]">Ctrl+Enter</kbd>
              </p>
            </div>
          </div>

          <div ref={outputRef} className="flex min-h-[300px] flex-1 flex-col md:min-h-0">
            <OutputPanel output={output} mode={mode} />
          </div>
        </main>
      </div>
    </ErrorBoundary>
  );
}

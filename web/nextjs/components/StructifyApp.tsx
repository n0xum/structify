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
import { generateSQL, generateCode, fetchVersion } from "@/lib/api";
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
  const [version, setVersion] = useState<string | null>(null);

  const outputRef = useRef<HTMLDivElement>(null);
  const pendingGenerate = useRef(false);

  // Restore input from localStorage
  useEffect(() => {
    try {
      const saved = localStorage.getItem(STORAGE_KEY);
      if (saved) setInput(saved);
    } catch {
      // localStorage unavailable
    }
  }, []);

  // Load example from ?load= query param (runs after localStorage so it can override)
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

  // Persist input to localStorage
  useEffect(() => {
    try {
      localStorage.setItem(STORAGE_KEY, input);
    } catch {
      // storage full or unavailable
    }
    setWarnings(validateInput(input));
  }, [input]);

  // Sync mode to URL query param only when it differs
  useEffect(() => {
    if (searchParams.get("mode") === mode) return;
    const params = new URLSearchParams(searchParams.toString());
    params.set("mode", mode);
    router.replace(`?${params.toString()}`, { scroll: false });
  }, [mode, router, searchParams]);

  // Fetch backend version
  useEffect(() => {
    fetchVersion().then(setVersion);
  }, []);

  const handleGenerate = useCallback(async (sourceOverride?: string) => {
    const source = sourceOverride ?? input;
    if (!source.trim() || loading) return;

    if (mode === "code") {
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
          : await generateCode(source, packageName);
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

  // Auto-generate after example is loaded
  useEffect(() => {
    if (pendingGenerate.current && input) {
      pendingGenerate.current = false;
      handleGenerate(input);
    }
  }, [input, handleGenerate]);

  // Keyboard shortcut Ctrl/Cmd+Enter
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
      <div className="flex flex-col min-h-screen bg-zinc-950 text-zinc-100">
        {/* Header */}
        <header className="border-b border-zinc-800 px-4 py-3 flex items-center justify-between gap-4 flex-wrap">
          <div className="flex items-center gap-3">
            <span className="font-mono font-bold text-lg text-zinc-100">structify</span>
            {version && (
              <span className="text-xs bg-zinc-800 text-zinc-400 px-2 py-0.5 rounded-full border border-zinc-700">
                v{version}
              </span>
            )}
            <p className="hidden sm:block text-sm text-zinc-400">
              Go structs to PostgreSQL, instantly.
            </p>
          </div>
          <div className="flex items-center gap-3">
            <ExampleLoader onSelect={handleExampleSelect} />
            <Link
              href="/docs"
              className="text-zinc-400 hover:text-zinc-200 transition-colors text-sm focus-visible:outline focus-visible:outline-2 focus-visible:outline-zinc-400 rounded"
            >
              Docs
            </Link>
            <a
              href="https://github.com/n0xum/structify"
              target="_blank"
              rel="noopener noreferrer"
              aria-label="GitHub repository"
              className="text-zinc-400 hover:text-zinc-200 transition-colors text-sm focus-visible:outline focus-visible:outline-2 focus-visible:outline-zinc-400 rounded"
            >
              GitHub
            </a>
          </div>
        </header>

        {/* Main */}
        <main className="flex flex-col flex-1 p-4 gap-4 md:flex-row">
          {/* Left panel — input */}
          <div className="flex flex-col gap-3 flex-1 min-h-[400px] md:min-h-0">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <span className="text-xs font-medium text-zinc-400 uppercase tracking-widest">
                  Go Struct
                </span>
                <TagReference />
              </div>
              <span className="text-xs text-zinc-600">
                {(input.length / 1024).toFixed(1)} KB
              </span>
            </div>

            <div className="flex-1 min-h-0">
              <Editor
                value={input}
                onChange={setInput}
                language="go"
                placeholder={DEFAULT_PLACEHOLDER}
                label="Go struct input editor"
                id="input-editor"
              />
            </div>

            {/* Warnings */}
            {warnings.length > 0 && (
              <ul className="flex flex-col gap-1" aria-describedby="input-editor">
                {warnings.map((w) => (
                  <li key={w} className="text-xs text-zinc-400">
                    {w}
                  </li>
                ))}
              </ul>
            )}

            {/* Error */}
            {error && (
              <p
                className="text-xs text-zinc-300 bg-zinc-900 border border-zinc-800 rounded-md px-3 py-2"
                role="alert"
                aria-live="polite"
              >
                {error}
              </p>
            )}

            {/* Controls */}
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
                className="w-full py-2.5 rounded-lg bg-zinc-100 text-zinc-950 text-sm font-semibold transition-colors hover:bg-white disabled:opacity-40 disabled:cursor-not-allowed focus-visible:outline focus-visible:outline-2 focus-visible:outline-zinc-400 shadow-sm"
              >
                {loading ? (
                  <span className="flex items-center justify-center gap-2">
                    <svg className="animate-spin h-4 w-4" viewBox="0 0 24 24" fill="none">
                      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z" />
                    </svg>
                    Generating…
                  </span>
                ) : (
                  "Generate"
                )}
              </button>
              <p className="text-xs text-zinc-600 text-center">
                or press{" "}
                <kbd className="font-mono bg-zinc-800 px-1 rounded text-zinc-400">Ctrl+Enter</kbd>
              </p>
            </div>
          </div>

          {/* Right panel — output */}
          <div ref={outputRef} className="flex flex-col flex-1 min-h-[300px] md:min-h-0">
            <OutputPanel output={output} mode={mode} />
          </div>
        </main>
      </div>
    </ErrorBoundary>
  );
}

"use client";

import { Component, type ReactNode } from "react";

type Props = { children: ReactNode };
type State = { error: Error | null };

export class ErrorBoundary extends Component<Props, State> {
  state: State = { error: null };

  static getDerivedStateFromError(error: Error): State {
    return { error };
  }

  render() {
    if (this.state.error) {
      return (
        <div className="flex min-h-screen flex-col items-center justify-center gap-4 bg-[var(--color-bg)] p-8 text-center">
          <h1 className="text-xl font-semibold text-[var(--color-text-primary)]">Something went wrong</h1>
          <p className="text-sm text-[var(--color-text-secondary)]">
            Please{" "}
            <a
              href="https://github.com/n0xum/structify/issues"
              className="underline text-[var(--color-accent)] hover:opacity-80"
              target="_blank"
              rel="noopener noreferrer"
            >
              open an issue
            </a>{" "}
            if this keeps happening.
          </p>
          <details className="mt-2 w-full max-w-xl text-left">
            <summary className="cursor-pointer text-xs text-[var(--color-text-muted)]">Error details</summary>
            <pre className="mt-2 overflow-auto rounded bg-[var(--color-bg-elevated)] p-3 text-xs text-[var(--color-danger)]">
              {this.state.error.message}
            </pre>
          </details>
        </div>
      );
    }
    return this.props.children;
  }
}

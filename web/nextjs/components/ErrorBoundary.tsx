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
        <div className="flex min-h-screen flex-col items-center justify-center gap-4 bg-zinc-950 p-8 text-center">
          <h1 className="text-xl font-semibold text-zinc-100">Something went wrong</h1>
          <p className="text-sm text-zinc-400">
            Please{" "}
            <a
              href="https://github.com/n0xum/structify/issues"
              className="underline text-sky-400 hover:text-sky-300"
              target="_blank"
              rel="noopener noreferrer"
            >
              open an issue
            </a>{" "}
            if this keeps happening.
          </p>
          <details className="mt-2 max-w-xl w-full text-left">
            <summary className="cursor-pointer text-xs text-zinc-500">Error details</summary>
            <pre className="mt-2 overflow-auto rounded bg-zinc-900 p-3 text-xs text-red-400">
              {this.state.error.message}
            </pre>
          </details>
        </div>
      );
    }
    return this.props.children;
  }
}

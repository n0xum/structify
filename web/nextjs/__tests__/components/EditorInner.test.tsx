import { useEffect, useState, type ComponentType } from "react";
import { afterEach, beforeEach, describe, it, expect, vi } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import { githubLight } from "@uiw/codemirror-theme-github";
import { intellijTheme } from "@/lib/intellij-theme";
import { THEME_STORAGE_KEY } from "@/lib/themes";
import { ThemeProvider } from "@/components/ThemeProvider";

type CodeMirrorProps = {
  value: string;
  onChange?: (value: string) => void;
  extensions: unknown[];
  theme: unknown;
  readOnly?: boolean;
  placeholder?: string;
  basicSetup: { lineNumbers: boolean; foldGutter: boolean; highlightActiveLine: boolean };
};

const codeMirrorCalls: CodeMirrorProps[] = [];

// Run the real loader of next/dynamic so that EditorInner (the code behind the lazy import) is exercised.
vi.mock("next/dynamic", () => ({
  default: (loader: () => Promise<ComponentType<Record<string, unknown>>>) =>
    function Loaded(props: Record<string, unknown>) {
      const [Inner, setInner] = useState<ComponentType<Record<string, unknown>> | null>(null);
      useEffect(() => {
        void loader().then((component) => setInner(() => component));
      }, []);
      return Inner ? <Inner {...props} /> : null;
    },
}));

// CodeMirror itself needs a real browser layout; the wiring (theme, extensions, options) is what is checked here.
vi.mock("@uiw/react-codemirror", () => ({
  default: (props: CodeMirrorProps) => {
    codeMirrorCalls.push(props);
    return <div data-testid="codemirror">{props.value}</div>;
  },
}));

import { Editor } from "@/components/Editor";

function lastCall(): CodeMirrorProps {
  return codeMirrorCalls[codeMirrorCalls.length - 1];
}

describe("Editor (loaded CodeMirror)", () => {
  beforeEach(() => {
    codeMirrorCalls.length = 0;
    localStorage.clear();
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it("passes the value, the language extension and the dark theme by default", async () => {
    render(
      <ThemeProvider>
        <Editor value="package main" language="go" label="Go input" id="in" onChange={vi.fn()} placeholder="type" />
      </ThemeProvider>
    );
    await waitFor(() => expect(screen.getByTestId("codemirror")).toHaveTextContent("package main"));
    const props = lastCall();
    expect(props.extensions).toHaveLength(1);
    expect(props.theme).toBe(intellijTheme);
    expect(props.placeholder).toBe("type");
    expect(props.basicSetup).toEqual({ lineNumbers: true, foldGutter: false, highlightActiveLine: true });
  });

  it("uses the light theme when the stored mode is light", async () => {
    localStorage.setItem(THEME_STORAGE_KEY, "light");
    render(
      <ThemeProvider>
        <Editor value="SELECT 1;" language="sql" label="SQL output" id="out" readOnly />
      </ThemeProvider>
    );
    await waitFor(() => expect(screen.getByTestId("codemirror")).toBeInTheDocument());
    const props = lastCall();
    expect(props.theme).toBe(githubLight);
    expect(props.readOnly).toBe(true);
    expect(props.basicSetup.highlightActiveLine).toBe(false);
  });

  it("builds a different extension for sql than for go", async () => {
    const { rerender } = render(
      <ThemeProvider>
        <Editor value="a" language="go" label="x" id="a" />
      </ThemeProvider>
    );
    await waitFor(() => expect(screen.getByTestId("codemirror")).toBeInTheDocument());
    const goExtension = lastCall().extensions[0];

    codeMirrorCalls.length = 0;
    rerender(
      <ThemeProvider>
        <Editor value="a" language="sql" label="x" id="a" />
      </ThemeProvider>
    );
    await waitFor(() => expect(codeMirrorCalls.length).toBeGreaterThan(0));
    expect(lastCall().extensions[0]).not.toBe(goExtension);
  });
});

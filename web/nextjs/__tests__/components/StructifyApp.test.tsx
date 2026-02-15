import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";

// Mock next/navigation
vi.mock("next/navigation", () => ({
  useRouter: () => ({ replace: vi.fn() }),
  useSearchParams: () => new URLSearchParams(),
}));

// Mock the Editor since CodeMirror doesn't work in jsdom
vi.mock("@/components/Editor", () => ({
  Editor: ({
    value,
    onChange,
    placeholder,
    id,
  }: {
    value: string;
    onChange?: (v: string) => void;
    placeholder?: string;
    id: string;
  }) => (
    <textarea
      data-testid={id}
      value={value}
      placeholder={placeholder}
      onChange={(e) => onChange?.(e.target.value)}
      readOnly={!onChange}
    />
  ),
}));

// Mock API
vi.mock("@/lib/api", () => ({
  generateSQL: vi.fn(),
  generateCode: vi.fn(),
  fetchVersion: vi.fn().mockResolvedValue("0.1.0"),
}));

import { StructifyApp } from "@/components/StructifyApp";
import { generateSQL, generateCode } from "@/lib/api";

beforeEach(() => {
  vi.clearAllMocks();
  globalThis.localStorage?.clear?.();
});

describe("StructifyApp", () => {
  it("renders header with title", async () => {
    render(<StructifyApp />);
    expect(screen.getByText("structify")).toBeInTheDocument();
  });

  it("shows version from API", async () => {
    render(<StructifyApp />);
    await waitFor(() => {
      expect(screen.getByText("v0.1.0")).toBeInTheDocument();
    });
  });

  it("renders generate button", () => {
    render(<StructifyApp />);
    expect(screen.getByText("Generate")).toBeInTheDocument();
  });

  it("disables generate button when input is empty", () => {
    render(<StructifyApp />);
    expect(
      screen.getByLabelText("Generate output (Ctrl+Enter)")
    ).toBeDisabled();
  });

  it("enables generate button when input has text", async () => {
    const user = userEvent.setup();
    render(<StructifyApp />);

    const input = screen.getByTestId("input-editor");
    await user.type(input, "package main");

    expect(
      screen.getByLabelText("Generate output (Ctrl+Enter)")
    ).toBeEnabled();
  });

  it("calls generateSQL when generate is clicked in sql mode", async () => {
    const user = userEvent.setup();
    vi.mocked(generateSQL).mockResolvedValue("CREATE TABLE test ();");

    render(<StructifyApp />);

    const input = screen.getByTestId("input-editor");
    await user.type(input, 'package m\ntype T struct {{\n\tID int64 `db:"pk"`\n}');

    await user.click(screen.getByText("Generate"));

    await waitFor(() => {
      expect(generateSQL).toHaveBeenCalled();
    });
  });

  it("shows error when generation fails", async () => {
    const user = userEvent.setup();
    vi.mocked(generateSQL).mockRejectedValue(new Error("parse failed"));

    render(<StructifyApp />);

    const input = screen.getByTestId("input-editor");
    await user.type(input, "bad input");
    await user.click(screen.getByText("Generate"));

    await waitFor(() => {
      expect(screen.getByText("parse failed")).toBeInTheDocument();
    });
  });

  it("clears output when mode changes", async () => {
    const user = userEvent.setup();
    vi.mocked(generateSQL).mockResolvedValue("SQL output");

    render(<StructifyApp />);

    const input = screen.getByTestId("input-editor");
    await user.type(input, "test");
    await user.click(screen.getByText("Generate"));

    await waitFor(() => {
      expect(generateSQL).toHaveBeenCalled();
    });

    await user.click(screen.getByText("Repository Code"));
    expect(screen.getByText("Output will appear here")).toBeInTheDocument();
  });

  it("persists input to localStorage", async () => {
    const store: Record<string, string> = {};
    vi.stubGlobal("localStorage", {
      getItem: vi.fn((key: string) => store[key] ?? null),
      setItem: vi.fn((key: string, val: string) => { store[key] = val; }),
      removeItem: vi.fn((key: string) => { delete store[key]; }),
      clear: vi.fn(),
    });

    const user = userEvent.setup();
    render(<StructifyApp />);

    const input = screen.getByTestId("input-editor");
    await user.type(input, "stored text");

    expect(store["structify_input"]).toBe("stored text");
  });

  it("renders GitHub link", () => {
    render(<StructifyApp />);
    expect(screen.getByText("GitHub")).toHaveAttribute(
      "href",
      "https://github.com/n0xum/structify"
    );
  });
});

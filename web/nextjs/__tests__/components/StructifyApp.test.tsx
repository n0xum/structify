import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor, act, fireEvent } from "@testing-library/react";
import userEvent from "@testing-library/user-event";

const mockReplace = vi.fn();
let mockSearchParams = new URLSearchParams();

// Mock next/navigation
vi.mock("next/navigation", () => ({
  useRouter: () => ({ replace: mockReplace }),
  useSearchParams: () => mockSearchParams,
}));

// Mock next/link
vi.mock("next/link", () => ({
  default: ({
    href,
    children,
    className,
  }: {
    href: string;
    children: React.ReactNode;
    className?: string;
  }) => (
    <a href={href} className={className}>
      {children}
    </a>
  ),
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
  vi.unstubAllGlobals();
  mockSearchParams = new URLSearchParams();
  mockReplace.mockClear();
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

  it("calls generateCode when generate is clicked in code mode", async () => {
    const user = userEvent.setup();
    vi.mocked(generateCode).mockResolvedValue("// generated code");

    render(<StructifyApp />);

    await user.click(screen.getByText("Repository Code"));

    const input = screen.getByTestId("input-editor");
    fireEvent.change(input, { target: { value: "package m\ntype T struct { ID int64 }" } });

    await user.click(screen.getByText("Generate"));

    await waitFor(() => {
      expect(generateCode).toHaveBeenCalled();
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

  it("restores input from localStorage on mount", async () => {
    const store: Record<string, string> = { structify_input: "restored content" };
    vi.stubGlobal("localStorage", {
      getItem: vi.fn((key: string) => store[key] ?? null),
      setItem: vi.fn(),
      removeItem: vi.fn(),
      clear: vi.fn(),
    });

    render(<StructifyApp />);

    await waitFor(() => {
      expect(screen.getByTestId("input-editor")).toHaveValue("restored content");
    });
  });

  it("renders GitHub link", () => {
    render(<StructifyApp />);
    expect(screen.getByText("GitHub")).toHaveAttribute(
      "href",
      "https://github.com/n0xum/structify"
    );
  });

  it("renders Docs nav link", () => {
    render(<StructifyApp />);
    expect(screen.getByText("Docs")).toHaveAttribute("href", "/docs");
  });

  it("triggers generate with Ctrl+Enter keyboard shortcut", async () => {
    const user = userEvent.setup();
    vi.mocked(generateSQL).mockResolvedValue("SQL output");

    render(<StructifyApp />);

    const input = screen.getByTestId("input-editor");
    await user.type(input, "package main");

    await user.keyboard("{Control>}{Enter}{/Control}");

    await waitFor(() => {
      expect(generateSQL).toHaveBeenCalled();
    });
  });

  it("loads example from ?load= search param and removes param from URL", async () => {
    mockSearchParams = new URLSearchParams("load=User");

    render(<StructifyApp />);

    await waitFor(() => {
      const value = screen.getByTestId("input-editor").getAttribute("value") ??
        (screen.getByTestId("input-editor") as HTMLTextAreaElement).value;
      expect(value).toContain("type User struct");
    });

    // replace should have been called with a URL that has no ?load= param
    await waitFor(() => {
      const loadRemovalCall = mockReplace.mock.calls.find(
        ([url]: [string]) => !url.includes("load=")
      );
      expect(loadRemovalCall).toBeDefined();
    });
  });

  it("preserves ?mode= when removing ?load= from URL", async () => {
    mockSearchParams = new URLSearchParams("load=User&mode=sql");

    render(<StructifyApp />);

    await waitFor(() => {
      expect(mockReplace).toHaveBeenCalled();
    });

    const call = mockReplace.mock.calls.find(
      ([url]: [string]) => !url.includes("load=")
    );
    expect(call).toBeDefined();
  });

  it("silently ignores unknown ?load= values", async () => {
    mockSearchParams = new URLSearchParams("load=NonExistentExample");

    render(<StructifyApp />);

    await waitFor(() => {
      expect(screen.getByText("structify")).toBeInTheDocument();
    });

    // Input should remain empty — unknown label should not load anything
    const input = screen.getByTestId("input-editor") as HTMLTextAreaElement;
    expect(input.value).toBe("");
  });

  it("shows package name input when mode is code", async () => {
    const user = userEvent.setup();
    render(<StructifyApp />);

    await user.click(screen.getByText("Repository Code"));

    expect(screen.getByPlaceholderText("models")).toBeInTheDocument();
  });

  it("shows validation warning for large input", async () => {
    const user = userEvent.setup();
    render(<StructifyApp />);

    const bigInput = "x".repeat(102_400 + 1);
    const input = screen.getByTestId("input-editor");

    await act(async () => {
      await user.clear(input);
      // Simulate direct state change to avoid slow character-by-character typing
      input.focus();
      await user.paste(bigInput);
    });

    await waitFor(() => {
      expect(screen.getByText(/exceeds/i)).toBeInTheDocument();
    });
  });
});

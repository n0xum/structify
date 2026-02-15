import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { OutputPanel } from "@/components/OutputPanel";

// Mock the Editor since CodeMirror doesn't work in jsdom
vi.mock("@/components/Editor", () => ({
  Editor: ({ value, id }: { value: string; id: string }) => (
    <pre data-testid={id}>{value}</pre>
  ),
}));

describe("OutputPanel", () => {
  it("shows placeholder when output is empty", () => {
    render(<OutputPanel output="" mode="sql" />);
    expect(screen.getByText("Output will appear here")).toBeInTheDocument();
  });

  it("renders output in editor when provided", () => {
    render(<OutputPanel output="CREATE TABLE users ();" mode="sql" />);
    expect(screen.getByTestId("output-editor")).toHaveTextContent(
      "CREATE TABLE users ();"
    );
  });

  it("shows copy button when output exists", () => {
    render(<OutputPanel output="some output" mode="sql" />);
    expect(screen.getByLabelText("Copy output to clipboard")).toBeInTheDocument();
  });

  it("hides copy button when output is empty", () => {
    render(<OutputPanel output="" mode="sql" />);
    expect(
      screen.queryByLabelText("Copy output to clipboard")
    ).not.toBeInTheDocument();
  });

  it("copies output to clipboard on click", async () => {
    const user = userEvent.setup();
    const writeText = vi.fn().mockResolvedValue(undefined);
    Object.defineProperty(navigator, "clipboard", {
      value: { writeText },
      writable: true,
      configurable: true,
    });

    render(<OutputPanel output="test output" mode="sql" />);
    await user.click(screen.getByLabelText("Copy output to clipboard"));

    expect(writeText).toHaveBeenCalledWith("test output");
    expect(screen.getByText("Copied!")).toBeInTheDocument();
  });
});

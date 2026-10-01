import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { ModeSelector } from "@/components/ModeSelector";

describe("ModeSelector", () => {
  const defaultProps = {
    mode: "sql" as const,
    onChange: vi.fn(),
    packageName: "models",
    onPackageChange: vi.fn(),
    packageError: null,
  };

  it("renders SQL and Repo tabs", () => {
    render(<ModeSelector {...defaultProps} />);
    expect(screen.getByText("SQL Schema")).toBeInTheDocument();
    expect(screen.getByText("Interface Repository")).toBeInTheDocument();
  });

  it("calls onChange when clicking repo tab", async () => {
    const user = userEvent.setup();
    const onChange = vi.fn();
    render(<ModeSelector {...defaultProps} onChange={onChange} />);

    await user.click(screen.getByText("Interface Repository"));
    expect(onChange).toHaveBeenCalledWith("repo");
  });

  it("shows package input when mode is repo", () => {
    render(<ModeSelector {...defaultProps} mode="repo" />);
    expect(screen.getByLabelText("Package name")).toBeInTheDocument();
  });

  it("hides package input when mode is sql", () => {
    render(<ModeSelector {...defaultProps} mode="sql" />);
    expect(screen.queryByLabelText("Package name")).not.toBeInTheDocument();
  });

  it("shows package error when provided", () => {
    render(
      <ModeSelector {...defaultProps} mode="repo" packageError="Invalid name" />
    );
    expect(screen.getByText("Invalid name")).toBeInTheDocument();
  });

  it("calls onChange when clicking the SQL tab", async () => {
    const user = userEvent.setup();
    const onChange = vi.fn();
    render(<ModeSelector {...defaultProps} mode="repo" onChange={onChange} />);

    await user.click(screen.getByText("SQL Schema"));
    expect(onChange).toHaveBeenCalledWith("sql");
  });

  it("marks only the active tab as selected", () => {
    render(<ModeSelector {...defaultProps} mode="repo" />);
    expect(screen.getByRole("tab", { name: "Interface Repository" })).toHaveAttribute("aria-selected", "true");
    expect(screen.getByRole("tab", { name: "SQL Schema" })).toHaveAttribute("aria-selected", "false");
  });

  it("reports every change of the package name", async () => {
    const user = userEvent.setup();
    const onPackageChange = vi.fn();
    render(<ModeSelector {...defaultProps} mode="repo" packageName="" onPackageChange={onPackageChange} />);

    await user.type(screen.getByLabelText("Package name"), "db");
    expect(onPackageChange).toHaveBeenCalledTimes(2);
    expect(onPackageChange).toHaveBeenLastCalledWith("b");
  });

  it("links the package input to its error message", () => {
    render(<ModeSelector {...defaultProps} mode="repo" packageError="Invalid name" />);
    expect(screen.getByLabelText("Package name")).toHaveAttribute("aria-describedby", "pkg-error");
  });

  it("has no error description when there is no error", () => {
    render(<ModeSelector {...defaultProps} mode="repo" />);
    expect(screen.getByLabelText("Package name")).not.toHaveAttribute("aria-describedby");
  });
});

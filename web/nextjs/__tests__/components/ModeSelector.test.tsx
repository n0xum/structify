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

  it("renders SQL and Code tabs", () => {
    render(<ModeSelector {...defaultProps} />);
    expect(screen.getByText("SQL Schema")).toBeInTheDocument();
    expect(screen.getByText("Repository Code")).toBeInTheDocument();
  });

  it("calls onChange when clicking code tab", async () => {
    const user = userEvent.setup();
    const onChange = vi.fn();
    render(<ModeSelector {...defaultProps} onChange={onChange} />);

    await user.click(screen.getByText("Repository Code"));
    expect(onChange).toHaveBeenCalledWith("code");
  });

  it("shows package input when mode is code", () => {
    render(<ModeSelector {...defaultProps} mode="code" />);
    expect(screen.getByLabelText("Package name")).toBeInTheDocument();
  });

  it("hides package input when mode is sql", () => {
    render(<ModeSelector {...defaultProps} mode="sql" />);
    expect(screen.queryByLabelText("Package name")).not.toBeInTheDocument();
  });

  it("shows package error when provided", () => {
    render(
      <ModeSelector {...defaultProps} mode="code" packageError="Invalid name" />
    );
    expect(screen.getByText("Invalid name")).toBeInTheDocument();
  });
});

import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { TagReference } from "@/components/TagReference";

describe("TagReference", () => {
  it("renders toggle button", () => {
    render(<TagReference />);
    expect(screen.getByLabelText("Show db tag reference")).toBeInTheDocument();
  });

  it("does not show tooltip initially", () => {
    render(<TagReference />);
    expect(screen.queryByRole("tooltip")).not.toBeInTheDocument();
  });

  it("shows tooltip on click", async () => {
    const user = userEvent.setup();
    render(<TagReference />);

    await user.click(screen.getByLabelText("Show db tag reference"));
    expect(screen.getByRole("tooltip")).toBeInTheDocument();
    expect(screen.getByText('db:"pk"')).toBeInTheDocument();
  });

  it("hides tooltip on second click", async () => {
    const user = userEvent.setup();
    render(<TagReference />);

    const btn = screen.getByLabelText("Show db tag reference");
    await user.click(btn);
    expect(screen.getByRole("tooltip")).toBeInTheDocument();

    await user.click(btn);
    expect(screen.queryByRole("tooltip")).not.toBeInTheDocument();
  });
});

import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { ExampleLoader } from "@/components/ExampleLoader";
import { EXAMPLES } from "@/lib/examples";

describe("ExampleLoader", () => {
  it("renders a select element", () => {
    render(<ExampleLoader onSelect={() => {}} />);
    expect(screen.getByLabelText("Load example struct")).toBeInTheDocument();
  });

  it("lists all examples as options", () => {
    render(<ExampleLoader onSelect={() => {}} />);
    for (const ex of EXAMPLES) {
      expect(screen.getByText(new RegExp(ex.label))).toBeInTheDocument();
    }
  });

  it("calls onSelect with example source when selected", async () => {
    const user = userEvent.setup();
    const onSelect = vi.fn();
    render(<ExampleLoader onSelect={onSelect} />);

    const select = screen.getByLabelText("Load example struct");
    await user.selectOptions(select, EXAMPLES[0].label);

    expect(onSelect).toHaveBeenCalledWith(EXAMPLES[0].source);
  });
});

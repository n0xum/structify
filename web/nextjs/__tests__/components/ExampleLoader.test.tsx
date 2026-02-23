import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { ExampleLoader } from "@/components/ExampleLoader";
import { EXAMPLES, REPO_EXAMPLES } from "@/lib/examples";

describe("ExampleLoader — sql mode", () => {
  it("renders a select element with sql aria-label", () => {
    render(<ExampleLoader mode="sql" onSelect={() => {}} />);
    expect(screen.getByLabelText("Load example struct")).toBeInTheDocument();
  });

  it("lists all SQL examples as options", () => {
    render(<ExampleLoader mode="sql" onSelect={() => {}} />);
    for (const ex of EXAMPLES) {
      expect(screen.getByText(new RegExp(ex.label))).toBeInTheDocument();
    }
  });

  it("does not show repo examples in sql mode", () => {
    render(<ExampleLoader mode="sql" onSelect={() => {}} />);
    // Repo examples have distinct labels not present in SQL examples
    for (const ex of REPO_EXAMPLES) {
      const inSql = EXAMPLES.some((e) => e.label === ex.label);
      if (!inSql) {
        expect(screen.queryByText(new RegExp(ex.label))).not.toBeInTheDocument();
      }
    }
  });

  it("calls onSelect with SQL example source when selected", async () => {
    const user = userEvent.setup();
    const onSelect = vi.fn();
    render(<ExampleLoader mode="sql" onSelect={onSelect} />);

    const select = screen.getByLabelText("Load example struct");
    await user.selectOptions(select, EXAMPLES[0].label);

    expect(onSelect).toHaveBeenCalledWith(EXAMPLES[0].source);
  });
});

describe("ExampleLoader — repo mode", () => {
  it("renders a select element with interface aria-label", () => {
    render(<ExampleLoader mode="repo" onSelect={() => {}} />);
    expect(screen.getByLabelText("Load example interface")).toBeInTheDocument();
  });

  it("lists all repo examples as options", () => {
    render(<ExampleLoader mode="repo" onSelect={() => {}} />);
    for (const ex of REPO_EXAMPLES) {
      // Use exact string match to avoid regex special-char issues in labels
      expect(screen.getByText(`${ex.label} — ${ex.description}`)).toBeInTheDocument();
    }
  });

  it("calls onSelect with repo example source when selected", async () => {
    const user = userEvent.setup();
    const onSelect = vi.fn();
    render(<ExampleLoader mode="repo" onSelect={onSelect} />);

    const select = screen.getByLabelText("Load example interface");
    await user.selectOptions(select, REPO_EXAMPLES[0].label);

    expect(onSelect).toHaveBeenCalledWith(REPO_EXAMPLES[0].source);
  });
});

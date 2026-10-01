import { describe, it, expect } from "vitest";
import { render } from "@testing-library/react";
import { HighlightedCode } from "@/components/HighlightedCode";

describe("HighlightedCode", () => {
  it("renders Go code and keeps its text", () => {
    const { container } = render(<HighlightedCode code={"func Add(a, b int) int {\n\treturn a + b\n}"} language="go" />);
    expect(container.textContent).toContain("func Add(a, b int) int");
    expect(container.querySelector("pre.syntax-highlight")).toBeInTheDocument();
  });

  it("renders SQL code and keeps its text", () => {
    const { container } = render(<HighlightedCode code="SELECT id FROM users WHERE id = 1;" language="sql" />);
    expect(container.textContent).toContain("SELECT id FROM users WHERE id = 1;");
  });

  it("highlights keywords with inline styles for a registered language", () => {
    const { container } = render(<HighlightedCode code="SELECT 1;" language="sql" />);
    const styled = container.querySelectorAll("span[style]");
    expect(styled.length).toBeGreaterThan(0);
  });

  it("falls back to plain text for a language that is not registered", () => {
    const { container } = render(<HighlightedCode code="echo hello" language="bash" />);
    expect(container.textContent).toContain("echo hello");
  });

  it("uses a transparent background without margin or padding", () => {
    const { container } = render(<HighlightedCode code="x" language="go" />);
    const pre = container.querySelector("pre") as HTMLElement;
    expect(pre.style.margin).toBe("0px");
    expect(pre.style.padding).toBe("0px");
  });
});

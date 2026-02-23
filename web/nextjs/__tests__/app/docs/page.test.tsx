import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";

// next/link renders a plain <a> in tests
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

// HighlightedCode uses react-syntax-highlighter which doesn't work in jsdom
vi.mock("@/components/HighlightedCode", () => ({
  HighlightedCode: ({ code }: { code: string }) => <pre>{code}</pre>,
}));

import DocsPage from "@/app/docs/page";

describe("DocsPage", () => {
  it("renders the page heading", () => {
    render(<DocsPage />);
    expect(screen.getByText("Docs & Features")).toBeInTheDocument();
  });

  it("renders the db: tag description", () => {
    render(<DocsPage />);
    expect(screen.getAllByText(/db:/).length).toBeGreaterThanOrEqual(1);
  });

  it("renders links back to the editor", () => {
    render(<DocsPage />);
    const links = screen.getAllByText(/Back to Editor/);
    expect(links.length).toBeGreaterThanOrEqual(1);
    expect(links[0].closest("a")).toHaveAttribute("href", "/");
  });

  // Current sections in the docs page
  const sections = [
    { title: "Basic Mapping", id: "basic", label: "User" },
    { title: "Constraints", id: "constraints", label: "Constraints" },
    { title: "Indexes", id: "indexes", label: "Indexes" },
    { title: "Foreign Keys", id: "foreign-keys", label: "Foreign Keys" },
    { title: "Repository Generation", id: "repo-generation", label: "Indexes" },
  ];

  it("renders all 5 feature section headings", () => {
    render(<DocsPage />);
    for (const s of sections) {
      expect(screen.getAllByText(s.title).length).toBeGreaterThanOrEqual(1);
    }
  });

  it("renders sidebar anchor links for each section", () => {
    render(<DocsPage />);
    for (const s of sections) {
      const anchors = screen
        .getAllByRole("link")
        .filter((a) => a.getAttribute("href") === `#${s.id}`);
      expect(anchors.length).toBeGreaterThanOrEqual(1);
    }
  });

  it("renders Try it links for SQL sections pointing to /?load=&mode=sql", () => {
    render(<DocsPage />);
    const sqlSections = sections.filter((s) => s.id !== "repo-generation");
    for (const s of sqlSections) {
      const expected = `/?load=${encodeURIComponent(s.label)}&mode=sql`;
      const links = screen
        .getAllByRole("link")
        .filter((a) => a.getAttribute("href") === expected);
      expect(links.length).toBe(1);
    }
  });

  it("renders each section as a scrollable element with correct id", () => {
    const { container } = render(<DocsPage />);
    for (const s of sections) {
      expect(container.querySelector(`#${s.id}`)).toBeInTheDocument();
    }
  });

  it("renders Go code for the Basic section", () => {
    const { container } = render(<DocsPage />);
    const basicSection = container.querySelector("#basic");
    expect(basicSection?.textContent).toContain('db:"pk"');
  });

  it("renders SQL code for the Basic section", () => {
    const { container } = render(<DocsPage />);
    const basicSection = container.querySelector("#basic");
    expect(basicSection?.textContent).toContain("CREATE TABLE users");
  });

  it("renders CHECK constraint code for Constraints section", () => {
    const { container } = render(<DocsPage />);
    const section = container.querySelector("#constraints");
    expect(section?.textContent).toContain("CHECK (length(name) > 0)");
  });

  it("renders CREATE INDEX for Indexes section", () => {
    const { container } = render(<DocsPage />);
    const section = container.querySelector("#indexes");
    expect(section?.textContent).toContain("CREATE UNIQUE INDEX uq_slug");
  });

  it("renders FOREIGN KEY for Foreign Keys section", () => {
    const { container } = render(<DocsPage />);
    const section = container.querySelector("#foreign-keys");
    expect(section?.textContent).toContain("FOREIGN KEY (user_id)");
  });

  it("renders generated repository code in repo-generation section", () => {
    const { container } = render(<DocsPage />);
    const section = container.querySelector("#repo-generation");
    expect(section?.textContent).toContain("FindByTitleAndLanguage");
  });

  it("repo-generation section shows the --to-repo mention", () => {
    render(<DocsPage />);
    expect(screen.getAllByText(/--to-repo/).length).toBeGreaterThanOrEqual(1);
  });
});

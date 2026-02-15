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

import DocsPage from "@/app/docs/page";

describe("DocsPage", () => {
  it("renders the page title", () => {
    render(<DocsPage />);
    expect(screen.getByText("Docs")).toBeInTheDocument();
  });

  it("renders the db: tag description", () => {
    render(<DocsPage />);
    expect(screen.getAllByText(/db:/).length).toBeGreaterThanOrEqual(1);
  });

  it("renders a link back to the home page", () => {
    render(<DocsPage />);
    const links = screen.getAllByText("← structify");
    expect(links.length).toBeGreaterThanOrEqual(1);
    expect(links[0]).toHaveAttribute("href", "/");
  });

  // Sidebar nav links
  const sections = [
    { title: "Basic", id: "basic", label: "User" },
    { title: "Constraints", id: "constraints", label: "Constraints" },
    { title: "Indexes", id: "indexes", label: "Indexes" },
    { title: "Foreign Keys", id: "foreign-keys", label: "Foreign Keys" },
    { title: "Composite Keys", id: "composite-keys", label: "Composite PK & FK" },
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

  it("renders a Try it link for each section pointing to /?load=", () => {
    render(<DocsPage />);
    for (const s of sections) {
      const expected = `/?load=${encodeURIComponent(s.label)}`;
      const links = screen
        .getAllByRole("link")
        .filter((a) => a.getAttribute("href") === expected);
      expect(links.length).toBe(1);
    }
  });

  it("renders Go and SQL code block labels for each section", () => {
    render(<DocsPage />);
    const goLabels = screen.getAllByText("Go");
    const sqlLabels = screen.getAllByText("SQL");
    expect(goLabels.length).toBe(sections.length);
    expect(sqlLabels.length).toBe(sections.length);
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
    render(<DocsPage />);
    expect(screen.getByText(/CHECK \(length\(name\) > 0\)/)).toBeInTheDocument();
  });

  it("renders CREATE INDEX for Indexes section", () => {
    render(<DocsPage />);
    expect(screen.getByText(/CREATE UNIQUE INDEX uq_slug/)).toBeInTheDocument();
  });

  it("renders FOREIGN KEY for Foreign Keys section", () => {
    render(<DocsPage />);
    expect(screen.getByText(/FOREIGN KEY \(user_id\)/)).toBeInTheDocument();
  });

  it("renders composite PRIMARY KEY for Composite Keys section", () => {
    render(<DocsPage />);
    expect(screen.getByText(/PRIMARY KEY \(order_id, product_id\)/)).toBeInTheDocument();
  });

  it("each section element has a scroll-mt-8 class", () => {
    const { container } = render(<DocsPage />);
    const sectionEls = container.querySelectorAll("section.scroll-mt-8");
    expect(sectionEls.length).toBe(sections.length);
  });

  it("each section has the correct id attribute", () => {
    const { container } = render(<DocsPage />);
    for (const s of sections) {
      expect(container.querySelector(`#${s.id}`)).toBeInTheDocument();
    }
  });
});

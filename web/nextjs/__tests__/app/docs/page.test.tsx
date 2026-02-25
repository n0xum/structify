import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";

const { redirectMock } = vi.hoisted(() => ({
  redirectMock: vi.fn(),
}));

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

vi.mock("next/navigation", () => ({
  usePathname: () => "/docs/sql-schema",
  redirect: redirectMock,
}));

vi.mock("@/components/HighlightedCode", () => ({
  HighlightedCode: ({ code }: { code: string }) => <pre>{code}</pre>,
}));

import DocsLayout from "@/app/docs/layout";
import DocsPage from "@/app/docs/page";
import DocsSqlSchemaPage from "@/app/docs/sql-schema/page";
import DocsRepositoryPage from "@/app/docs/repository/page";
import DocsProjectGuidePage from "@/app/docs/project-guide/page";

describe("Split docs routes", () => {
  it("redirects /docs to /docs/sql-schema", () => {
    DocsPage();

    expect(redirectMock).toHaveBeenCalledWith("/docs/sql-schema");
  });

  it("renders sidebar links to all split routes", () => {
    render(
      <DocsLayout>
        <DocsSqlSchemaPage />
      </DocsLayout>,
    );

    expect(screen.getAllByText("SQL Schema").length).toBeGreaterThanOrEqual(1);
    expect(screen.getAllByText("Repository").length).toBeGreaterThanOrEqual(1);
    expect(screen.getAllByText("Project Guide").length).toBeGreaterThanOrEqual(1);

    expect(screen.getByRole("link", { name: "SQL Schema" })).toHaveAttribute("href", "/docs/sql-schema");
    expect(screen.getByRole("link", { name: "Repository" })).toHaveAttribute("href", "/docs/repository");
    expect(screen.getByRole("link", { name: "Project Guide" })).toHaveAttribute("href", "/docs/project-guide");

    expect(
      screen
        .getAllByRole("link")
        .some((a) => a.getAttribute("href") === "/docs/sql-schema#sql-struct-basics"),
    ).toBe(true);
    expect(
      screen
        .getAllByRole("link")
        .some((a) => a.getAttribute("href") === "/docs/repository#repo-prerequisites"),
    ).toBe(true);
    expect(
      screen
        .getAllByRole("link")
        .some((a) => a.getAttribute("href") === "/docs/project-guide#project-feature-overview"),
    ).toBe(true);
  });

  it("renders SQL schema page content only", () => {
    const { container } = render(<DocsSqlSchemaPage />);

    expect(screen.getByText("Docs & Features")).toBeInTheDocument();
    expect(screen.getByRole("heading", { level: 2, name: "SQL Schema" })).toBeInTheDocument();
    expect(container.querySelector("#sql-struct-basics")).toBeInTheDocument();
    expect(container.querySelector("#repo-prerequisites")).not.toBeInTheDocument();
  });

  it("renders repository page content only", () => {
    const { container } = render(<DocsRepositoryPage />);

    expect(screen.getByRole("heading", { level: 2, name: "Repository" })).toBeInTheDocument();
    expect(container.querySelector("#repo-prerequisites")).toBeInTheDocument();
    expect(container.querySelector("#project-feature-overview")).not.toBeInTheDocument();
    expect(screen.getAllByText(/--to-repo/).length).toBeGreaterThanOrEqual(1);
  });

  it("renders project guide page content only", () => {
    const { container } = render(<DocsProjectGuidePage />);

    expect(screen.getByRole("heading", { level: 2, name: "Project Guide" })).toBeInTheDocument();
    expect(container.querySelector("#project-local-commands")).toBeInTheDocument();
    expect(container.querySelector("#sql-struct-basics")).not.toBeInTheDocument();
    expect(container.querySelector("#project-troubleshooting")?.textContent).toContain("Quick debug loop");
  });
});

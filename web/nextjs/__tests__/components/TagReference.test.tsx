import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { TagReference } from "@/components/TagReference";

describe("TagReference — sql mode", () => {
  it("renders toggle button with db tag label", () => {
    render(<TagReference mode="sql" />);
    expect(screen.getByLabelText("Show db tag reference")).toBeInTheDocument();
  });

  it("does not show tooltip initially", () => {
    render(<TagReference mode="sql" />);
    expect(screen.queryByRole("tooltip")).not.toBeInTheDocument();
  });

  it("shows tooltip on click", async () => {
    const user = userEvent.setup();
    render(<TagReference mode="sql" />);

    await user.click(screen.getByLabelText("Show db tag reference"));
    expect(screen.getByRole("tooltip")).toBeInTheDocument();
    expect(screen.getByText('db:"pk"')).toBeInTheDocument();
  });

  it("hides tooltip on second click", async () => {
    const user = userEvent.setup();
    render(<TagReference mode="sql" />);

    const btn = screen.getByLabelText("Show db tag reference");
    await user.click(btn);
    expect(screen.getByRole("tooltip")).toBeInTheDocument();

    await user.click(btn);
    expect(screen.queryByRole("tooltip")).not.toBeInTheDocument();
  });

  it("shows all tag groups when open", async () => {
    const user = userEvent.setup();
    render(<TagReference mode="sql" />);

    await user.click(screen.getByLabelText("Show db tag reference"));
    expect(screen.getByText("Constraints")).toBeInTheDocument();
    expect(screen.getByText("Indexes")).toBeInTheDocument();
    expect(screen.getByText("Foreign Keys")).toBeInTheDocument();
    expect(screen.getByText("Composite")).toBeInTheDocument();
  });

  it("shows check constraint tag", async () => {
    const user = userEvent.setup();
    render(<TagReference mode="sql" />);

    await user.click(screen.getByLabelText("Show db tag reference"));
    expect(screen.getByText('db:"check:age >= 18"')).toBeInTheDocument();
  });

  it("shows foreign key tag", async () => {
    const user = userEvent.setup();
    render(<TagReference mode="sql" />);

    await user.click(screen.getByLabelText("Show db tag reference"));
    expect(screen.getByText('db:"fk:users,id"')).toBeInTheDocument();
  });

  it("shows enum tag", async () => {
    const user = userEvent.setup();
    render(<TagReference mode="sql" />);

    await user.click(screen.getByLabelText("Show db tag reference"));
    expect(screen.getByText('db:"enum:a,b,c"')).toBeInTheDocument();
  });
});

describe("TagReference — repo mode", () => {
  it("renders toggle button with method naming label", () => {
    render(<TagReference mode="repo" />);
    expect(screen.getByLabelText("Show method naming reference")).toBeInTheDocument();
  });

  it("shows method naming reference heading when open", async () => {
    const user = userEvent.setup();
    render(<TagReference mode="repo" />);

    await user.click(screen.getByLabelText("Show method naming reference"));
    expect(screen.getByText("method naming reference")).toBeInTheDocument();
  });

  it("shows SmartQuery group", async () => {
    const user = userEvent.setup();
    render(<TagReference mode="repo" />);

    await user.click(screen.getByLabelText("Show method naming reference"));
    expect(screen.getByText("SmartQuery methods")).toBeInTheDocument();
  });

  it("shows FindBy group", async () => {
    const user = userEvent.setup();
    render(<TagReference mode="repo" />);

    await user.click(screen.getByLabelText("Show method naming reference"));
    expect(screen.getByText("FindBy filters")).toBeInTheDocument();
  });

  it("shows custom SQL override group", async () => {
    const user = userEvent.setup();
    render(<TagReference mode="repo" />);

    await user.click(screen.getByLabelText("Show method naming reference"));
    expect(screen.getByText("Custom SQL override")).toBeInTheDocument();
  });

  it("shows repository method signature guidance", async () => {
    const user = userEvent.setup();
    render(<TagReference mode="repo" />);

    await user.click(screen.getByLabelText("Show method naming reference"));
    expect(screen.getByText(/context\.Context/)).toBeInTheDocument();
  });

  it("does not show db tag content", async () => {
    const user = userEvent.setup();
    render(<TagReference mode="repo" />);

    await user.click(screen.getByLabelText("Show method naming reference"));
    expect(screen.queryByText('db:"pk"')).not.toBeInTheDocument();
  });
});

import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import GlobalError from "@/app/global-error";

describe("GlobalError", () => {
  it("renders error message", () => {
    render(
      <GlobalError error={new Error("boom")} reset={() => {}} />
    );
    expect(screen.getByText("Something went wrong")).toBeInTheDocument();
  });

  it("renders try again button", () => {
    render(
      <GlobalError error={new Error("boom")} reset={() => {}} />
    );
    expect(screen.getByText("Try again")).toBeInTheDocument();
  });

  it("calls reset on button click", async () => {
    const user = userEvent.setup();
    const reset = vi.fn();
    render(<GlobalError error={new Error("boom")} reset={reset} />);

    await user.click(screen.getByText("Try again"));
    expect(reset).toHaveBeenCalledOnce();
  });

  it("links to GitHub issues", () => {
    render(
      <GlobalError error={new Error("boom")} reset={() => {}} />
    );
    expect(screen.getByText("open an issue")).toHaveAttribute(
      "href",
      "https://github.com/n0xum/structify/issues"
    );
  });
});

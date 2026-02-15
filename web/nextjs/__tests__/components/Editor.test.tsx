import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";

// Mock the dynamic CodeMirror import
vi.mock("next/dynamic", () => ({
  default: (_loader: unknown, _opts: unknown) =>
    function MockDynamicEditor({
      value,
      placeholder,
    }: {
      value: string;
      placeholder?: string;
    }) {
      return <textarea defaultValue={value} placeholder={placeholder} />;
    },
}));

import { Editor } from "@/components/Editor";

describe("Editor", () => {
  it("renders a section with the aria-label", () => {
    render(
      <Editor
        value=""
        onChange={vi.fn()}
        language="go"
        label="Go struct input editor"
        id="input-editor"
      />
    );
    expect(screen.getByRole("region", { name: "Go struct input editor" })).toBeInTheDocument();
  });

  it("renders the id on the inner div", () => {
    const { container } = render(
      <Editor
        value="hello"
        onChange={vi.fn()}
        language="go"
        label="Go struct input editor"
        id="input-editor"
      />
    );
    expect(container.querySelector("#input-editor")).toBeInTheDocument();
  });

  it("passes placeholder to the dynamic editor", () => {
    render(
      <Editor
        value=""
        language="go"
        placeholder="type here"
        label="Go struct input editor"
        id="input-editor"
      />
    );
    expect(screen.getByPlaceholderText("type here")).toBeInTheDocument();
  });

  it("renders with sql language without error", () => {
    render(
      <Editor
        value="SELECT 1"
        language="sql"
        readOnly
        label="SQL output editor"
        id="output-editor"
      />
    );
    expect(screen.getByRole("region", { name: "SQL output editor" })).toBeInTheDocument();
  });
});

import { describe, it, expect } from "vitest";
import { EditorState } from "@codemirror/state";
import { intellijTheme } from "@/lib/intellij-theme";

describe("intellijTheme", () => {
  it("is a CodeMirror extension that an editor state accepts", () => {
    expect(intellijTheme).toBeDefined();
    const state = EditorState.create({ doc: "package main", extensions: [intellijTheme] });
    expect(state.doc.toString()).toBe("package main");
  });
});

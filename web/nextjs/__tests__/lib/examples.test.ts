import { describe, it, expect } from "vitest";
import { EXAMPLES } from "@/lib/examples";

describe("EXAMPLES", () => {
  it("contains at least one example", () => {
    expect(EXAMPLES.length).toBeGreaterThan(0);
  });

  it("each example has label, description, and source", () => {
    for (const ex of EXAMPLES) {
      expect(ex.label).toBeTruthy();
      expect(ex.description).toBeTruthy();
      expect(ex.source).toContain("struct");
    }
  });

  it("each example has unique labels", () => {
    const labels = EXAMPLES.map((e) => e.label);
    expect(new Set(labels).size).toBe(labels.length);
  });
});

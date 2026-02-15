import { describe, it, expect } from "vitest";
import { validateInput, validatePackageName } from "@/lib/validation";

describe("validateInput", () => {
  it("warns on empty input", () => {
    expect(validateInput("")).toContain("Input is empty.");
    expect(validateInput("   ")).toContain("Input is empty.");
  });

  it("warns when input exceeds 100 KB", () => {
    const large = "a".repeat(100 * 1024 + 1);
    const warnings = validateInput(large);
    expect(warnings).toContain("Input exceeds 100 KB limit.");
  });

  it("warns when no exported struct found", () => {
    const warnings = validateInput("package main\n\nvar x = 1");
    expect(warnings).toContain(
      "No exported struct found. Struct names must start with an uppercase letter."
    );
  });

  it("returns no warnings for valid struct input", () => {
    const input = 'package models\n\ntype User struct {\n\tID int64 `db:"pk"`\n}';
    expect(validateInput(input)).toHaveLength(0);
  });

  it("detects unexported struct as missing", () => {
    const input = "package models\n\ntype user struct {\n\tID int64\n}";
    expect(validateInput(input)).toContain(
      "No exported struct found. Struct names must start with an uppercase letter."
    );
  });
});

describe("validatePackageName", () => {
  it("returns error for empty name", () => {
    expect(validatePackageName("")).toBe("Package name is required.");
  });

  it("returns error for uppercase name", () => {
    expect(validatePackageName("Models")).not.toBeNull();
  });

  it("returns error for name starting with digit", () => {
    expect(validatePackageName("1models")).not.toBeNull();
  });

  it("returns error for name with hyphens", () => {
    expect(validatePackageName("my-models")).not.toBeNull();
  });

  it("returns null for valid name", () => {
    expect(validatePackageName("models")).toBeNull();
    expect(validatePackageName("my_models")).toBeNull();
    expect(validatePackageName("db2")).toBeNull();
  });
});

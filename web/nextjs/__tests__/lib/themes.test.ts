import { describe, it, expect } from "vitest";
import { THEME_STORAGE_KEY, resolveInitialTheme, themes, type ThemeMode } from "@/lib/themes";

describe("resolveInitialTheme", () => {
  it("returns a stored light theme regardless of the system preference", () => {
    expect(resolveInitialTheme("light", true)).toBe("light");
    expect(resolveInitialTheme("light", false)).toBe("light");
  });

  it("returns a stored dark theme regardless of the system preference", () => {
    expect(resolveInitialTheme("dark", true)).toBe("dark");
    expect(resolveInitialTheme("dark", false)).toBe("dark");
  });

  it("follows the system preference when nothing is stored", () => {
    expect(resolveInitialTheme(null, true)).toBe("dark");
    expect(resolveInitialTheme(null, false)).toBe("light");
  });

  it("ignores an unknown stored value and follows the system preference", () => {
    expect(resolveInitialTheme("sepia", true)).toBe("dark");
    expect(resolveInitialTheme("", false)).toBe("light");
  });
});

describe("themes", () => {
  it("uses a stable storage key", () => {
    expect(THEME_STORAGE_KEY).toBe("structify_theme");
  });

  it("defines the same tokens for light and dark, none of them empty", () => {
    const modes: ThemeMode[] = ["light", "dark"];
    expect(Object.keys(themes.light).sort()).toEqual(Object.keys(themes.dark).sort());
    for (const mode of modes) {
      for (const [token, value] of Object.entries(themes[mode])) {
        expect(value, `${mode}.${token}`).not.toBe("");
      }
    }
  });

  it("gives the two modes different backgrounds and text colors", () => {
    expect(themes.light.background).not.toBe(themes.dark.background);
    expect(themes.light.textPrimary).not.toBe(themes.dark.textPrimary);
  });
});

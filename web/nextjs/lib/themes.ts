export const THEME_STORAGE_KEY = "structify_theme";

export type ThemeMode = "light" | "dark";

export type ThemeTokens = {
  background: string;
  backgroundElevated: string;
  backgroundSubtle: string;
  backgroundOverlay: string;
  textPrimary: string;
  textSecondary: string;
  textMuted: string;
  textInverse: string;
  border: string;
  borderStrong: string;
  accent: string;
  accentSoft: string;
  danger: string;
  syntaxText: string;
};

export const themes: Record<ThemeMode, ThemeTokens> = {
  dark: {
    background: "#09090b",
    backgroundElevated: "#111216",
    backgroundSubtle: "#18181c",
    backgroundOverlay: "rgba(9, 9, 11, 0.88)",
    textPrimary: "#f4f4f5",
    textSecondary: "#c9c9d1",
    textMuted: "#878793",
    textInverse: "#09090b",
    border: "#27272f",
    borderStrong: "#3b3b46",
    accent: "#74b3ff",
    accentSoft: "rgba(116, 179, 255, 0.2)",
    danger: "#ff6f7d",
    syntaxText: "#bcbec4",
  },
  light: {
    background: "#f7f7fa",
    backgroundElevated: "#ffffff",
    backgroundSubtle: "#ececf2",
    backgroundOverlay: "rgba(255, 255, 255, 0.78)",
    textPrimary: "#111218",
    textSecondary: "#2e3040",
    textMuted: "#66697a",
    textInverse: "#ffffff",
    border: "#d7d9e2",
    borderStrong: "#bec2d3",
    accent: "#235ddc",
    accentSoft: "rgba(35, 93, 220, 0.12)",
    danger: "#bb1f3f",
    syntaxText: "#2e3040",
  },
};

export function resolveInitialTheme(storedTheme: string | null, prefersDark: boolean): ThemeMode {
  if (storedTheme === "light" || storedTheme === "dark") {
    return storedTheme;
  }

  return prefersDark ? "dark" : "light";
}

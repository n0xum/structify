"use client";

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from "react";
import {
  THEME_STORAGE_KEY,
  resolveInitialTheme,
  themes,
  type ThemeMode,
} from "@/lib/themes";

type ThemeContextValue = {
  mode: ThemeMode;
  setMode: (mode: ThemeMode) => void;
  toggleMode: () => void;
};

const ThemeContext = createContext<ThemeContextValue>({
  mode: "dark",
  setMode: () => {},
  toggleMode: () => {},
});

function getInitialTheme(): ThemeMode {
  if (typeof window === "undefined") {
    return "dark";
  }

  const hasMatchMedia = typeof globalThis.matchMedia === "function";
  const prefersDark = hasMatchMedia
    ? globalThis.matchMedia("(prefers-color-scheme: dark)").matches
    : true;

  let storedTheme: string | null = null;
  try {
    storedTheme = localStorage.getItem(THEME_STORAGE_KEY);
  } catch {
    storedTheme = null;
  }

  return resolveInitialTheme(storedTheme, prefersDark);
}

function applyTheme(mode: ThemeMode) {
  const root = document.documentElement;
  const selectedTheme = themes[mode];

  root.dataset.theme = mode;
  root.style.colorScheme = mode;
  root.style.setProperty("--color-bg", selectedTheme.background);
  root.style.setProperty("--color-bg-elevated", selectedTheme.backgroundElevated);
  root.style.setProperty("--color-bg-subtle", selectedTheme.backgroundSubtle);
  root.style.setProperty("--color-bg-overlay", selectedTheme.backgroundOverlay);
  root.style.setProperty("--color-text-primary", selectedTheme.textPrimary);
  root.style.setProperty("--color-text-secondary", selectedTheme.textSecondary);
  root.style.setProperty("--color-text-muted", selectedTheme.textMuted);
  root.style.setProperty("--color-text-inverse", selectedTheme.textInverse);
  root.style.setProperty("--color-border", selectedTheme.border);
  root.style.setProperty("--color-border-strong", selectedTheme.borderStrong);
  root.style.setProperty("--color-accent", selectedTheme.accent);
  root.style.setProperty("--color-accent-soft", selectedTheme.accentSoft);
  root.style.setProperty("--color-danger", selectedTheme.danger);
  root.style.setProperty("--color-syntax-text", selectedTheme.syntaxText);
}

export function ThemeProvider({ children }: Readonly<{ children: ReactNode }>) {
  const [mode, setModeState] = useState<ThemeMode>(getInitialTheme);

  useEffect(() => {
    applyTheme(mode);
  }, [mode]);

  const setMode = useCallback((nextMode: ThemeMode) => {
    setModeState(nextMode);
    try {
      localStorage.setItem(THEME_STORAGE_KEY, nextMode);
    } catch {
      // Ignore persistence errors and keep in-memory mode.
    }
    applyTheme(nextMode);
  }, []);

  const toggleMode = useCallback(() => {
    setMode(mode === "dark" ? "light" : "dark");
  }, [mode, setMode]);

  const value = useMemo(
    () => ({
      mode,
      setMode,
      toggleMode,
    }),
    [mode, setMode, toggleMode]
  );

  return <ThemeContext.Provider value={value}>{children}</ThemeContext.Provider>;
}

export function useTheme() {
  return useContext(ThemeContext);
}

import { afterEach, beforeEach, describe, it, expect, vi } from "vitest";
import { render, screen, act } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { ThemeProvider, useTheme } from "@/components/ThemeProvider";
import { THEME_STORAGE_KEY, themes } from "@/lib/themes";

function Probe() {
  const { mode, setMode, toggleMode } = useTheme();
  return (
    <div>
      <span data-testid="mode">{mode}</span>
      <button onClick={toggleMode}>toggle</button>
      <button onClick={() => setMode("light")}>set light</button>
      <button onClick={() => setMode("dark")}>set dark</button>
    </div>
  );
}

function mockSystemPrefersDark(prefersDark: boolean) {
  vi.stubGlobal(
    "matchMedia",
    vi.fn().mockReturnValue({ matches: prefersDark }) as unknown as typeof globalThis.matchMedia
  );
}

describe("ThemeProvider", () => {
  beforeEach(() => {
    localStorage.clear();
    document.documentElement.removeAttribute("data-theme");
    document.documentElement.removeAttribute("style");
  });

  afterEach(() => {
    vi.unstubAllGlobals();
    vi.restoreAllMocks();
  });

  it("renders its children", () => {
    mockSystemPrefersDark(true);
    render(
      <ThemeProvider>
        <p>child</p>
      </ThemeProvider>
    );
    expect(screen.getByText("child")).toBeInTheDocument();
  });

  it("starts dark when the system prefers dark and nothing is stored", () => {
    mockSystemPrefersDark(true);
    render(<ThemeProvider><Probe /></ThemeProvider>);
    expect(screen.getByTestId("mode")).toHaveTextContent("dark");
  });

  it("starts light when the system prefers light and nothing is stored", () => {
    mockSystemPrefersDark(false);
    render(<ThemeProvider><Probe /></ThemeProvider>);
    expect(screen.getByTestId("mode")).toHaveTextContent("light");
  });

  it("starts dark when matchMedia does not exist", () => {
    vi.stubGlobal("matchMedia", undefined);
    render(<ThemeProvider><Probe /></ThemeProvider>);
    expect(screen.getByTestId("mode")).toHaveTextContent("dark");
  });

  it("prefers the stored theme over the system preference", () => {
    mockSystemPrefersDark(true);
    localStorage.setItem(THEME_STORAGE_KEY, "light");
    render(<ThemeProvider><Probe /></ThemeProvider>);
    expect(screen.getByTestId("mode")).toHaveTextContent("light");
  });

  it("applies the theme tokens to the document element", () => {
    mockSystemPrefersDark(true);
    render(<ThemeProvider><Probe /></ThemeProvider>);
    const root = document.documentElement;
    expect(root.dataset.theme).toBe("dark");
    expect(root.style.colorScheme).toBe("dark");
    expect(root.style.getPropertyValue("--color-bg")).toBe(themes.dark.background);
    expect(root.style.getPropertyValue("--color-accent")).toBe(themes.dark.accent);
    expect(root.style.getPropertyValue("--color-syntax-text")).toBe(themes.dark.syntaxText);
  });

  it("toggles the mode, persists it and re-applies the tokens", async () => {
    mockSystemPrefersDark(true);
    const user = userEvent.setup();
    render(<ThemeProvider><Probe /></ThemeProvider>);

    await user.click(screen.getByText("toggle"));
    expect(screen.getByTestId("mode")).toHaveTextContent("light");
    expect(localStorage.getItem(THEME_STORAGE_KEY)).toBe("light");
    expect(document.documentElement.dataset.theme).toBe("light");
    expect(document.documentElement.style.getPropertyValue("--color-bg")).toBe(themes.light.background);

    await user.click(screen.getByText("toggle"));
    expect(screen.getByTestId("mode")).toHaveTextContent("dark");
    expect(localStorage.getItem(THEME_STORAGE_KEY)).toBe("dark");
  });

  it("sets an explicit mode", async () => {
    mockSystemPrefersDark(true);
    const user = userEvent.setup();
    render(<ThemeProvider><Probe /></ThemeProvider>);

    await user.click(screen.getByText("set light"));
    expect(screen.getByTestId("mode")).toHaveTextContent("light");
    await user.click(screen.getByText("set dark"));
    expect(screen.getByTestId("mode")).toHaveTextContent("dark");
  });

  it("keeps working when reading the stored theme throws", () => {
    mockSystemPrefersDark(false);
    vi.spyOn(Storage.prototype, "getItem").mockImplementation(() => {
      throw new Error("storage blocked");
    });
    render(<ThemeProvider><Probe /></ThemeProvider>);
    expect(screen.getByTestId("mode")).toHaveTextContent("light");
  });

  it("keeps the in-memory mode when persisting the theme throws", async () => {
    mockSystemPrefersDark(true);
    vi.spyOn(Storage.prototype, "setItem").mockImplementation(() => {
      throw new Error("quota exceeded");
    });
    const user = userEvent.setup();
    render(<ThemeProvider><Probe /></ThemeProvider>);

    await act(async () => {
      await user.click(screen.getByText("toggle"));
    });
    expect(screen.getByTestId("mode")).toHaveTextContent("light");
    expect(document.documentElement.dataset.theme).toBe("light");
  });
});

describe("useTheme without a provider", () => {
  it("returns the dark default and inert setters", async () => {
    const user = userEvent.setup();
    render(<Probe />);
    expect(screen.getByTestId("mode")).toHaveTextContent("dark");
    await user.click(screen.getByText("toggle"));
    await user.click(screen.getByText("set light"));
    expect(screen.getByTestId("mode")).toHaveTextContent("dark");
  });
});

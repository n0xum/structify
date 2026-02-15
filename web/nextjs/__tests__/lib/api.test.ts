import { describe, it, expect, vi, beforeEach } from "vitest";
import { generateSQL, generateCode, fetchVersion } from "@/lib/api";

beforeEach(() => {
  vi.restoreAllMocks();
});

describe("generateSQL", () => {
  it("returns output on success", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: () => Promise.resolve({ output: "CREATE TABLE users ();" }),
      })
    );

    const result = await generateSQL("package main\ntype User struct{}");
    expect(result).toBe("CREATE TABLE users ();");
  });

  it("throws on server error", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: false,
        status: 500,
        json: () => Promise.resolve({}),
      })
    );

    await expect(generateSQL("source")).rejects.toThrow("Server error (500)");
  });

  it("throws on API error response", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: true,
        status: 422,
        json: () => Promise.resolve({ error: "no structs found" }),
      })
    );

    await expect(generateSQL("bad source")).rejects.toThrow("no structs found");
  });

  it("throws when backend unreachable", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockRejectedValue(new TypeError("fetch failed"))
    );

    await expect(generateSQL("source")).rejects.toThrow(
      "Could not reach the backend"
    );
  });
});

describe("generateCode", () => {
  it("sends source and package name", async () => {
    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      status: 200,
      json: () => Promise.resolve({ output: "package models\n" }),
    });
    vi.stubGlobal("fetch", mockFetch);

    await generateCode("source", "models");

    const body = JSON.parse(mockFetch.mock.calls[0][1].body);
    expect(body.source).toBe("source");
    expect(body.package).toBe("models");
  });
});

describe("fetchVersion", () => {
  it("returns version string", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ version: "0.1.0" }),
      })
    );

    expect(await fetchVersion()).toBe("0.1.0");
  });

  it("returns unknown on failure", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockRejectedValue(new Error("network error"))
    );

    expect(await fetchVersion()).toBe("unknown");
  });
});

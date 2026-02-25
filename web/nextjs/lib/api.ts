const BASE_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";
const STRUCTIFY_VERSION = "0.1.0";

async function post(path: string, body: object): Promise<string> {
  let res: Response;
  try {
    res = await fetch(`${BASE_URL}${path}`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    });
  } catch {
    throw new Error("Could not reach the backend. Is the server running?");
  }

  if (!res.ok && res.status >= 500) {
    throw new Error(`Server error (${res.status})`);
  }

  const data = await res.json();

  if (data && typeof data === "object" && "error" in data) {
    throw new Error(String(data.error));
  }
  if (data && typeof data === "object" && "output" in data) {
    return String(data.output);
  }
  throw new Error("Unexpected response from server.");
}

export async function generateSQL(source: string): Promise<string> {
  return post("/api/generate/sql", { source });
}

export async function generateRepository(source: string, pkg: string): Promise<string> {
  return post("/api/generate/repo", { source, package: pkg });
}

export async function fetchVersion(): Promise<string> {
  return STRUCTIFY_VERSION;
}

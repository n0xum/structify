"use client";

export default function GlobalError({
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  return (
    <html lang="en">
      <body
        style={{
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
          justifyContent: "center",
          minHeight: "100vh",
          background: "#09090b",
          color: "#f4f4f5",
          fontFamily: "sans-serif",
          textAlign: "center",
          padding: "2rem",
          gap: "1rem",
          margin: 0,
        }}
      >
        <h1 style={{ fontSize: "1.25rem", fontWeight: 600, margin: 0 }}>
          Something went wrong
        </h1>
        <p style={{ fontSize: "0.875rem", color: "#a1a1aa", margin: 0 }}>
          Please{" "}
          <a
            href="https://github.com/n0xum/structify/issues"
            style={{ color: "#38bdf8" }}
            target="_blank"
            rel="noopener noreferrer"
          >
            open an issue
          </a>{" "}
          if this keeps happening.
        </p>
        <button
          onClick={reset}
          style={{
            marginTop: "0.5rem",
            padding: "0.5rem 1rem",
            borderRadius: "0.5rem",
            background: "#0284c7",
            color: "white",
            fontSize: "0.875rem",
            border: "none",
            cursor: "pointer",
          }}
        >
          Try again
        </button>
      </body>
    </html>
  );
}

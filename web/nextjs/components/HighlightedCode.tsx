import React from "react";
import { PrismLight as SyntaxHighlighter } from "react-syntax-highlighter";
import go from "react-syntax-highlighter/dist/cjs/languages/prism/go";
import sql from "react-syntax-highlighter/dist/cjs/languages/prism/sql";

SyntaxHighlighter.registerLanguage("go", go);
SyntaxHighlighter.registerLanguage("sql", sql);

const prismTheme = {
  'code[class*="language-"]': {
    color: "var(--color-syntax-text)",
    background: "none",
    fontFamily:
      'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace',
    fontSize: "0.875rem",
    lineHeight: "1.5",
    direction: "ltr",
    textAlign: "left",
    whiteSpace: "pre",
    wordSpacing: "normal",
    wordBreak: "normal",
    tabSize: 4,
    hyphens: "none",
  },
  'pre[class*="language-"]': {
    color: "var(--color-syntax-text)",
    background: "transparent",
    margin: 0,
    overflow: "auto",
  },
  comment: { color: "var(--color-text-muted)", fontStyle: "italic" },
  prolog: { color: "var(--color-text-muted)" },
  doctype: { color: "var(--color-text-muted)" },
  cdata: { color: "var(--color-text-muted)" },
  punctuation: { color: "var(--color-syntax-text)" },
  operator: { color: "var(--color-syntax-text)" },
  keyword: { color: "#cf8e6d" },
  atrule: { color: "#cf8e6d" },
  builtin: { color: "#56a8f5" },
  function: { color: "#56a8f5" },
  "class-name": { color: "#6aab73" },
  property: { color: "#56a8f5" },
  variable: { color: "var(--color-syntax-text)" },
  constant: { color: "#2aacb8" },
  symbol: { color: "#2aacb8" },
  boolean: { color: "#cf8e6d" },
  number: { color: "#2aacb8" },
  string: { color: "#6aab73" },
  char: { color: "#6aab73" },
  regex: { color: "#6aab73" },
  namespace: { opacity: 0.8 },
  deleted: { color: "#f44747" },
  inserted: { color: "#6aab73" },
  important: { color: "#cf8e6d", fontWeight: "bold" },
  bold: { fontWeight: "bold" },
  italic: { fontStyle: "italic" },
} as Record<string, React.CSSProperties>;

type HighlightedCodeProps = {
  code: string;
  language: "go" | "sql";
};

export function HighlightedCode({ code, language }: Readonly<HighlightedCodeProps>) {
  return (
    <SyntaxHighlighter
      language={language}
      style={prismTheme}
      className="syntax-highlight"
      useInlineStyles
      customStyle={{
        margin: 0,
        padding: 0,
        background: "transparent",
      }}
    >
      {code}
    </SyntaxHighlighter>
  );
}

"use client";

import dynamic from "next/dynamic";

type EditorProps = {
  value: string;
  onChange?: (value: string) => void;
  language: "go" | "sql";
  readOnly?: boolean;
  placeholder?: string;
  label: string;
  id: string;
};

const DynamicEditor = dynamic(
  async () => {
    const [{ default: CodeMirror }, { go }, { sql }, { githubDark }] =
      await Promise.all([
        import("@uiw/react-codemirror"),
        import("@codemirror/lang-go"),
        import("@codemirror/lang-sql"),
        import("@uiw/codemirror-theme-github"),
      ]);

    return function EditorInner({
      value,
      onChange,
      language,
      readOnly,
      placeholder,
    }: Omit<EditorProps, "label" | "id">) {
      const extensions = [language === "go" ? go() : sql()];
      return (
        <CodeMirror
          value={value}
          onChange={onChange}
          extensions={extensions}
          theme={githubDark}
          readOnly={readOnly}
          placeholder={placeholder}
          basicSetup={{
            lineNumbers: true,
            foldGutter: false,
            highlightActiveLine: !readOnly,
          }}
          style={{ fontSize: "13px", minHeight: "100%" }}
          className="h-full"
        />
      );
    };
  },
  { ssr: false }
);

export function Editor({ value, onChange, language, readOnly, placeholder, label, id }: Readonly<EditorProps>) {
  return (
    <section className="flex flex-col h-full" aria-label={label}>
      <div
        id={id}
        className="flex-1 overflow-auto rounded-lg border border-zinc-700 text-sm font-mono"
      >
        <DynamicEditor
          value={value}
          onChange={onChange}
          language={language}
          readOnly={readOnly}
          placeholder={placeholder}
        />
      </div>
    </section>
  );
}

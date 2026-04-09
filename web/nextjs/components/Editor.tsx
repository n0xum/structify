"use client";

import dynamic from "next/dynamic";
import { useTheme } from "@/components/ThemeProvider";

type EditorProps = {
  value: string;
  onChange?: (value: string) => void;
  language: "go" | "sql";
  readOnly?: boolean;
  placeholder?: string;
  label: string;
  id: string;
};

type DynamicEditorProps = Omit<EditorProps, "label" | "id"> & {
  mode: "light" | "dark";
};

const DynamicEditor = dynamic(
  async () => {
    const [{ default: CodeMirror }, { go }, { sql }, { intellijTheme }, { githubLight }] =
      await Promise.all([
        import("@uiw/react-codemirror"),
        import("@codemirror/lang-go"),
        import("@codemirror/lang-sql"),
        import("@/lib/intellij-theme"),
        import("@uiw/codemirror-theme-github"),
      ]);

    return function EditorInner({
      value,
      onChange,
      language,
      readOnly,
      placeholder,
      mode,
    }: DynamicEditorProps) {
      const extensions = [language === "go" ? go() : sql()];
      return (
        <CodeMirror
          value={value}
          onChange={onChange}
          extensions={extensions}
          theme={mode === "dark" ? intellijTheme : githubLight}
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
  const { mode } = useTheme();

  return (
    <section className="flex h-full flex-col" aria-label={label}>
      <div
        id={id}
        className="flex-1 overflow-auto rounded-lg border border-[var(--color-border)] bg-[var(--color-bg-elevated)] text-sm font-mono"
      >
        <DynamicEditor
          value={value}
          onChange={onChange}
          language={language}
          readOnly={readOnly}
          placeholder={placeholder}
          mode={mode}
        />
      </div>
    </section>
  );
}

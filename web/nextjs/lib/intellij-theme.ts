import { createTheme } from "@uiw/codemirror-themes";
import { tags as t } from "@lezer/highlight";

/**
 * IntelliJ IDEA "New UI Dark" inspired theme.
 *
 * Lezer-Go tag reference (from @lezer/go/src/highlight.js):
 *   t.definitionKeyword    func interface struct chan map const type var
 *   t.moduleKeyword        import package
 *   t.controlKeyword       switch for go select return break continue goto
 *                          fallthrough case if else defer
 *   t.keyword              range
 *   t.typeName             TypeName  (int64, string, bool, error, Context, DB …)
 *   t.function(t.definition(t.variableName))   FunctionDecl/DefName  ← function being declared
 *   t.function(t.variableName)                 CallExpr/VariableName ← function being called
 *   t.definition(t.variableName)               DefName (other definitions)
 *   t.definition(t.typeName)                   TypeSpec/DefName (type name being declared)
 *   t.variableName         VariableName (parameters, local vars, package-qualified prefix)
 *   t.propertyName         FieldName (struct fields, method selector)
 *   t.modifier             ~ and * (pointer / dereference)
 *   t.string               String literals
 *   t.number               Number literals
 *   t.bool                 Bool (true / false)
 *   t.null                 Nil
 *   t.lineComment          // …
 *   t.blockComment         /* … *\/
 */
export const intellijTheme = createTheme({
  theme: "dark",
  settings: {
    background: "#1e1f22",
    backgroundImage: "",
    foreground: "#bcbec4",
    caret: "#ced0d6",
    selection: "#214283",
    selectionMatch: "#2d4a8a",
    lineHighlight: "#26282e",
    gutterBackground: "#1e1f22",
    gutterForeground: "#6f737a",
    gutterBorder: "transparent",
    gutterActiveForeground: "#bcbec4",
    fontFamily:
      'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace',
  },
  styles: [
    // ── Comments ──────────────────────────────────────────────────────────
    { tag: [t.lineComment, t.blockComment, t.docComment], color: "#7a7e85", fontStyle: "italic" },

    // ── Keywords (all variants) ────────────────────────────────────────────
    // func interface struct chan map const type var  →  definitionKeyword
    // import package                                →  moduleKeyword
    // return if for switch defer break continue …   →  controlKeyword
    // range                                         →  keyword
    {
      tag: [t.keyword, t.definitionKeyword, t.moduleKeyword, t.controlKeyword, t.operatorKeyword],
      color: "#cf8e6d",
    },

    // ── Type names ────────────────────────────────────────────────────────
    // Built-in types (int64, string, bool, error) AND user-defined types
    // (Context, DB, Order …) all share t.typeName in lezer-go.
    { tag: [t.typeName, t.typeOperator], color: "#cf8e6d" },

    // Type name being *declared* (after "type Foo …") → soft green
    { tag: t.definition(t.typeName), color: "#6aab73" },

    // ── Function names ────────────────────────────────────────────────────
    // Function being *declared*  →  FunctionDecl/DefName
    { tag: t.function(t.definition(t.variableName)), color: "#56a8f5" },
    // Function being *called*    →  CallExpr/VariableName
    { tag: t.function(t.variableName), color: "#56a8f5" },

    // ── Variables & parameters ────────────────────────────────────────────
    // Other definitions (non-function DefName), param names, local vars
    { tag: [t.variableName, t.definition(t.variableName)], color: "#bcbec4" },

    // ── Struct fields / method selector ──────────────────────────────────
    { tag: t.propertyName, color: "#bcbec4" },

    // ── Operators ─────────────────────────────────────────────────────────
    // Pointer/deref  *  and  ~  →  t.modifier (keep as text, not orange)
    { tag: t.modifier, color: "#bcbec4" },
    {
      tag: [
        t.operator,
        t.logicOperator,
        t.arithmeticOperator,
        t.bitwiseOperator,
        t.compareOperator,
        t.updateOperator,
        t.derefOperator,
        t.definitionOperator,
        t.punctuation,
        t.separator,
        t.paren,
        t.squareBracket,
        t.brace,
      ],
      color: "#bcbec4",
    },

    // ── Literals ──────────────────────────────────────────────────────────
    { tag: [t.string, t.special(t.string), t.character], color: "#6aab73" },
    { tag: t.regexp, color: "#6aab73" },
    { tag: [t.number, t.integer, t.float], color: "#2aacb8" },
    // bool & nil share the keyword orange — matches IntelliJ
    { tag: [t.bool, t.null], color: "#cf8e6d" },

    // ── Misc ──────────────────────────────────────────────────────────────
    { tag: t.labelName, color: "#bcbec4" },
    { tag: [t.meta, t.attributeName, t.attributeValue], color: "#6aab73" },
    { tag: t.invalid, color: "#f44747" },

    // ── SQL keywords (reuse keyword orange) ───────────────────────────────
    { tag: t.special(t.keyword), color: "#cf8e6d" },
    { tag: [t.heading, t.strong], color: "#bcbec4", fontWeight: "bold" },
  ],
});

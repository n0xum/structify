import React from 'react';
import { PrismLight as SyntaxHighlighter } from 'react-syntax-highlighter';
import go from 'react-syntax-highlighter/dist/cjs/languages/prism/go';
import sql from 'react-syntax-highlighter/dist/cjs/languages/prism/sql';

// Register languages
SyntaxHighlighter.registerLanguage('go', go);
SyntaxHighlighter.registerLanguage('sql', sql);

// Custom minimal monochromatic syntax theme
const customTheme = {
    'code[class*="language-"]': {
        color: '#e4e4e7', // zinc-200
        background: 'none',
        fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace',
        fontSize: '0.875rem',
        lineHeight: '1.5',
        direction: 'ltr',
        textAlign: 'left',
        whiteSpace: 'pre',
        wordSpacing: 'normal',
        wordBreak: 'normal',
        tabSize: 4,
        hyphens: 'none',
    },
    'pre[class*="language-"]': {
        color: '#e4e4e7',
        background: 'transparent',
        margin: 0,
        overflow: 'auto',
    },
    'comment': { color: '#71717a', fontStyle: 'italic' }, // zinc-500
    'prolog': { color: '#71717a' },
    'doctype': { color: '#71717a' },
    'cdata': { color: '#71717a' },
    'punctuation': { color: '#a1a1aa' }, // zinc-400
    'namespace': { opacity: 0.7 },
    'property': { color: '#e4e4e7' },
    'keyword': { color: '#ffffff', fontWeight: 'bold' }, // white
    'tag': { color: '#ffffff' },
    'class-name': { color: '#f4f4f5', textDecoration: 'underline' }, // zinc-100
    'boolean': { color: '#e4e4e7', fontWeight: 'bold' },
    'constant': { color: '#e4e4e7' },
    'symbol': { color: '#e4e4e7' },
    'deleted': { color: '#ef4444' }, // red-500 for diffs if ever used
    'selector': { color: '#e4e4e7' },
    'attr-name': { color: '#a1a1aa' },
    'string': { color: '#d4d4d8' }, // zinc-300
    'char': { color: '#d4d4d8' },
    'builtin': { color: '#ffffff' },
    'inserted': { color: '#10b981' }, // emerald-500
    'variable': { color: '#e4e4e7' },
    'operator': { color: '#a1a1aa' },
    'entity': { color: '#e4e4e7', cursor: 'help' },
    'url': { color: '#a1a1aa' },
    '.language-css .token.string': { color: '#d4d4d8' },
    '.style .token.string': { color: '#d4d4d8' },
    'atrule': { color: '#ffffff' },
    'attr-value': { color: '#d4d4d8' },
    'function': { color: '#f4f4f5', fontWeight: '500' },
    'regex': { color: '#d4d4d8' },
    'important': { color: '#ffffff', fontWeight: 'bold' },
    'bold': { fontWeight: 'bold' },
    'italic': { fontStyle: 'italic' }
} as Record<string, React.CSSProperties>;

type HighlightedCodeProps = {
    code: string;
    language: "go" | "sql";
};

export function HighlightedCode({ code, language }: Readonly<HighlightedCodeProps>) {
    return (
        <SyntaxHighlighter
            language={language}
            style={customTheme}
            customStyle={{
                margin: 0,
                padding: 0,
                background: 'transparent',
            }}
        >
            {code}
        </SyntaxHighlighter>
    );
}

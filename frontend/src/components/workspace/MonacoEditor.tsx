'use client';

import React, { useEffect, useRef } from 'react';
import Editor, { Monaco } from '@monaco-editor/react';
import type { editor } from 'monaco-editor';

interface MonacoEditorProps {
  value: string;
  onChange: (value: string | undefined) => void;
  language?: string;
  readOnly?: boolean;
  fontSize?: number;
  minimap?: boolean;
  height?: string;
  theme?: 'vs-dark' | 'light';
  onError?: (lineNumber: number) => void;
}

export default function MonacoEditor({
  value,
  onChange,
  language = 'go',
  readOnly = false,
  fontSize = 14,
  minimap = false,
  height = '100%',
  theme = 'vs-dark',
  onError,
}: MonacoEditorProps) {
  const editorRef = useRef<editor.IStandaloneCodeEditor | null>(null);

  const handleEditorDidMount = (editor: editor.IStandaloneCodeEditor, monaco: Monaco) => {
    editorRef.current = editor;

    // Configure Go language support
    monaco.languages.setMonarchTokensProvider('go', {
      keywords: [
        'break', 'case', 'chan', 'const', 'continue', 'default', 'defer', 'else',
        'fallthrough', 'for', 'func', 'go', 'goto', 'if', 'import', 'interface',
        'map', 'package', 'range', 'return', 'select', 'struct', 'switch', 'type',
        'var',
      ],
      operators: [
        '+', '-', '*', '/', '%', '&', '|', '^', '<<', '>>', '&^', '+=', '-=',
        '*=', '/=', '%=', '&=', '|=', '^=', '<<=', '>>=', '&^=', '&&', '||',
        '<-', '++', '--', '==', '<', '>', '=', '!', '!=', '<=', '>=', ':=', '...',
        '(', ')', '[', ']', '{', '}', ',', ';', '.', ':',
      ],
      symbols: /[=><!~?:&|+\-*\/\^%]+/,
      escapes: /\\(?:[abfnrtv\\"']|x[0-9A-Fa-f]{1,4}|u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8})/,
      tokenizer: {
        root: [
          [/[a-z_$][\w$]*/, {
            cases: {
              '@keywords': 'keyword',
              '@default': 'identifier',
            },
          }],
          [/[A-Z][\w\$]*/, 'type.identifier'],
          { include: '@whitespace' },
          [/[:{}]/, '@brackets'],
          [/@symbols/, {
            cases: {
              '@operators': 'operator',
              '@default': '',
            },
          }],
          [/\d*\.\d+([eE][\-+]?\d+)?/, 'number.float'],
          [/\d+/, 'number'],
          [/[;,.]/, 'delimiter'],
          [/"([^"\\]|\\.)*$/, 'string.invalid'],
          [/"/, 'string', '@string'],
          [/`/, 'string', '@rawstring'],
        ],
        whitespace: [
          [/[ \t\r\n]+/, 'white'],
          [/\/\/.*$/, 'comment'],
        ],
        string: [
          [/[^\\"]+/, 'string'],
          [/@escapes/, 'string.escape'],
          [/\\./, 'string.escape.invalid'],
          [/"/, 'string', '@pop'],
        ],
        rawstring: [
          [/[^\`]+/, 'string'],
          [/`/, 'string', '@pop'],
        ],
      },
    });

    // Configure IntelliSense for Go standard library
    monaco.languages.registerCompletionItemProvider('go', {
      provideCompletionItems: (model, position) => {
        const word = model.getWordUntilPosition(position);
        const range = {
          startLineNumber: position.lineNumber,
          endLineNumber: position.lineNumber,
          startColumn: word.startColumn,
          endColumn: word.endColumn,
        };

        // Go standard library suggestions
        const suggestions = [
          // Common packages
          { label: 'fmt', kind: monaco.languages.CompletionItemKind.Module, insertText: 'fmt', detail: 'Format package', range },
          { label: 'os', kind: monaco.languages.CompletionItemKind.Module, insertText: 'os', detail: 'OS package', range },
          { label: 'io', kind: monaco.languages.CompletionItemKind.Module, insertText: 'io', detail: 'IO package', range },
          { label: 'strings', kind: monaco.languages.CompletionItemKind.Module, insertText: 'strings', detail: 'Strings package', range },
          { label: 'strconv', kind: monaco.languages.CompletionItemKind.Module, insertText: 'strconv', detail: 'String conversion package', range },
          { label: 'net/http', kind: monaco.languages.CompletionItemKind.Module, insertText: 'net/http', detail: 'HTTP package', range },
          { label: 'encoding/json', kind: monaco.languages.CompletionItemKind.Module, insertText: 'encoding/json', detail: 'JSON package', range },

          // Common functions
          { label: 'main', kind: monaco.languages.CompletionItemKind.Function, insertText: 'main()', detail: 'Main function', range },
          { label: 'init', kind: monaco.languages.CompletionItemKind.Function, insertText: 'init()', detail: 'Init function', range },
          { label: 'func', kind: monaco.languages.CompletionItemKind.Keyword, insertText: 'func ', detail: 'Function keyword', range },
          { label: 'var', kind: monaco.languages.CompletionItemKind.Keyword, insertText: 'var ', detail: 'Variable declaration', range },
          { label: 'const', kind: monaco.languages.CompletionItemKind.Keyword, insertText: 'const ', detail: 'Constant declaration', range },
          { label: 'type', kind: monaco.languages.CompletionItemKind.Keyword, insertText: 'type ', detail: 'Type declaration', range },
          { label: 'struct', kind: monaco.languages.CompletionItemKind.Keyword, insertText: 'struct', detail: 'Struct keyword', range },
          { label: 'interface', kind: monaco.languages.CompletionItemKind.Keyword, insertText: 'interface', detail: 'Interface keyword', range },
          { label: 'map', kind: monaco.languages.CompletionItemKind.Keyword, insertText: 'map[]', detail: 'Map keyword', range },

          // Common fmt functions
          { label: 'fmt.Println', kind: monaco.languages.CompletionItemKind.Function, insertText: 'fmt.Println(', detail: 'Print to stdout', range },
          { label: 'fmt.Printf', kind: monaco.languages.CompletionItemKind.Function, insertText: 'fmt.Printf(', detail: 'Formatted print', range },
          { label: 'fmt.Sprintf', kind: monaco.languages.CompletionItemKind.Function, insertText: 'fmt.Sprintf(', detail: 'Formatted string', range },
          { label: 'fmt.Scanln', kind: monaco.languages.CompletionItemKind.Function, insertText: 'fmt.Scanln(', detail: 'Scan from stdin', range },
        ];

        return { suggestions };
      },
    });

    // Enable error navigation
    if (onError) {
      editor.onMouseDown((e) => {
        const target = e.target;
        if (target && 'position' in target) {
          const position = (target as { position: { lineNumber: number } | null }).position;
          if (position) {
            const model = editor.getModel();
            if (model) {
              const line = model.getLineContent(position.lineNumber);
              // Check if line contains error indicator
              if (line.includes('error:') || line.includes('Error')) {
                onError(position.lineNumber);
              }
            }
          }
        }
      });
    }

    // Configure editor options
    editor.updateOptions({
      fontSize,
      minimap: { enabled: minimap },
      readOnly,
      scrollBeyondLastLine: false,
      automaticLayout: true,
      tabSize: 4,
      insertSpaces: false,
      wordWrap: 'off',
      lineNumbers: 'on',
      renderWhitespace: 'selection',
      renderLineHighlight: 'all',
      cursorBlinking: 'smooth',
      cursorSmoothCaretAnimation: 'on',
      smoothScrolling: true,
      folding: true,
      foldingStrategy: 'indentation',
      showFoldingControls: 'always',
      formatOnPaste: true,
      formatOnType: true,
      autoIndent: 'full',
    });
  };

  return (
    <div className="h-full w-full">
      <Editor
        height={height}
        language={language}
        value={value}
        onChange={onChange}
        theme={theme}
        onMount={handleEditorDidMount}
        options={{
          fontSize,
          minimap: { enabled: minimap },
          readOnly,
          scrollBeyondLastLine: false,
          automaticLayout: true,
          tabSize: 4,
          insertSpaces: false,
          wordWrap: 'off',
          lineNumbers: 'on',
          renderWhitespace: 'selection',
          renderLineHighlight: 'all',
        }}
        loading={<div className="flex items-center justify-center h-full">Loading editor...</div>}
      />
    </div>
  );
}

import Editor, { OnMount } from '@monaco-editor/react';
import * as monaco from 'monaco-editor';

interface PlaygroundEditorProps {
  value: string;
  onChange: (value: string | undefined) => void;
  language: string;
  height: string;
}

export default function PlaygroundEditor({
  value,
  onChange,
  language,
  height,
}: PlaygroundEditorProps) {
  const handleEditorMount: OnMount = (editor: monaco.editor.IStandaloneCodeEditor) => {
    editor.updateOptions({
      tabSize: 4,
      insertSpaces: true,
    });
  };

  return (
    <div style={{ height, width: '100%' }}>
      <Editor
        height={height}
        language={language}
        value={value}
        onChange={onChange}
        onMount={handleEditorMount}
        theme="vs-dark"
        options={{
          minimap: { enabled: false },
          fontSize: 14,
          lineNumbers: 'on',
          scrollBeyondLastLine: false,
          automaticLayout: true,
          padding: { top: 16, bottom: 16 },
        }}
      />
    </div>
  );
}

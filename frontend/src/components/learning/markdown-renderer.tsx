"use client";

import { useState, useEffect } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { 
  Copy, 
  Check, 
  Play, 
  Eye, 
  EyeOff,
  Code2,
  BookOpen,
  Lightbulb,
  AlertTriangle,
  Info,
  CheckCircle,
  Zap
} from "lucide-react";

interface MarkdownRendererProps {
  content: string;
  className?: string;
  enableCodeHighlight?: boolean;
  enableInteractiveExamples?: boolean;
}

interface CodeBlock {
  language: string;
  code: string;
  title?: string;
  runnable?: boolean;
}

export default function MarkdownRenderer({ 
  content, 
  className = "",
  enableCodeHighlight = true,
  enableInteractiveExamples = true
}: MarkdownRendererProps) {
  const [copiedCode, setCopiedCode] = useState<string | null>(null);
  const [expandedSections, setExpandedSections] = useState<Set<string>>(new Set());

  // Parse markdown content into structured elements
  const parseMarkdown = (text: string) => {
    const elements: any[] = [];
    const lines = text.split('\n');
    let currentElement: any = null;
    let inCodeBlock = false;
    let codeBlockContent = '';
    let codeBlockLanguage = '';

    for (let i = 0; i < lines.length; i++) {
      const line = lines[i];

      // Code blocks
      if (line.startsWith('```')) {
        if (!inCodeBlock) {
          inCodeBlock = true;
          codeBlockLanguage = line.slice(3).trim();
          codeBlockContent = '';
        } else {
          inCodeBlock = false;
          elements.push({
            type: 'code',
            language: codeBlockLanguage,
            content: codeBlockContent.trim(),
            runnable: codeBlockLanguage === 'go'
          });
          codeBlockContent = '';
        }
        continue;
      }

      if (inCodeBlock) {
        codeBlockContent += line + '\n';
        continue;
      }

      // Headers
      if (line.startsWith('#')) {
        const level = line.match(/^#+/)?.[0].length || 1;
        const text = line.replace(/^#+\s*/, '');
        elements.push({
          type: 'header',
          level,
          content: text,
          id: text.toLowerCase().replace(/\s+/g, '-').replace(/[^\w-]/g, '')
        });
        continue;
      }

      // Lists
      if (line.match(/^[\s]*[-*+]\s/)) {
        if (currentElement?.type !== 'list') {
          currentElement = { type: 'list', items: [] };
          elements.push(currentElement);
        }
        currentElement.items.push(line.replace(/^[\s]*[-*+]\s/, ''));
        continue;
      }

      // Callouts/Alerts
      if (line.startsWith('> ')) {
        const alertType = line.match(/^>\s*\*\*(Note|Warning|Tip|Important)\*\*:/i)?.[1]?.toLowerCase();
        if (alertType) {
          elements.push({
            type: 'alert',
            alertType,
            content: line.replace(/^>\s*\*\*[^*]+\*\*:\s*/, '')
          });
          continue;
        }
      }

      // Paragraphs
      if (line.trim()) {
        if (currentElement?.type !== 'paragraph') {
          currentElement = { type: 'paragraph', content: '' };
          elements.push(currentElement);
        }
        currentElement.content += (currentElement.content ? ' ' : '') + line.trim();
      } else {
        currentElement = null;
      }
    }

    return elements;
  };

  const copyToClipboard = async (text: string, id: string) => {
    try {
      await navigator.clipboard.writeText(text);
      setCopiedCode(id);
      setTimeout(() => setCopiedCode(null), 2000);
    } catch (err) {
      console.error('Failed to copy:', err);
    }
  };

  const toggleSection = (id: string) => {
    const newExpanded = new Set(expandedSections);
    if (newExpanded.has(id)) {
      newExpanded.delete(id);
    } else {
      newExpanded.add(id);
    }
    setExpandedSections(newExpanded);
  };

  const renderCodeBlock = (element: any, index: number) => {
    const codeId = `code-${index}`;
    const isCopied = copiedCode === codeId;

    return (
      <Card key={index} className="my-4 overflow-hidden border-2 hover:border-primary/30 transition-colors">
        <div className="flex items-center justify-between p-3 bg-muted/50 border-b">
          <div className="flex items-center space-x-2">
            <Code2 className="h-4 w-4 text-primary" />
            <Badge variant="outline" className="text-xs">
              {element.language}
            </Badge>
            {element.runnable && (
              <Badge variant="secondary" className="text-xs">
                <Play className="mr-1 h-3 w-3" />
                Runnable
              </Badge>
            )}
          </div>
          <div className="flex items-center space-x-2">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => copyToClipboard(element.content, codeId)}
              className="h-8 w-8 p-0"
            >
              {isCopied ? (
                <Check className="h-4 w-4 text-green-500" />
              ) : (
                <Copy className="h-4 w-4" />
              )}
            </Button>
            {element.runnable && enableInteractiveExamples && (
              <Button
                variant="ghost"
                size="sm"
                className="h-8 w-8 p-0"
                onClick={() => console.log('Run code:', element.content)}
              >
                <Play className="h-4 w-4" />
              </Button>
            )}
          </div>
        </div>
        <CardContent className="p-0">
          <pre className="p-4 overflow-x-auto text-sm bg-background">
            <code className={`language-${element.language}`}>
              {element.content}
            </code>
          </pre>
        </CardContent>
      </Card>
    );
  };

  const renderAlert = (element: any, index: number) => {
    const alertConfig = {
      note: { icon: Info, className: "border-blue-200 bg-blue-50 dark:bg-blue-950 dark:border-blue-800", iconColor: "text-blue-500" },
      warning: { icon: AlertTriangle, className: "border-yellow-200 bg-yellow-50 dark:bg-yellow-950 dark:border-yellow-800", iconColor: "text-yellow-500" },
      tip: { icon: Lightbulb, className: "border-green-200 bg-green-50 dark:bg-green-950 dark:border-green-800", iconColor: "text-green-500" },
      important: { icon: Zap, className: "border-red-200 bg-red-50 dark:bg-red-950 dark:border-red-800", iconColor: "text-red-500" }
    };

    const config = alertConfig[element.alertType as keyof typeof alertConfig] || alertConfig.note;
    const Icon = config.icon;

    return (
      <Card key={index} className={`my-4 border-2 ${config.className}`}>
        <CardContent className="p-4">
          <div className="flex items-start space-x-3">
            <Icon className={`h-5 w-5 mt-0.5 ${config.iconColor}`} />
            <div className="flex-1">
              <p className="text-sm leading-relaxed">{element.content}</p>
            </div>
          </div>
        </CardContent>
      </Card>
    );
  };

  const renderElement = (element: any, index: number) => {
    switch (element.type) {
      case 'header':
        const HeaderTag = `h${Math.min(element.level, 6)}` as keyof JSX.IntrinsicElements;
        const headerClasses = {
          1: "text-3xl font-bold mb-4 mt-8",
          2: "text-2xl font-semibold mb-3 mt-6",
          3: "text-xl font-medium mb-2 mt-4",
          4: "text-lg font-medium mb-2 mt-3",
          5: "text-base font-medium mb-1 mt-2",
          6: "text-sm font-medium mb-1 mt-2"
        };
        return (
          <HeaderTag 
            key={index} 
            id={element.id}
            className={`${headerClasses[element.level as keyof typeof headerClasses]} scroll-mt-20 stagger-item`}
          >
            {element.content}
          </HeaderTag>
        );

      case 'paragraph':
        return (
          <p key={index} className="mb-4 leading-relaxed text-muted-foreground stagger-item">
            {element.content}
          </p>
        );

      case 'list':
        return (
          <ul key={index} className="mb-4 space-y-2 stagger-item">
            {element.items.map((item: string, itemIndex: number) => (
              <li key={itemIndex} className="flex items-start space-x-2">
                <div className="w-1.5 h-1.5 rounded-full bg-primary mt-2 flex-shrink-0" />
                <span className="text-muted-foreground">{item}</span>
              </li>
            ))}
          </ul>
        );

      case 'code':
        return renderCodeBlock(element, index);

      case 'alert':
        return renderAlert(element, index);

      default:
        return null;
    }
  };

  const elements = parseMarkdown(content);

  return (
    <div className={`prose dark:prose-invert max-w-none ${className}`}>
      {elements.map((element, index) => renderElement(element, index))}
    </div>
  );
}

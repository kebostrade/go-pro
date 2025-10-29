"use client";

import { useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { 
  GitCompare, 
  Eye, 
  Code2, 
  ArrowLeftRight,
  Check,
  X,
  Info
} from "lucide-react";

interface CodeDiffViewerProps {
  originalCode: string;
  modifiedCode: string;
  originalLabel?: string;
  modifiedLabel?: string;
  language?: string;
  showLineNumbers?: boolean;
}

interface DiffLine {
  type: 'added' | 'removed' | 'unchanged' | 'modified';
  originalLineNumber?: number;
  modifiedLineNumber?: number;
  content: string;
  originalContent?: string;
}

/**
 * Calculate diff between two code strings
 */
function calculateDiff(original: string, modified: string): DiffLine[] {
  const originalLines = original.split('\n');
  const modifiedLines = modified.split('\n');
  const diff: DiffLine[] = [];
  
  let originalIndex = 0;
  let modifiedIndex = 0;
  
  while (originalIndex < originalLines.length || modifiedIndex < modifiedLines.length) {
    const originalLine = originalLines[originalIndex];
    const modifiedLine = modifiedLines[modifiedIndex];
    
    if (originalLine === modifiedLine) {
      // Lines are identical
      diff.push({
        type: 'unchanged',
        originalLineNumber: originalIndex + 1,
        modifiedLineNumber: modifiedIndex + 1,
        content: originalLine || '',
      });
      originalIndex++;
      modifiedIndex++;
    } else if (originalIndex >= originalLines.length) {
      // Added line
      diff.push({
        type: 'added',
        modifiedLineNumber: modifiedIndex + 1,
        content: modifiedLine || '',
      });
      modifiedIndex++;
    } else if (modifiedIndex >= modifiedLines.length) {
      // Removed line
      diff.push({
        type: 'removed',
        originalLineNumber: originalIndex + 1,
        content: originalLine || '',
      });
      originalIndex++;
    } else {
      // Lines are different - check if it's a modification or add/remove
      const nextOriginalMatch = modifiedLines.slice(modifiedIndex).indexOf(originalLine);
      const nextModifiedMatch = originalLines.slice(originalIndex).indexOf(modifiedLine);
      
      if (nextOriginalMatch === 0) {
        // Current modified line was added
        diff.push({
          type: 'added',
          modifiedLineNumber: modifiedIndex + 1,
          content: modifiedLine,
        });
        modifiedIndex++;
      } else if (nextModifiedMatch === 0) {
        // Current original line was removed
        diff.push({
          type: 'removed',
          originalLineNumber: originalIndex + 1,
          content: originalLine,
        });
        originalIndex++;
      } else {
        // Lines are modified
        diff.push({
          type: 'modified',
          originalLineNumber: originalIndex + 1,
          modifiedLineNumber: modifiedIndex + 1,
          content: modifiedLine,
          originalContent: originalLine,
        });
        originalIndex++;
        modifiedIndex++;
      }
    }
  }
  
  return diff;
}

/**
 * Calculate diff statistics
 */
function calculateStats(diff: DiffLine[]) {
  const added = diff.filter(line => line.type === 'added').length;
  const removed = diff.filter(line => line.type === 'removed').length;
  const modified = diff.filter(line => line.type === 'modified').length;
  const unchanged = diff.filter(line => line.type === 'unchanged').length;
  
  return { added, removed, modified, unchanged, total: diff.length };
}

export default function CodeDiffViewer({
  originalCode,
  modifiedCode,
  originalLabel = "Original",
  modifiedLabel = "Modified",
  language = "go",
  showLineNumbers = true,
}: CodeDiffViewerProps) {
  const [viewMode, setViewMode] = useState<'unified' | 'split'>('unified');
  const [showOnlyChanges, setShowOnlyChanges] = useState(false);
  
  const diff = calculateDiff(originalCode, modifiedCode);
  const stats = calculateStats(diff);
  
  const filteredDiff = showOnlyChanges 
    ? diff.filter(line => line.type !== 'unchanged')
    : diff;
  
  const getLineClass = (type: DiffLine['type']) => {
    switch (type) {
      case 'added':
        return 'bg-green-50 dark:bg-green-950 border-l-4 border-green-500';
      case 'removed':
        return 'bg-red-50 dark:bg-red-950 border-l-4 border-red-500';
      case 'modified':
        return 'bg-yellow-50 dark:bg-yellow-950 border-l-4 border-yellow-500';
      default:
        return 'bg-background';
    }
  };
  
  const getLineIcon = (type: DiffLine['type']) => {
    switch (type) {
      case 'added':
        return <span className="text-green-600 dark:text-green-400 font-bold">+</span>;
      case 'removed':
        return <span className="text-red-600 dark:text-red-400 font-bold">-</span>;
      case 'modified':
        return <span className="text-yellow-600 dark:text-yellow-400 font-bold">~</span>;
      default:
        return <span className="text-muted-foreground"> </span>;
    }
  };

  return (
    <Card className="glass-card border-2">
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="flex items-center">
              <GitCompare className="mr-2 h-5 w-5 text-primary" />
              Code Comparison
            </CardTitle>
            <CardDescription>
              Compare your code with the solution
            </CardDescription>
          </div>
          
          <div className="flex items-center space-x-2">
            <Button
              variant={viewMode === 'unified' ? 'default' : 'outline'}
              size="sm"
              onClick={() => setViewMode('unified')}
            >
              <Eye className="mr-2 h-4 w-4" />
              Unified
            </Button>
            <Button
              variant={viewMode === 'split' ? 'default' : 'outline'}
              size="sm"
              onClick={() => setViewMode('split')}
            >
              <ArrowLeftRight className="mr-2 h-4 w-4" />
              Split
            </Button>
          </div>
        </div>
        
        {/* Statistics */}
        <div className="flex flex-wrap gap-2 mt-4">
          <Badge variant="outline" className="border-green-500 text-green-600">
            <Check className="mr-1 h-3 w-3" />
            {stats.added} added
          </Badge>
          <Badge variant="outline" className="border-red-500 text-red-600">
            <X className="mr-1 h-3 w-3" />
            {stats.removed} removed
          </Badge>
          <Badge variant="outline" className="border-yellow-500 text-yellow-600">
            <Info className="mr-1 h-3 w-3" />
            {stats.modified} modified
          </Badge>
          <Badge variant="outline">
            {stats.unchanged} unchanged
          </Badge>
        </div>
      </CardHeader>
      
      <CardContent>
        <div className="space-y-4">
          {/* Controls */}
          <div className="flex items-center justify-between">
            <Button
              variant="outline"
              size="sm"
              onClick={() => setShowOnlyChanges(!showOnlyChanges)}
            >
              {showOnlyChanges ? 'Show All Lines' : 'Show Only Changes'}
            </Button>
            
            <div className="text-sm text-muted-foreground">
              {filteredDiff.length} lines displayed
            </div>
          </div>
          
          {/* Diff Display */}
          {viewMode === 'unified' ? (
            <div className="border rounded-lg overflow-hidden">
              <div className="bg-muted px-4 py-2 text-sm font-medium border-b">
                {originalLabel} → {modifiedLabel}
              </div>
              <div className="max-h-[500px] overflow-y-auto custom-scrollbar">
                {filteredDiff.map((line, index) => (
                  <div
                    key={index}
                    className={`flex items-start font-mono text-sm ${getLineClass(line.type)} transition-colors`}
                  >
                    {showLineNumbers && (
                      <div className="flex-shrink-0 w-20 px-2 py-1 text-xs text-muted-foreground text-right border-r">
                        <span className="inline-block w-8">
                          {line.originalLineNumber || ''}
                        </span>
                        <span className="inline-block w-8">
                          {line.modifiedLineNumber || ''}
                        </span>
                      </div>
                    )}
                    <div className="flex-shrink-0 w-8 px-2 py-1 text-center">
                      {getLineIcon(line.type)}
                    </div>
                    <div className="flex-1 px-2 py-1 overflow-x-auto">
                      {line.type === 'modified' && line.originalContent ? (
                        <div className="space-y-1">
                          <div className="text-red-600 dark:text-red-400 line-through opacity-70">
                            {line.originalContent}
                          </div>
                          <div className="text-green-600 dark:text-green-400">
                            {line.content}
                          </div>
                        </div>
                      ) : (
                        <div className={
                          line.type === 'added' ? 'text-green-600 dark:text-green-400' :
                          line.type === 'removed' ? 'text-red-600 dark:text-red-400' :
                          ''
                        }>
                          {line.content || ' '}
                        </div>
                      )}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          ) : (
            <div className="grid grid-cols-2 gap-4">
              {/* Original Code */}
              <div className="border rounded-lg overflow-hidden">
                <div className="bg-muted px-4 py-2 text-sm font-medium border-b">
                  {originalLabel}
                </div>
                <div className="max-h-[500px] overflow-y-auto custom-scrollbar">
                  {originalCode.split('\n').map((line, index) => (
                    <div
                      key={index}
                      className="flex items-start font-mono text-sm hover:bg-muted/50"
                    >
                      {showLineNumbers && (
                        <div className="flex-shrink-0 w-12 px-2 py-1 text-xs text-muted-foreground text-right border-r">
                          {index + 1}
                        </div>
                      )}
                      <div className="flex-1 px-2 py-1 overflow-x-auto">
                        {line || ' '}
                      </div>
                    </div>
                  ))}
                </div>
              </div>
              
              {/* Modified Code */}
              <div className="border rounded-lg overflow-hidden">
                <div className="bg-muted px-4 py-2 text-sm font-medium border-b">
                  {modifiedLabel}
                </div>
                <div className="max-h-[500px] overflow-y-auto custom-scrollbar">
                  {modifiedCode.split('\n').map((line, index) => (
                    <div
                      key={index}
                      className="flex items-start font-mono text-sm hover:bg-muted/50"
                    >
                      {showLineNumbers && (
                        <div className="flex-shrink-0 w-12 px-2 py-1 text-xs text-muted-foreground text-right border-r">
                          {index + 1}
                        </div>
                      )}
                      <div className="flex-1 px-2 py-1 overflow-x-auto">
                        {line || ' '}
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}


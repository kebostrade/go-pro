"use client";

import { useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Code2,
  Lightbulb,
  BookOpen,
  ChevronRight,
  Plus,
  Play
} from "lucide-react";

interface Pattern {
  id: string;
  name: string;
  description: string;
  template: string;
  problems: string[];
  whenToUse: string[];
  confidence: number;
}

interface PatternNotesProps {}

const PatternNotes = ({}: PatternNotesProps) => {
  const [expandedPattern, setExpandedPattern] = useState<string | null>(null);

  const patterns: Pattern[] = [
    {
      id: 'sliding-window',
      name: 'Sliding Window',
      description: 'For contiguous subarrays/substrings with constraints',
      template: `left := 0
for right := 0; right < len(arr); right++ {
  // expand window
  for !valid(window) {
    // shrink from left
    left++
  }
  // update result
}`,
      problems: ['Longest Substring Without Repeating', 'Minimum Window Substring'],
      whenToUse: ['Contiguous elements', 'Max/min length', 'Fixed window size'],
      confidence: 3,
    },
    {
      id: 'two-pointers',
      name: 'Two Pointers',
      description: 'For sorted arrays, pairs with sum, palindromes',
      template: `left, right := 0, len(arr)-1
for left < right {
  sum := arr[left] + arr[right]
  if sum == target { return }
  if sum < target { left++ } else { right-- }
}`,
      problems: ['Two Sum II', '3Sum', 'Container With Most Water'],
      whenToUse: ['Sorted array', 'Finding pairs', 'Palindrome check'],
      confidence: 4,
    },
    {
      id: 'fast-slow',
      name: 'Fast & Slow Pointers',
      description: 'For cycle detection, finding middle',
      template: `slow, fast := head, head
for fast != nil && fast.Next != nil {
  slow = slow.Next
  fast = fast.Next.Next
  if slow == fast { return true } // cycle
}`,
      problems: ['Linked List Cycle', 'Middle of Linked List'],
      whenToUse: ['Cycle detection', 'Find middle element', 'Linked list problems'],
      confidence: 2,
    },
    {
      id: 'bfs',
      name: 'BFS (Level Order)',
      description: 'For shortest path, level-by-level processing',
      template: `queue := [root]
while len(queue) > 0 {
  levelSize := len(queue)
  for i := 0; i < levelSize; i++ {
    node := queue[0]
    queue = queue[1:]
    // process node
  }
}`,
      problems: ['Binary Tree Level Order', 'Number of Islands'],
      whenToUse: ['Shortest path', 'Level traversal', 'Tree by levels'],
      confidence: 4,
    },
    {
      id: 'dfs',
      name: 'DFS (Depth First)',
      description: 'For exhaustive search, backtracking',
      template: `func dfs(node *Node) {
  if baseCase { return }
  for _, neighbor := range neighbors {
    dfs(neighbor)
  }`,
      problems: ['Max Path Sum', 'Path Sum III'],
      whenToUse: ['Exhaustive search', 'All paths', 'Backtracking base'],
      confidence: 3,
    },
    {
      id: 'binary-search',
      name: 'Binary Search',
      description: 'For sorted data, find specific value',
      template: `left, right := 0, len(arr)-1
for left <= right {
  mid := left + (right-left)/2
  if arr[mid] == target { return mid }
  if arr[mid] < target { left = mid + 1 } else { right = mid - 1 }
}`,
      problems: ['Search in Rotated Array', 'Find Minimum in Rotated'],
      whenToUse: ['Sorted array', 'Find target', 'Rotated sorted array'],
      confidence: 5,
    },
  ];

  const getConfidenceColor = (confidence: number) => {
    if (confidence >= 5) return 'text-green-600 bg-green-100';
    if (confidence >= 4) return 'text-blue-600 bg-blue-100';
    if (confidence >= 3) return 'text-yellow-600 bg-yellow-100';
    return 'text-gray-600 bg-gray-100';
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-semibold flex items-center gap-2">
          <Lightbulb className="h-5 w-5 text-yellow-500" />
          Pattern Reference
        </h3>
        <Button size="sm" variant="outline">
          <Plus className="h-4 w-4 mr-1" />
          Add Pattern
        </Button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {patterns.map((pattern) => {
          const isExpanded = expandedPattern === pattern.id;
          const problemCount = pattern.problems.length;

          return (
            <Card
              key={pattern.id}
              className={`cursor-pointer transition-all hover:shadow-sm ${
                isExpanded ? 'ring-2 ring-primary' : ''
              }`}
              onClick={() => setExpandedPattern(isExpanded ? null : pattern.id)}
            >
              <CardHeader className="pb-2">
                <div className="flex items-center justify-between">
                  <CardTitle className="text-base flex items-center gap-2">
                    <Code2 className="h-4 w-4" />
                    {pattern.name}
                  </CardTitle>
                  <Badge className={getConfidenceColor(pattern.confidence)}>
                    {pattern.confidence}/5
                  </Badge>
                </div>
                <CardDescription className="text-xs mt-1">
                  {pattern.description}
                </CardDescription>
              </CardHeader>
              <CardContent className="pt-2">
                <div className="space-y-3">
                  <div className="flex items-center gap-2">
                    <BookOpen className="h-4 w-4 text-muted-foreground" />
                    <span className="text-xs text-muted-foreground">
                      {problemCount} problem{problemCount !== 1 ? 's' : ''} practiced
                    </span>
                  </div>

                  <div className="text-xs text-muted-foreground">
                    <strong>When to use:</strong>
                    <ul className="mt-1 space-y-1 list-disc pl-4">
                      {pattern.whenToUse.map((use, idx) => (
                        <li key={idx} className="flex items-center gap-1">
                          <ChevronRight className="h-3 w-3" />
                          {use}
                        </li>
                      ))}
                    </ul>
                  </div>

                  {isExpanded && (
                    <div className="pt-2 border-t mt-2">
                      <p className="text-xs font-medium mb-2">Template:</p>
                      <pre className="text-xs bg-muted p-2 rounded overflow-x-auto font-mono">
                        {pattern.template}
                      </pre>
                      {problemCount > 0 && (
                        <div className="mt-3">
                          <p className="text-xs font-medium mb-2">Related Problems:</p>
                          <div className="flex flex-wrap gap-1">
                            {pattern.problems.map((problem) => (
                            <Badge key={problem} variant="outline" className="text-xs">
                              {problem}
                            </Badge>
                            ))}
                          </div>
                        </div>
                      )}
                    </div>
                  )}
                </div>
              </CardContent>
            </Card>
          );
        })}
      </div>

      {/* Quick start section */}
      <Card className="border-dashed">
        <CardHeader>
          <CardTitle className="text-sm">Learn New Pattern</CardTitle>
          <CardDescription>
            Master common algorithm patterns to solve problems faster
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex gap-2">
            <Button variant="outline" size="sm" className="flex-1">
              <BookOpen className="h-4 w-4 mr-1" />
              Study Guide
            </Button>
            <Button size="sm" className="flex-1">
              <Play className="h-4 w-4 mr-1" />
              Practice
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
};

export default PatternNotes;

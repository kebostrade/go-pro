"use client";

import { useState, useEffect, useRef } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Textarea } from "@/components/ui/textarea";
import { Brain, Code2, Timer, Copy, CheckCircle, Lightbulb, Zap, Clock, Target, Database, Play, Pause, RotateCcw, AlertCircle, BookOpen, Calculator, BarChart3, CheckCircle2 } from "lucide-react";

interface ComplexityResult { time: string; space: string; explanation: string; }
interface Snippet { id: string; name: string; pattern: string; code: string; whenToUse: string; }

const codeSnippets: Snippet[] = [
  { id: "two-pointer", name: "Two Pointers", pattern: "O(n) time, O(1) space", code: "function twoPointers(arr, target) {\n  let left = 0, right = arr.length - 1;\n  while (left < right) {\n    const sum = arr[left] + arr[right];\n    if (sum === target) return [left, right];\n    if (sum < target) left++;\n    else right--;\n  }\n  return [-1, -1];\n}", whenToUse: "Sorted arrays, sum problems" },
  { id: "sliding-window", name: "Sliding Window", pattern: "O(n) time, O(k) space", code: "function maxSum(arr, k) {\n  let windowSum = 0, maxSum = 0;\n  for (let i = 0; i < arr.length; i++) {\n    windowSum += arr[i];\n    if (i >= k - 1) {\n      maxSum = Math.max(maxSum, windowSum);\n      windowSum -= arr[i - k + 1];\n    }\n  }\n  return maxSum;\n}", whenToUse: "Maximum subarray, fixed window" },
  { id: "binary-search", name: "Binary Search", pattern: "O(log n) time, O(1) space", code: "function binarySearch(arr, target) {\n  let left = 0, right = arr.length - 1;\n  while (left <= right) {\n    const mid = Math.floor((left + right) / 2);\n    if (arr[mid] === target) return mid;\n    if (arr[mid] < target) left = mid + 1;\n    else right = mid - 1;\n  }\n  return -1;\n}", whenToUse: "Sorted data, search space" },
  { id: "dfs", name: "DFS", pattern: "O(V+E) time, O(V) space", code: "function dfs(graph, start, visited = new Set()) {\n  visited.add(start);\n  for (const neighbor of graph[start]) {\n    if (!visited.has(neighbor)) {\n      dfs(graph, neighbor, visited);\n    }\n  }\n}", whenToUse: "Tree traversal, cycle detection" },
  { id: "bfs", name: "BFS", pattern: "O(V+E) time, O(V) space", code: "function bfs(graph, start) {\n  const visited = new Set([start]);\n  const queue = [start];\n  while (queue.length > 0) {\n    const node = queue.shift();\n    for (const neighbor of graph[node]) {\n      if (!visited.has(neighbor)) {\n        visited.add(neighbor);\n        queue.push(neighbor);\n      }\n    }\n  }\n}", whenToUse: "Shortest path, level-order" }
];

const companyPatterns: Record<string, string[]> = {
  "Google": ["Two Pointers", "Sliding Window", "Binary Search", "DFS/BFS", "Dynamic Programming"],
  "Meta": ["Arrays", "Hash Tables", "Dynamic Programming", "Trees", "Linked Lists"],
  "Amazon": ["Trees", "Graphs", "Dynamic Programming", "Arrays", "System Design"],
  "Microsoft": ["Arrays", "Strings", "Dynamic Programming", "Trees", "System Design"],
  "Apple": ["Arrays", "Trees", "Dynamic Programming", "Graphs", "System Design"],
  "Netflix": ["Arrays", "Binary Search", "Dynamic Programming", "System Design"],
};

const bigOComplexities = [
  { name: "O(1)", description: "Constant", color: "bg-green-500" },
  { name: "O(log n)", description: "Logarithmic", color: "bg-green-400" },
  { name: "O(n)", description: "Linear", color: "bg-yellow-500" },
  { name: "O(n log n)", description: "Linearithmic", color: "bg-yellow-400" },
  { name: "O(n²)", description: "Quadratic", color: "bg-orange-500" },
  { name: "O(2ⁿ)", description: "Exponential", color: "bg-red-500" },
  { name: "O(n!)", description: "Factorial", color: "bg-red-600" },
];

export default function HackInterviewTools() {
  const [activeTab, setActiveTab] = useState("patterns");
  const [copiedId, setCopiedId] = useState<string | null>(null);
  const [timerMinutes, setTimerMinutes] = useState(30);
  const [timerSeconds, setTimerSeconds] = useState(0);
  const [isTimerRunning, setIsTimerRunning] = useState(false);
  const [showTimerAlert, setShowTimerAlert] = useState(false);
  const timerRef = useRef<NodeJS.Timeout | null>(null);
  const [codeInput, setCodeInput] = useState("");
  const [complexityResult, setComplexityResult] = useState<ComplexityResult | null>(null);

  useEffect(() => {
    if (isTimerRunning && (timerMinutes > 0 || timerSeconds > 0)) {
      timerRef.current = setInterval(() => {
        if (timerSeconds === 0) {
          if (timerMinutes === 0) { setIsTimerRunning(false); setShowTimerAlert(true); setTimeout(() => setShowTimerAlert(false), 5000); }
          else { setTimerMinutes(timerMinutes - 1); setTimerSeconds(59); }
        } else { setTimerSeconds(timerSeconds - 1); }
      }, 1000);
    }
    return () => { if (timerRef.current) clearInterval(timerRef.current); };
  }, [isTimerRunning, timerMinutes, timerSeconds]);

  const copyToClipboard = (text: string, id: string) => { navigator.clipboard.writeText(text); setCopiedId(id); setTimeout(() => setCopiedId(null), 2000); };
  const analyzeComplexity = () => {
    const code = codeInput.toLowerCase();
    let time = "O(1)", space = "O(1)", explanation = "";
    if (code.includes("for") && code.includes("for")) { time = "O(n²)"; explanation = "Nested loops"; }
    else if (code.includes("for") || code.includes("while")) { time = "O(n)"; explanation = "Single loop"; }
    setComplexityResult({ time, space, explanation: explanation || "Basic operations" });
  };
  const resetTimer = () => { setIsTimerRunning(false); setTimerMinutes(30); setTimerSeconds(0); };

  return (
    <div className="min-h-screen bg-gradient-to-b from-background via-accent/5 to-background">
      <div className="container mx-auto px-4 py-8 max-w-7xl">
        <div className="mb-8">
          <h1 className="text-4xl font-bold tracking-tight mb-2"><span className="bg-gradient-to-r from-primary to-blue-600 bg-clip-text text-transparent">Hack Interview Tools</span></h1>
          <p className="text-muted-foreground text-lg">Master coding interviews with these essential tools</p>
        </div>
        <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-6">
          <TabsList className="flex flex-wrap gap-2 h-auto p-2">
            <TabsTrigger value="patterns"><Code2 className="h-4 w-4 mr-2"/>Code Patterns</TabsTrigger>
            <TabsTrigger value="companies"><Target className="h-4 w-4 mr-2"/>Company Focus</TabsTrigger>
            <TabsTrigger value="timer"><Timer className="h-4 w-4 mr-2"/>Interview Timer</TabsTrigger>
            <TabsTrigger value="complexity"><BarChart3 className="h-4 w-4 mr-2"/>Complexity</TabsTrigger>
            <TabsTrigger value="cheatsheet"><BookOpen className="h-4 w-4 mr-2"/>Cheatsheet</TabsTrigger>
          </TabsList>
          <TabsContent value="patterns" className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {codeSnippets.map((snippet) => (
                <Card key={snippet.id} className="hover:shadow-lg">
                  <CardHeader className="pb-3">
                    <div className="flex items-center justify-between"><CardTitle className="text-lg">{snippet.name}</CardTitle><Badge variant="secondary">{snippet.pattern}</Badge></div>
                    <CardDescription className="text-xs">{snippet.whenToUse}</CardDescription>
                  </CardHeader>
                  <CardContent>
                    <pre className="bg-muted p-3 rounded-lg text-xs overflow-x-auto max-h-48"><code>{snippet.code}</code></pre>
                    <Button variant="outline" size="sm" className="mt-3 w-full" onClick={() => copyToClipboard(snippet.code, snippet.id)}>
                      {copiedId === snippet.id ? <><CheckCircle className="h-4 w-4 mr-2"/>Copied!</> : <><Copy className="h-4 w-4 mr-2"/>Copy Code</>}
                    </Button>
                  </CardContent>
                </Card>
              ))}
            </div>
          </TabsContent>
          <TabsContent value="companies" className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {Object.entries(companyPatterns).map(([company, topics]) => (
                <Card key={company} className="hover:shadow-lg">
                  <CardHeader><CardTitle className="text-xl flex items-center gap-2"><div className="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center"><Target className="h-4 w-4 text-primary"/></div>{company}</CardTitle></CardHeader>
                  <CardContent><div className="flex flex-wrap gap-2">{topics.map((topic, idx) => <Badge key={idx} variant="outline" className="text-xs">{topic}</Badge>)}</div></CardContent>
                </Card>
              ))}
            </div>
          </TabsContent>
          <TabsContent value="timer" className="space-y-6">
            <div className="max-w-2xl mx-auto">
              <Card>
                <CardHeader className="text-center"><CardTitle className="text-2xl">Interview Timer</CardTitle><CardDescription>Practice with a timed coding session</CardDescription></CardHeader>
                <CardContent className="space-y-6">
                  {showTimerAlert && <div className="bg-red-100 dark:bg-red-900/20 border border-red-500 text-red-700 px-4 py-3 rounded-lg flex items-center gap-2"><AlertCircle className="h-5 w-5"/>Time's up!</div>}
                  <div className="text-center"><div className="text-7xl font-mono font-bold tabular-nums">{String(timerMinutes).padStart(2,'0')}:{String(timerSeconds).padStart(2,'0')}</div></div>
                  <div className="flex justify-center gap-4">
                    <Button size="lg" onClick={() => setIsTimerRunning(!isTimerRunning)} className="w-32">{isTimerRunning ? <><Pause className="h-5 w-5 mr-2"/>Pause</> : <><Play className="h-5 w-5 mr-2"/>Start</>}</Button>
                    <Button variant="outline" size="lg" onClick={resetTimer} className="w-32"><RotateCcw className="h-5 w-5 mr-2"/>Reset</Button>
                  </div>
                  <div className="border-t pt-6"><div className="flex gap-2 justify-center">{[5,10,15,20,30,45,60].map(mins => <Button key={mins} variant={timerMinutes===mins?"default":"outline"} size="sm" onClick={()=>{setTimerMinutes(mins);setTimerSeconds(0);setIsTimerRunning(false);}}>{mins}m</Button>)}</div></div>
                </CardContent>
              </Card>
            </div>
          </TabsContent>
          <TabsContent value="complexity" className="space-y-6">
            <div className="max-w-4xl mx-auto">
              <Card>
                <CardHeader><CardTitle className="flex items-center gap-2"><Calculator className="h-5 w-5"/>Code Complexity Analyzer</CardTitle></CardHeader>
                <CardContent className="space-y-4">
                  <Textarea placeholder="Paste your code here..." value={codeInput} onChange={(e)=>setCodeInput(e.target.value)} className="font-mono min-h-[200px]"/>
                  <Button onClick={analyzeComplexity} className="w-full"><Zap className="h-4 w-4 mr-2"/>Analyze Complexity</Button>
                  {complexityResult && (
                    <div className="grid grid-cols-2 gap-4">
                      <div className="bg-muted p-4 rounded-lg"><div className="text-sm text-muted-foreground">Time</div><div className="text-2xl font-bold">{complexityResult.time}</div></div>
                      <div className="bg-muted p-4 rounded-lg"><div className="text-sm text-muted-foreground">Space</div><div className="text-2xl font-bold">{complexityResult.space}</div></div>
                      <div className="col-span-2 bg-muted p-4 rounded-lg"><div className="text-sm text-muted-foreground">Analysis</div><div className="text-sm">{complexityResult.explanation}</div></div>
                    </div>
                  )}
                </CardContent>
              </Card>
            </div>
          </TabsContent>
          <TabsContent value="cheatsheet" className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <Card><CardHeader><CardTitle className="flex items-center gap-2"><Clock className="h-5 w-5"/>Time Complexities</CardTitle></CardHeader>
                <CardContent className="space-y-3">{bigOComplexities.map(c => <div key={c.name} className="flex items-center gap-3"><div className={`w-3 h-3 rounded-full ${c.color}`}/><span className="font-mono font-bold">{c.name}</span><span className="text-muted-foreground ml-2">- {c.description}</span></div>)}</CardContent></Card>
              <Card><CardHeader><CardTitle className="flex items-center gap-2"><Database className="h-5 w-5"/>Data Structures</CardTitle></CardHeader>
                <CardContent className="space-y-2"><div className="flex justify-between"><span>Array</span><span className="text-muted-foreground">O(1) access</span></div><div className="flex justify-between"><span>Hash Map</span><span className="text-muted-foreground">O(1) avg</span></div><div className="flex justify-between"><span>Linked List</span><span className="text-muted-foreground">O(1) insert</span></div><div className="flex justify-between"><span>Binary Tree</span><span className="text-muted-foreground">O(log n)</span></div><div className="flex justify-between"><span>Heap</span><span className="text-muted-foreground">O(1) max/min</span></div><div className="flex justify-between"><span>Stack</span><span className="text-muted-foreground">O(1) push/pop</span></div></CardContent></Card>
              <Card className="md:col-span-2"><CardHeader><CardTitle className="flex items-center gap-2"><Lightbulb className="h-5 w-5"/>Quick Tips</CardTitle></CardHeader>
                <CardContent><div className="grid grid-cols-1 md:grid-cols-2 gap-4"><ul className="space-y-2 text-sm"><li className="flex items-center gap-2"><CheckCircle2 className="h-4 w-4 text-green-500"/>Start with brute force</li><li className="flex items-center gap-2"><CheckCircle2 className="h-4 w-4 text-green-500"/>Ask clarifying questions</li><li className="flex items-center gap-2"><CheckCircle2 className="h-4 w-4 text-green-500"/>Think about edge cases</li><li className="flex items-center gap-2"><CheckCircle2 className="h-4 w-4 text-green-500"/>Communicate your process</li></ul><ul className="space-y-2 text-sm"><li className="flex items-center gap-2"><CheckCircle2 className="h-4 w-4 text-green-500"/>Test with examples</li><li className="flex items-center gap-2"><CheckCircle2 className="h-4 w-4 text-green-500"/>Consider tradeoffs</li><li className="flex items-center gap-2"><CheckCircle2 className="h-4 w-4 text-green-500"/>Dry run your code</li><li className="flex items-center gap-2"><CheckCircle2 className="h-4 w-4 text-green-500"/>Know patterns by heart</li></ul></div></CardContent></Card>
            </div>
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
}

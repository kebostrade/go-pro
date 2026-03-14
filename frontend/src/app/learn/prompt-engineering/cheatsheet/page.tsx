"use client"

import Link from "next/link"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { BookOpen, Code2, Zap, AlertCircle, Check, X, ArrowLeft } from "lucide-react"

export default function CheatsheetPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-900 dark:to-gray-800">
      <div className="max-w-6xl mx-auto px-4 py-8">
        {/* Header */}
        <div className="mb-8">
          <Link href="/learn/prompt-engineering">
            <Button variant="ghost" size="sm" className="mb-4">
              <ArrowLeft className="h-4 w-4 mr-2" />
              Back to Course
            </Button>
          </Link>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">
            Prompt Engineering Cheatsheet
          </h1>
          <p className="text-gray-600 dark:text-gray-400">
            Quick reference guide for effective prompt engineering
          </p>
        </div>

        <div className="grid gap-6">
          {/* PIERS Framework */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-xl">
                <Zap className="h-5 w-5 text-blue-500" />
                PIERS Framework
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="overflow-x-auto">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-gray-200 dark:border-gray-700">
                      <th className="text-left py-3 px-4 font-semibold text-gray-700 dark:text-gray-300">Letter</th>
                      <th className="text-left py-3 px-4 font-semibold text-gray-700 dark:text-gray-300">Component</th>
                      <th className="text-left py-3 px-4 font-semibold text-gray-700 dark:text-gray-300">Description</th>
                      <th className="text-left py-3 px-4 font-semibold text-gray-700 dark:text-gray-300">Example</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
                    <tr className="hover:bg-gray-50 dark:hover:bg-gray-800/50">
                      <td className="py-3 px-4"><Badge variant="info">P</Badge></td>
                      <td className="py-3 px-4 font-medium">Purpose</td>
                      <td className="py-3 px-4 text-gray-600 dark:text-gray-400">What do you want?</td>
                      <td className="py-3 px-4 text-gray-500 dark:text-gray-500 italic">"Summarize this article"</td>
                    </tr>
                    <tr className="hover:bg-gray-50 dark:hover:bg-gray-800/50">
                      <td className="py-3 px-4"><Badge variant="info">I</Badge></td>
                      <td className="py-3 px-4 font-medium">Information</td>
                      <td className="py-3 px-4 text-gray-600 dark:text-gray-400">Context needed</td>
                      <td className="py-3 px-4 text-gray-500 dark:text-gray-500 italic">"For a business audience"</td>
                    </tr>
                    <tr className="hover:bg-gray-50 dark:hover:bg-gray-800/50">
                      <td className="py-3 px-4"><Badge variant="info">E</Badge></td>
                      <td className="py-3 px-4 font-medium">Examples</td>
                      <td className="py-3 px-4 text-gray-600 dark:text-gray-400">Show, don&apos;t tell</td>
                      <td className="py-3 px-4 text-gray-500 dark:text-gray-500 italic">Few-shot examples</td>
                    </tr>
                    <tr className="hover:bg-gray-50 dark:hover:bg-gray-800/50">
                      <td className="py-3 px-4"><Badge variant="info">R</Badge></td>
                      <td className="py-3 px-4 font-medium">Role</td>
                      <td className="py-3 px-4 text-gray-600 dark:text-gray-400">Who should AI be?</td>
                      <td className="py-3 px-4 text-gray-500 dark:text-gray-500 italic">"You are a senior editor"</td>
                    </tr>
                    <tr className="hover:bg-gray-50 dark:hover:bg-gray-800/50">
                      <td className="py-3 px-4"><Badge variant="info">S</Badge></td>
                      <td className="py-3 px-4 font-medium">Style</td>
                      <td className="py-3 px-4 text-gray-600 dark:text-gray-400">Output format</td>
                      <td className="py-3 px-4 text-gray-500 dark:text-gray-500 italic">"Return as JSON"</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </CardContent>
          </Card>

          {/* Prompt Patterns */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-xl">
                <Code2 className="h-5 w-5 text-blue-500" />
                Prompt Patterns
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="overflow-x-auto">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-gray-200 dark:border-gray-700">
                      <th className="text-left py-3 px-4 font-semibold text-gray-700 dark:text-gray-300">Pattern</th>
                      <th className="text-left py-3 px-4 font-semibold text-gray-700 dark:text-gray-300">Use Case</th>
                      <th className="text-left py-3 px-4 font-semibold text-gray-700 dark:text-gray-300">Template</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
                    <tr className="hover:bg-gray-50 dark:hover:bg-gray-800/50">
                      <td className="py-3 px-4"><Badge variant="secondary">Zero-Shot</Badge></td>
                      <td className="py-3 px-4 text-gray-600 dark:text-gray-400">Simple tasks</td>
                      <td className="py-3 px-4 font-mono text-xs bg-gray-100 dark:bg-gray-800 rounded px-2 py-1">[Task description] + [Input]</td>
                    </tr>
                    <tr className="hover:bg-gray-50 dark:hover:bg-gray-800/50">
                      <td className="py-3 px-4"><Badge variant="secondary">Few-Shot</Badge></td>
                      <td className="py-3 px-4 text-gray-600 dark:text-gray-400">Complex formatting</td>
                      <td className="py-3 px-4 font-mono text-xs bg-gray-100 dark:bg-gray-800 rounded px-2 py-1">Examples + Input</td>
                    </tr>
                    <tr className="hover:bg-gray-50 dark:hover:bg-gray-800/50">
                      <td className="py-3 px-4"><Badge variant="secondary">CoT</Badge></td>
                      <td className="py-3 px-4 text-gray-600 dark:text-gray-400">Reasoning tasks</td>
                      <td className="py-3 px-4 font-mono text-xs bg-gray-100 dark:bg-gray-800 rounded px-2 py-1">&quot;Let&apos;s think step by step.&quot;</td>
                    </tr>
                    <tr className="hover:bg-gray-50 dark:hover:bg-gray-800/50">
                      <td className="py-3 px-4"><Badge variant="secondary">ReAct</Badge></td>
                      <td className="py-3 px-4 text-gray-600 dark:text-gray-400">Tool-using agents</td>
                      <td className="py-3 px-4 font-mono text-xs bg-gray-100 dark:bg-gray-800 rounded px-2 py-1">Thought + Action + Observation</td>
                    </tr>
                    <tr className="hover:bg-gray-50 dark:hover:bg-gray-800/50">
                      <td className="py-3 px-4"><Badge variant="secondary">Tree of Thoughts</Badge></td>
                      <td className="py-3 px-4 text-gray-600 dark:text-gray-400">Complex decisions</td>
                      <td className="py-3 px-4 font-mono text-xs bg-gray-100 dark:bg-gray-800 rounded px-2 py-1">Explore multiple paths</td>
                    </tr>
                    <tr className="hover:bg-gray-50 dark:hover:bg-gray-800/50">
                      <td className="py-3 px-4"><Badge variant="secondary">Self-Consistency</Badge></td>
                      <td className="py-3 px-4 text-gray-600 dark:text-gray-400">High-stakes answers</td>
                      <td className="py-3 px-4 font-mono text-xs bg-gray-100 dark:bg-gray-800 rounded px-2 py-1">Multiple reasoning paths + vote</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </CardContent>
          </Card>

          <div className="grid md:grid-cols-2 gap-6">
            {/* Temperature Guide */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2 text-xl">
                  <AlertCircle className="h-5 w-5 text-blue-500" />
                  Temperature Guide
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  <div className="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-800 rounded-lg">
                    <div className="flex items-center gap-3">
                      <span className="font-mono font-bold text-blue-600">0.0</span>
                      <span className="text-sm text-gray-600 dark:text-gray-400">Code, factual answers</span>
                    </div>
                    <Badge variant="outline" className="text-xs">Deterministic</Badge>
                  </div>
                  <div className="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-800 rounded-lg">
                    <div className="flex items-center gap-3">
                      <span className="font-mono font-bold text-blue-600">0.3</span>
                      <span className="text-sm text-gray-600 dark:text-gray-400">Technical writing</span>
                    </div>
                    <Badge variant="outline" className="text-xs">Low variance</Badge>
                  </div>
                  <div className="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-800 rounded-lg">
                    <div className="flex items-center gap-3">
                      <span className="font-mono font-bold text-blue-600">0.5</span>
                      <span className="text-sm text-gray-600 dark:text-gray-400">Balanced responses</span>
                    </div>
                    <Badge variant="info" className="text-xs">Default</Badge>
                  </div>
                  <div className="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-800 rounded-lg">
                    <div className="flex items-center gap-3">
                      <span className="font-mono font-bold text-blue-600">0.7</span>
                      <span className="text-sm text-gray-600 dark:text-gray-400">General conversation</span>
                    </div>
                    <Badge variant="outline" className="text-xs">Moderate</Badge>
                  </div>
                  <div className="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-800 rounded-lg">
                    <div className="flex items-center gap-3">
                      <span className="font-mono font-bold text-blue-600">1.0</span>
                      <span className="text-sm text-gray-600 dark:text-gray-400">Creative writing</span>
                    </div>
                    <Badge variant="warning" className="text-xs">High variance</Badge>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Quick Formula */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2 text-xl">
                  <BookOpen className="h-5 w-5 text-blue-500" />
                  Quick Formula
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="p-4 bg-gray-900 dark:bg-gray-950 rounded-lg">
                    <p className="text-green-400 font-mono text-sm mb-2"># Perfect Prompt</p>
                    <p className="text-white font-mono text-sm">
                      Role + Context + Task + Constraints + Format + Examples
                    </p>
                  </div>
                  <div className="p-4 bg-gray-900 dark:bg-gray-950 rounded-lg">
                    <p className="text-green-400 font-mono text-sm mb-2"># Token Budget</p>
                    <p className="text-white font-mono text-sm">
                      Budget = Context Window - Input - Output Buffer
                    </p>
                    <p className="text-gray-400 font-mono text-xs mt-1">
                      Example: 200,000 - prompt_tokens - 1,000
                    </p>
                  </div>
                  <div className="p-4 bg-gray-900 dark:bg-gray-950 rounded-lg">
                    <p className="text-green-400 font-mono text-sm mb-2"># Cost Calculation</p>
                    <p className="text-white font-mono text-sm">
                      Cost = (input_tokens x input_rate) + (output_tokens x output_rate)
                    </p>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Best Practices */}
          <div className="grid md:grid-cols-2 gap-6">
            {/* DO */}
            <Card className="border-green-200 dark:border-green-900">
              <CardHeader>
                <CardTitle className="flex items-center gap-2 text-xl text-green-600 dark:text-green-400">
                  <Check className="h-5 w-5" />
                  DO
                </CardTitle>
              </CardHeader>
              <CardContent>
                <ul className="space-y-2">
                  {[
                    "Be specific and explicit",
                    "Provide examples for complex tasks",
                    "Define output format clearly",
                    "Use consistent structure",
                    "Include constraints",
                    "State what NOT to do",
                  ].map((item, i) => (
                    <li key={i} className="flex items-start gap-2 text-sm">
                      <Check className="h-4 w-4 text-green-500 mt-0.5 shrink-0" />
                      <span className="text-gray-700 dark:text-gray-300">{item}</span>
                    </li>
                  ))}
                </ul>
              </CardContent>
            </Card>

            {/* DON'T */}
            <Card className="border-red-200 dark:border-red-900">
              <CardHeader>
                <CardTitle className="flex items-center gap-2 text-xl text-red-600 dark:text-red-400">
                  <X className="h-5 w-5" />
                  DON&apos;T
                </CardTitle>
              </CardHeader>
              <CardContent>
                <ul className="space-y-2">
                  {[
                    'Be vague ("make it better")',
                    "Assume context is understood",
                    "Request multiple formats",
                    "Use ambiguous language",
                    "Skip edge cases",
                    "Over-constrain (too many rules)",
                  ].map((item, i) => (
                    <li key={i} className="flex items-start gap-2 text-sm">
                      <X className="h-4 w-4 text-red-500 mt-0.5 shrink-0" />
                      <span className="text-gray-700 dark:text-gray-300">{item}</span>
                    </li>
                  ))}
                </ul>
              </CardContent>
            </Card>
          </div>

          {/* Quick Templates */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-xl">
                <Code2 className="h-5 w-5 text-blue-500" />
                Quick Templates
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid md:grid-cols-2 gap-4">
                {/* Few-Shot Template */}
                <div>
                  <Badge variant="secondary" className="mb-2">Few-Shot</Badge>
                  <pre className="p-4 bg-gray-900 dark:bg-gray-950 rounded-lg text-sm overflow-x-auto">
                    <code className="text-gray-300">{`Example 1:
Input: [example]
Output: [example]

Example 2:
Input: [example]
Output: [example]

Now:
Input: [actual input]
Output:`}</code>
                  </pre>
                </div>

                {/* Chain-of-Thought Template */}
                <div>
                  <Badge variant="secondary" className="mb-2">Chain-of-Thought</Badge>
                  <pre className="p-4 bg-gray-900 dark:bg-gray-950 rounded-lg text-sm overflow-x-auto">
                    <code className="text-gray-300">{`[Task description]

Let's think step by step.
1. First, ...
2. Then, ...
3. Finally, ...

Answer:`}</code>
                  </pre>
                </div>

                {/* JSON Output Template */}
                <div>
                  <Badge variant="secondary" className="mb-2">JSON Output</Badge>
                  <pre className="p-4 bg-gray-900 dark:bg-gray-950 rounded-lg text-sm overflow-x-auto">
                    <code className="text-gray-300">{`Return as JSON:
{
  "field1": "type",
  "field2": ["array"],
  "field3": { "nested": "value" }
}

No markdown, only valid JSON.`}</code>
                  </pre>
                </div>

                {/* ReAct Template */}
                <div>
                  <Badge variant="secondary" className="mb-2">ReAct Agent</Badge>
                  <pre className="p-4 bg-gray-900 dark:bg-gray-950 rounded-lg text-sm overflow-x-auto">
                    <code className="text-gray-300">{`Thought: [reasoning]
Action: [tool_name]
Action Input: [parameters]
Observation: [result]
[repeat as needed]
Final Answer: [response]`}</code>
                  </pre>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Debugging Prompts */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-xl">
                <AlertCircle className="h-5 w-5 text-blue-500" />
                Debugging Prompts
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid md:grid-cols-2 gap-4">
                <div className="p-4 bg-gray-50 dark:bg-gray-800 rounded-lg">
                  <p className="font-medium text-sm text-gray-700 dark:text-gray-300 mb-2">Output too long?</p>
                  <code className="text-xs bg-gray-200 dark:bg-gray-700 px-2 py-1 rounded">Keep response under [N] words. Be concise.</code>
                </div>
                <div className="p-4 bg-gray-50 dark:bg-gray-800 rounded-lg">
                  <p className="font-medium text-sm text-gray-700 dark:text-gray-300 mb-2">Wrong format?</p>
                  <code className="text-xs bg-gray-200 dark:bg-gray-700 px-2 py-1 rounded">Output ONLY [format]. No explanations.</code>
                </div>
                <div className="p-4 bg-gray-50 dark:bg-gray-800 rounded-lg">
                  <p className="font-medium text-sm text-gray-700 dark:text-gray-300 mb-2">Inconsistent?</p>
                  <code className="text-xs bg-gray-200 dark:bg-gray-700 px-2 py-1 rounded">Follow this EXACT format: [show format]</code>
                </div>
                <div className="p-4 bg-gray-50 dark:bg-gray-800 rounded-lg">
                  <p className="font-medium text-sm text-gray-700 dark:text-gray-300 mb-2">Task misunderstood?</p>
                  <code className="text-xs bg-gray-200 dark:bg-gray-700 px-2 py-1 rounded">Your goal is [action]. NOT: [what it&apos;s not].</code>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Footer */}
        <div className="mt-8 pt-6 border-t border-gray-200 dark:border-gray-700 text-center text-sm text-gray-500 dark:text-gray-400">
          <p>Last Updated: March 2026</p>
        </div>
      </div>
    </div>
  )
}

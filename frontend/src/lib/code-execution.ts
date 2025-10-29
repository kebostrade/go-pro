/**
 * Enhanced Go Code Execution Engine
 * Provides advanced simulation of Go code execution with better parsing and analysis
 */

export interface ExecutionResult {
  output: string;
  error?: string;
  executionTime: number;
  memoryUsed?: number;
  complexity?: {
    time: string;
    space: string;
  };
  warnings?: string[];
  suggestions?: string[];
}

export interface CodeAnalysis {
  lineCount: number;
  functionCount: number;
  loopCount: number;
  conditionalCount: number;
  complexity: number;
  hasMain: boolean;
  imports: string[];
  functions: string[];
  variables: string[];
}

/**
 * Analyze Go code structure and complexity
 */
export function analyzeCode(code: string): CodeAnalysis {
  const lines = code.split('\n').filter(line => line.trim() && !line.trim().startsWith('//'));
  
  const functionMatches = code.match(/func\s+\w+\s*\([^)]*\)/g) || [];
  const loopMatches = code.match(/\b(for|range)\b/g) || [];
  const conditionalMatches = code.match(/\b(if|else|switch|case)\b/g) || [];
  const importMatches = code.match(/import\s+(?:"[^"]+"|`[^`]+`|\([^)]+\))/g) || [];
  
  // Extract function names
  const functions = functionMatches.map(match => {
    const nameMatch = match.match(/func\s+(\w+)/);
    return nameMatch ? nameMatch[1] : '';
  }).filter(Boolean);
  
  // Extract imports
  const imports: string[] = [];
  importMatches.forEach(match => {
    const singleImport = match.match(/"([^"]+)"/g);
    if (singleImport) {
      imports.push(...singleImport.map(imp => imp.replace(/"/g, '')));
    }
  });
  
  // Calculate cyclomatic complexity (simplified)
  const complexity = 1 + loopMatches.length + conditionalMatches.length;
  
  return {
    lineCount: lines.length,
    functionCount: functions.length,
    loopCount: loopMatches.length,
    conditionalCount: conditionalMatches.length,
    complexity,
    hasMain: code.includes('func main()'),
    imports,
    functions,
    variables: extractVariables(code),
  };
}

/**
 * Extract variable declarations from code
 */
function extractVariables(code: string): string[] {
  const varMatches = code.match(/(?:var|:=)\s+(\w+)/g) || [];
  return varMatches.map(match => {
    const nameMatch = match.match(/(?:var|:=)\s+(\w+)/);
    return nameMatch ? nameMatch[1] : '';
  }).filter(Boolean);
}

/**
 * Enhanced Go code execution simulator
 */
export async function executeGoCode(
  code: string,
  inputs: Record<string, any> = {}
): Promise<ExecutionResult> {
  const startTime = performance.now();
  
  // Analyze code first
  const analysis = analyzeCode(code);
  
  // Simulate realistic execution time based on complexity
  const baseTime = 50;
  const complexityTime = analysis.complexity * 20;
  const loopTime = analysis.loopCount * 30;
  const executionTime = baseTime + complexityTime + loopTime + Math.random() * 100;
  
  await new Promise(resolve => setTimeout(resolve, executionTime));
  
  let output = "";
  let error: string | undefined;
  const warnings: string[] = [];
  const suggestions: string[] = [];
  
  try {
    // Check for common errors
    if (!analysis.hasMain && code.includes('package main')) {
      warnings.push('Package main declared but no main() function found');
    }
    
    // Simulate fmt.Println outputs
    output = simulatePrintStatements(code, inputs);
    
    // Simulate mathematical operations
    output += simulateMathOperations(code, inputs);
    
    // Simulate loops
    output += simulateLoops(code, inputs);
    
    // Simulate conditionals
    output += simulateConditionals(code, inputs);
    
    // Simulate string operations
    output += simulateStringOperations(code, inputs);
    
    // Simulate array/slice operations
    output += simulateArrayOperations(code, inputs);
    
    // Add suggestions based on code analysis
    if (analysis.complexity > 10) {
      suggestions.push('Consider refactoring: High cyclomatic complexity detected');
    }
    
    if (analysis.loopCount > 3) {
      suggestions.push('Multiple loops detected: Consider optimizing for performance');
    }
    
    if (!code.includes('error')) {
      suggestions.push('Consider adding error handling for production code');
    }
    
    // Default output if nothing specific found
    if (!output.trim()) {
      output = "Program executed successfully\n";
      if (Object.keys(inputs).length > 0) {
        output += `Inputs processed: ${JSON.stringify(inputs)}\n`;
      }
    }
    
  } catch (e) {
    error = `Runtime error: ${e instanceof Error ? e.message : 'Unknown error'}`;
    output = "";
  }
  
  const endTime = performance.now();
  const actualExecutionTime = Math.round(endTime - startTime);
  
  // Calculate estimated memory usage based on code analysis
  const baseMemory = 1024; // 1KB base
  const variableMemory = analysis.variables.length * 64;
  const functionMemory = analysis.functionCount * 128;
  const memoryUsed = baseMemory + variableMemory + functionMemory + Math.round(Math.random() * 512);
  
  // Estimate complexity
  const timeComplexity = estimateTimeComplexity(analysis);
  const spaceComplexity = estimateSpaceComplexity(analysis);
  
  return {
    output: output.trim(),
    error,
    executionTime: actualExecutionTime,
    memoryUsed,
    complexity: {
      time: timeComplexity,
      space: spaceComplexity,
    },
    warnings: warnings.length > 0 ? warnings : undefined,
    suggestions: suggestions.length > 0 ? suggestions : undefined,
  };
}

/**
 * Simulate fmt.Println and fmt.Printf statements
 */
function simulatePrintStatements(code: string, inputs: Record<string, any>): string {
  let output = "";
  
  // Match fmt.Println
  const printlnMatches = code.match(/fmt\.Println\([^)]+\)/g) || [];
  printlnMatches.forEach(match => {
    let content = match.replace(/fmt\.Println\(|\)/g, '');
    content = replaceVariables(content, inputs);
    content = content.replace(/"/g, '').replace(/`/g, '');
    output += content + '\n';
  });
  
  // Match fmt.Printf
  const printfMatches = code.match(/fmt\.Printf\([^)]+\)/g) || [];
  printfMatches.forEach(match => {
    let content = match.replace(/fmt\.Printf\(|\)/g, '');
    content = replaceVariables(content, inputs);
    // Simple format string handling
    content = content.replace(/%d/g, '').replace(/%s/g, '').replace(/%v/g, '');
    content = content.replace(/"/g, '').replace(/`/g, '');
    output += content + '\n';
  });
  
  return output;
}

/**
 * Simulate mathematical operations
 */
function simulateMathOperations(code: string, inputs: Record<string, any>): string {
  let output = "";
  
  const numbers = Object.values(inputs).filter(v => typeof v === 'number');
  if (numbers.length >= 2) {
    const [a, b] = numbers as number[];
    
    if (code.includes('+') && code.includes('fmt.Println')) {
      output += `${a + b}\n`;
    }
    if (code.includes('-') && code.includes('fmt.Println')) {
      output += `${a - b}\n`;
    }
    if (code.includes('*') && code.includes('fmt.Println')) {
      output += `${a * b}\n`;
    }
    if (code.includes('/') && b !== 0 && code.includes('fmt.Println')) {
      output += `${a / b}\n`;
    }
  }
  
  return output;
}

/**
 * Simulate loop execution
 */
function simulateLoops(code: string, inputs: Record<string, any>): string {
  let output = "";
  
  // Simulate for loops with range
  const rangeMatch = code.match(/for\s+\w+\s*:=\s*range\s+(\d+)/);
  if (rangeMatch) {
    const count = Math.min(parseInt(rangeMatch[1]), 10);
    for (let i = 0; i < count; i++) {
      output += `Iteration ${i}\n`;
    }
  }
  
  // Simulate for loops with condition
  const forMatch = code.match(/for\s+\w+\s*:=\s*(\d+);\s*\w+\s*<\s*(\d+)/);
  if (forMatch) {
    const start = parseInt(forMatch[1]);
    const end = Math.min(parseInt(forMatch[2]), start + 10);
    for (let i = start; i < end; i++) {
      output += `Loop ${i}\n`;
    }
  }
  
  return output;
}

/**
 * Simulate conditional statements
 */
function simulateConditionals(code: string, inputs: Record<string, any>): string {
  let output = "";
  
  // Simple if-else simulation
  if (code.includes('if') && code.includes('else')) {
    const numbers = Object.values(inputs).filter(v => typeof v === 'number');
    if (numbers.length >= 2) {
      const [a, b] = numbers as number[];
      if (a > b) {
        output += `${a} is greater than ${b}\n`;
      } else {
        output += `${b} is greater than or equal to ${a}\n`;
      }
    }
  }
  
  return output;
}

/**
 * Simulate string operations
 */
function simulateStringOperations(code: string, inputs: Record<string, any>): string {
  let output = "";
  
  const strings = Object.values(inputs).filter(v => typeof v === 'string');
  if (strings.length > 0 && code.includes('strings.')) {
    if (code.includes('strings.ToUpper')) {
      output += strings.map(s => (s as string).toUpperCase()).join('\n') + '\n';
    }
    if (code.includes('strings.ToLower')) {
      output += strings.map(s => (s as string).toLowerCase()).join('\n') + '\n';
    }
  }
  
  return output;
}

/**
 * Simulate array/slice operations
 */
function simulateArrayOperations(code: string, inputs: Record<string, any>): string {
  let output = "";
  
  if (code.includes('[]') && code.includes('append')) {
    output += "Array/slice operations executed\n";
  }
  
  return output;
}

/**
 * Replace variables in code with actual values
 */
function replaceVariables(content: string, inputs: Record<string, any>): string {
  Object.entries(inputs).forEach(([key, value]) => {
    const regex = new RegExp(`\\b${key}\\b`, 'g');
    content = content.replace(regex, String(value));
  });
  return content;
}

/**
 * Estimate time complexity based on code analysis
 */
function estimateTimeComplexity(analysis: CodeAnalysis): string {
  if (analysis.loopCount === 0) return 'O(1)';
  if (analysis.loopCount === 1) return 'O(n)';
  if (analysis.loopCount === 2) return 'O(n²)';
  if (analysis.loopCount >= 3) return 'O(n³)';
  return 'O(n)';
}

/**
 * Estimate space complexity based on code analysis
 */
function estimateSpaceComplexity(analysis: CodeAnalysis): string {
  if (analysis.variables.length === 0) return 'O(1)';
  // Check if code uses arrays or slices
  const hasArrays = analysis.variables.some(v => v.includes('[]'));
  if (hasArrays) return 'O(n)';
  return 'O(1)';
}


"use client";

import { useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { 
  Code2, 
  Copy, 
  Check, 
  Search,
  Filter,
  BookOpen,
  Zap,
  Database,
  Network,
  Lock,
  FileJson,
  Terminal
} from "lucide-react";

interface CodeSnippet {
  id: string;
  title: string;
  description: string;
  code: string;
  category: 'basics' | 'data-structures' | 'algorithms' | 'concurrency' | 'web' | 'database' | 'testing';
  tags: string[];
  difficulty: 'beginner' | 'intermediate' | 'advanced';
}

const SNIPPETS: CodeSnippet[] = [
  {
    id: 'hello-world',
    title: 'Hello World',
    description: 'Basic Go program structure',
    code: `package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}`,
    category: 'basics',
    tags: ['beginner', 'syntax'],
    difficulty: 'beginner',
  },
  {
    id: 'variables',
    title: 'Variable Declaration',
    description: 'Different ways to declare variables in Go',
    code: `package main

import "fmt"

func main() {
    // Method 1: var keyword
    var name string = "John"
    
    // Method 2: type inference
    var age = 30
    
    // Method 3: short declaration
    city := "New York"
    
    fmt.Println(name, age, city)
}`,
    category: 'basics',
    tags: ['variables', 'syntax'],
    difficulty: 'beginner',
  },
  {
    id: 'for-loop',
    title: 'For Loop',
    description: 'Different for loop patterns in Go',
    code: `package main

import "fmt"

func main() {
    // Classic for loop
    for i := 0; i < 5; i++ {
        fmt.Println(i)
    }
    
    // While-style loop
    j := 0
    for j < 5 {
        fmt.Println(j)
        j++
    }
    
    // Infinite loop with break
    k := 0
    for {
        if k >= 5 {
            break
        }
        fmt.Println(k)
        k++
    }
}`,
    category: 'basics',
    tags: ['loops', 'control-flow'],
    difficulty: 'beginner',
  },
  {
    id: 'slice-operations',
    title: 'Slice Operations',
    description: 'Common slice operations and manipulations',
    code: `package main

import "fmt"

func main() {
    // Create slice
    numbers := []int{1, 2, 3, 4, 5}
    
    // Append
    numbers = append(numbers, 6)
    
    // Slice
    subset := numbers[1:4]
    
    // Length and capacity
    fmt.Println("Length:", len(numbers))
    fmt.Println("Capacity:", cap(numbers))
    
    // Iterate
    for i, num := range numbers {
        fmt.Printf("Index %d: %d\\n", i, num)
    }
}`,
    category: 'data-structures',
    tags: ['slices', 'arrays'],
    difficulty: 'beginner',
  },
  {
    id: 'map-operations',
    title: 'Map Operations',
    description: 'Working with maps (hash tables)',
    code: `package main

import "fmt"

func main() {
    // Create map
    ages := make(map[string]int)
    
    // Add entries
    ages["Alice"] = 25
    ages["Bob"] = 30
    
    // Check if key exists
    age, exists := ages["Alice"]
    if exists {
        fmt.Println("Alice's age:", age)
    }
    
    // Delete entry
    delete(ages, "Bob")
    
    // Iterate
    for name, age := range ages {
        fmt.Printf("%s is %d years old\\n", name, age)
    }
}`,
    category: 'data-structures',
    tags: ['maps', 'hash-table'],
    difficulty: 'beginner',
  },
  {
    id: 'struct-definition',
    title: 'Struct Definition',
    description: 'Define and use custom types with structs',
    code: `package main

import "fmt"

type Person struct {
    Name string
    Age  int
    City string
}

func main() {
    // Create struct
    person := Person{
        Name: "Alice",
        Age:  25,
        City: "NYC",
    }
    
    // Access fields
    fmt.Println(person.Name)
    
    // Modify fields
    person.Age = 26
    
    fmt.Printf("%+v\\n", person)
}`,
    category: 'basics',
    tags: ['structs', 'types'],
    difficulty: 'beginner',
  },
  {
    id: 'goroutine-basic',
    title: 'Basic Goroutine',
    description: 'Launch concurrent goroutines',
    code: `package main

import (
    "fmt"
    "time"
)

func sayHello(name string) {
    for i := 0; i < 3; i++ {
        fmt.Printf("Hello, %s!\\n", name)
        time.Sleep(100 * time.Millisecond)
    }
}

func main() {
    // Launch goroutine
    go sayHello("Alice")
    go sayHello("Bob")
    
    // Wait for goroutines
    time.Sleep(500 * time.Millisecond)
}`,
    category: 'concurrency',
    tags: ['goroutines', 'concurrency'],
    difficulty: 'intermediate',
  },
  {
    id: 'channel-basic',
    title: 'Basic Channel',
    description: 'Use channels for goroutine communication',
    code: `package main

import "fmt"

func sum(numbers []int, ch chan int) {
    total := 0
    for _, num := range numbers {
        total += num
    }
    ch <- total // Send to channel
}

func main() {
    numbers := []int{1, 2, 3, 4, 5}
    ch := make(chan int)
    
    go sum(numbers[:len(numbers)/2], ch)
    go sum(numbers[len(numbers)/2:], ch)
    
    x, y := <-ch, <-ch // Receive from channel
    
    fmt.Println("Total:", x+y)
}`,
    category: 'concurrency',
    tags: ['channels', 'goroutines'],
    difficulty: 'intermediate',
  },
  {
    id: 'error-handling',
    title: 'Error Handling',
    description: 'Proper error handling pattern',
    code: `package main

import (
    "errors"
    "fmt"
)

func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

func main() {
    result, err := divide(10, 2)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Result:", result)
}`,
    category: 'basics',
    tags: ['errors', 'best-practices'],
    difficulty: 'beginner',
  },
  {
    id: 'http-server',
    title: 'Simple HTTP Server',
    description: 'Create a basic HTTP server',
    code: `package main

import (
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}

func main() {
    http.HandleFunc("/", handler)
    fmt.Println("Server starting on :8080")
    http.ListenAndServe(":8080", nil)
}`,
    category: 'web',
    tags: ['http', 'server', 'web'],
    difficulty: 'intermediate',
  },
];

interface CodeSnippetsLibraryProps {
  onInsert?: (code: string) => void;
}

export default function CodeSnippetsLibrary({ onInsert }: CodeSnippetsLibraryProps) {
  const [searchTerm, setSearchTerm] = useState("");
  const [selectedCategory, setSelectedCategory] = useState<string | null>(null);
  const [copiedId, setCopiedId] = useState<string | null>(null);
  
  const categories = Array.from(new Set(SNIPPETS.map(s => s.category)));
  
  const filteredSnippets = SNIPPETS.filter(snippet => {
    const matchesSearch = !searchTerm || 
      snippet.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
      snippet.description.toLowerCase().includes(searchTerm.toLowerCase()) ||
      snippet.tags.some(tag => tag.toLowerCase().includes(searchTerm.toLowerCase()));
    
    const matchesCategory = !selectedCategory || snippet.category === selectedCategory;
    
    return matchesSearch && matchesCategory;
  });
  
  const handleCopy = async (snippet: CodeSnippet) => {
    await navigator.clipboard.writeText(snippet.code);
    setCopiedId(snippet.id);
    setTimeout(() => setCopiedId(null), 2000);
  };
  
  const getCategoryIcon = (category: string) => {
    switch (category) {
      case 'basics': return BookOpen;
      case 'data-structures': return Database;
      case 'algorithms': return Zap;
      case 'concurrency': return Network;
      case 'web': return Terminal;
      case 'database': return Database;
      case 'testing': return FileJson;
      default: return Code2;
    }
  };
  
  const getDifficultyColor = (difficulty: string) => {
    switch (difficulty) {
      case 'beginner': return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200';
      case 'intermediate': return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200';
      case 'advanced': return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200';
      default: return 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-200';
    }
  };

  return (
    <Card className="glass-card border-2">
      <CardHeader>
        <CardTitle className="flex items-center">
          <Code2 className="mr-2 h-5 w-5 text-primary" />
          Code Snippets Library
        </CardTitle>
        <CardDescription>
          Quick access to common Go patterns and examples
        </CardDescription>
      </CardHeader>
      
      <CardContent className="space-y-4">
        {/* Search and Filter */}
        <div className="space-y-3">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Search snippets..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="pl-10"
            />
          </div>
          
          <div className="flex flex-wrap gap-2">
            <Button
              variant={selectedCategory === null ? 'default' : 'outline'}
              size="sm"
              onClick={() => setSelectedCategory(null)}
            >
              All
            </Button>
            {categories.map(category => {
              const Icon = getCategoryIcon(category);
              return (
                <Button
                  key={category}
                  variant={selectedCategory === category ? 'default' : 'outline'}
                  size="sm"
                  onClick={() => setSelectedCategory(category)}
                >
                  <Icon className="mr-2 h-3 w-3" />
                  {category}
                </Button>
              );
            })}
          </div>
        </div>
        
        {/* Snippets Grid */}
        <div className="grid grid-cols-1 gap-3 max-h-[600px] overflow-y-auto custom-scrollbar">
          {filteredSnippets.map(snippet => (
            <Card key={snippet.id} className="hover:border-primary/50 transition-colors">
              <CardContent className="p-4">
                <div className="space-y-3">
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <h4 className="font-medium text-sm mb-1">{snippet.title}</h4>
                      <p className="text-xs text-muted-foreground">{snippet.description}</p>
                    </div>
                    <Badge className={getDifficultyColor(snippet.difficulty)}>
                      {snippet.difficulty}
                    </Badge>
                  </div>
                  
                  <pre className="text-xs bg-muted p-3 rounded overflow-x-auto max-h-32">
                    <code>{snippet.code}</code>
                  </pre>
                  
                  <div className="flex items-center justify-between">
                    <div className="flex flex-wrap gap-1">
                      {snippet.tags.slice(0, 3).map(tag => (
                        <Badge key={tag} variant="outline" className="text-xs">
                          {tag}
                        </Badge>
                      ))}
                    </div>
                    
                    <div className="flex items-center space-x-2">
                      {onInsert && (
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => onInsert(snippet.code)}
                        >
                          Insert
                        </Button>
                      )}
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => handleCopy(snippet)}
                      >
                        {copiedId === snippet.id ? (
                          <Check className="h-4 w-4 text-green-500" />
                        ) : (
                          <Copy className="h-4 w-4" />
                        )}
                      </Button>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
        
        {filteredSnippets.length === 0 && (
          <div className="text-center py-8 text-muted-foreground">
            <Code2 className="mx-auto h-12 w-12 mb-4 opacity-50" />
            <p>No snippets found</p>
            <p className="text-sm">Try adjusting your search or filters</p>
          </div>
        )}
      </CardContent>
    </Card>
  );
}


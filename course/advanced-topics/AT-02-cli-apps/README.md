# Building CLI Applications with Go

Develop powerful command-line tools using Go's flag package and cobra.

## Learning Objectives

- Parse command-line arguments
- Structure CLI applications
- Handle configuration files
- Implement subcommands
- Add colored output and progress bars
- Distribute CLI tools

## Theory

### Basic CLI with flag

```go
func main() {
    name := flag.String("name", "World", "a name to greet")
    count := flag.Int("count", 1, "number of greetings")
    flag.Parse()

    for i := 0; i < *count; i++ {
        fmt.Printf("Hello, %s!\n", *name)
    }
}
```

### Cobra Framework

```go
var rootCmd = &cobra.Command{
    Use:   "myapp",
    Short: "My CLI application",
    Long:  "A longer description of my application",
}

var greetCmd = &cobra.Command{
    Use:   "greet [name]",
    Short: "Greet someone",
    Args:  cobra.MaximumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        name := "World"
        if len(args) > 0 {
            name = args[0]
        }
        fmt.Printf("Hello, %s!\n", name)
    },
}

func init() {
    rootCmd.AddCommand(greetCmd)
}

func main() {
    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}
```

### Configuration Management

```go
type Config struct {
    APIKey  string `yaml:"api_key" env:"API_KEY"`
    Timeout int    `yaml:"timeout" env:"TIMEOUT"`
    Debug   bool   `yaml:"debug" env:"DEBUG"`
}

func LoadConfig(path string) (*Config, error) {
    cfg := &Config{Timeout: 30}

    if data, err := os.ReadFile(path); err == nil {
        if err := yaml.Unmarshal(data, cfg); err != nil {
            return nil, err
        }
    }

    if err := env.Parse(cfg); err != nil {
        return nil, err
    }

    return cfg, nil
}
```

### Colored Output

```go
import "github.com/fatih/color"

var (
    success = color.New(color.FgGreen).Add(color.Bold)
    errorC  = color.New(color.FgRed).Add(color.Bold)
    info    = color.New(color.FgCyan)
)

func printResult(msg string, err error) {
    if err != nil {
        errorC.Printf("Error: %v\n", err)
        os.Exit(1)
    }
    success.Printf("Success: %s\n", msg)
}
```

### Progress Indicators

```go
import "github.com/schollz/progressbar/v3"

func processFiles(files []string) error {
    bar := progressbar.Default(int64(len(files)))
    for _, f := range files {
        if err := processFile(f); err != nil {
            return err
        }
        bar.Add(1)
    }
    return nil
}
```

## Security Considerations

```go
func maskSensitive(key string) string {
    if len(key) <= 8 {
        return "****"
    }
    return key[:4] + "****" + key[len(key)-4:]
}

func validatePath(path string) error {
    clean := filepath.Clean(path)
    if strings.Contains(clean, "..") {
        return errors.New("invalid path: directory traversal detected")
    }
    return nil
}
```

## Performance Tips

```go
func main() {
    if len(os.Args) < 2 {
        fmt.Fprintln(os.Stderr, "usage: myapp <command>")
        os.Exit(1)
    }

    switch os.Args[1] {
    case "version":
        fmt.Println("v1.0.0")
        os.Exit(0)
    }

    rootCmd.Execute()
}
```

## Exercises

1. Build a file converter CLI
2. Create a task manager with subcommands
3. Add configuration file support
4. Implement shell completion

## Validation

```bash
cd exercises
go build -o myapp
./myapp --help
./myapp greet World
```

## Key Takeaways

- Use cobra for complex CLIs
- Support multiple config sources
- Provide helpful error messages
- Include shell completion
- Version your releases

## Next Steps

**[AT-03: Testing & Debugging](../AT-03-testing-debugging/README.md)**

---

CLI tools developers love. ⌨️

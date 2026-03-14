package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

// ============================================================================
// TASK MODEL
// ============================================================================

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Priority    string    `json:"priority"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
}

type TaskList struct {
	Tasks []Task `json:"tasks"`
}

// ============================================================================
// FILE STORAGE
// ============================================================================

const taskFile = "tasks.json"

func loadTasks() (*TaskList, error) {
	data, err := os.ReadFile(taskFile)
	if err != nil {
		if os.IsNotExist(err) {
			return &TaskList{Tasks: []Task{}}, nil
		}
		return nil, err
	}

	var tasks TaskList
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}

	return &tasks, nil
}

func saveTasks(tasks *TaskList) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(taskFile, data, 0644)
}

// ============================================================================
// COMMANDS
// ============================================================================

func cmdAdd(args []string) error {
	flags := flag.NewFlagSet("add", flag.ExitOnError)
	description := flags.String("desc", "", "Task description")
	priority := flags.String("priority", "medium", "Task priority (low, medium, high)")

	if err := flags.Parse(args); err != nil {
		return err
	}

	title := flags.Arg(0)
	if title == "" {
		return fmt.Errorf("task title is required")
	}

	tasks, _ := loadTasks()

	task := Task{
		ID:          fmt.Sprintf("%d", time.Now().UnixNano()),
		Title:       title,
		Description: *description,
		Priority:    *priority,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	tasks.Tasks = append(tasks.Tasks, task)

	if err := saveTasks(tasks); err != nil {
		return fmt.Errorf("failed to save task: %w", err)
	}

	fmt.Printf("✅ Task added: %s (ID: %s)\n", task.Title, task.ID)
	return nil
}

func cmdList(args []string) error {
	flags := flag.NewFlagSet("list", flag.ExitOnError)
	status := flags.String("status", "", "Filter by status (pending, completed)")
	priority := flags.String("priority", "", "Filter by priority (low, medium, high)")

	if err := flags.Parse(args); err != nil {
		return err
	}

	tasks, err := loadTasks()
	if err != nil {
		return fmt.Errorf("failed to load tasks: %w", err)
	}

	filtered := filterTasks(tasks.Tasks, *status, *priority)

	if len(filtered) == 0 {
		fmt.Println("📝 No tasks found")
		return nil
	}

	fmt.Printf("\n📋 Tasks (%d):\n\n", len(filtered))
	for _, task := range filtered {
		statusIcon := "⏳"
		if task.Status == "completed" {
			statusIcon = "✅"
		}

		priorityIcon := "🔵"
		if task.Priority == "high" {
			priorityIcon = "🔴"
		} else if task.Priority == "medium" {
			priorityIcon = "🟡"
		}

		fmt.Printf("%s [%s] %s %s\n", statusIcon, task.ID, priorityIcon, task.Title)
		if task.Description != "" {
			fmt.Printf("   %s\n", task.Description)
		}
		fmt.Printf("   Created: %s\n\n", task.CreatedAt.Format("2006-01-02 15:04"))
	}

	return nil
}

func cmdComplete(args []string) error {
	flags := flag.NewFlagSet("complete", flag.ExitOnError)

	if err := flags.Parse(args); err != nil {
		return err
	}

	taskID := flags.Arg(0)
	if taskID == "" {
		return fmt.Errorf("task ID is required")
	}

	tasks, _ := loadTasks()

	found := false
	for i, task := range tasks.Tasks {
		if task.ID == taskID {
			tasks.Tasks[i].Status = "completed"
			tasks.Tasks[i].CompletedAt = time.Now()
			found = true
			fmt.Printf("✅ Task completed: %s\n", task.Title)
			break
		}
	}

	if !found {
		return fmt.Errorf("task not found: %s", taskID)
	}

	if err := saveTasks(tasks); err != nil {
		return fmt.Errorf("failed to save tasks: %w", err)
	}

	return nil
}

func cmdDelete(args []string) error {
	flags := flag.NewFlagSet("delete", flag.ExitOnError)

	if err := flags.Parse(args); err != nil {
		return err
	}

	taskID := flags.Arg(0)
	if taskID == "" {
		return fmt.Errorf("task ID is required")
	}

	tasks, _ := loadTasks()

	found := false
	var updated []Task
	for _, task := range tasks.Tasks {
		if task.ID == taskID {
			found = true
			fmt.Printf("🗑️  Task deleted: %s\n", task.Title)
		} else {
			updated = append(updated, task)
		}
	}

	if !found {
		return fmt.Errorf("task not found: %s", taskID)
	}

	tasks.Tasks = updated

	if err := saveTasks(tasks); err != nil {
		return fmt.Errorf("failed to save tasks: %w", err)
	}

	return nil
}

func cmdClear(args []string) error {
	tasks := &TaskList{Tasks: []Task{}}

	if err := saveTasks(tasks); err != nil {
		return fmt.Errorf("failed to clear tasks: %w", err)
	}

	fmt.Println("🗑️  All tasks cleared")
	return nil
}

func filterTasks(tasks []Task, status, priority string) []Task {
	var filtered []Task

	for _, task := range tasks {
		if status != "" && task.Status != status {
			continue
		}
		if priority != "" && task.Priority != priority {
			continue
		}
		filtered = append(filtered, task)
	}

	return filtered
}

// ============================================================================
// MAIN
// ============================================================================

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	var err error
	switch command {
	case "add":
		err = cmdAdd(args)
	case "list", "ls":
		err = cmdList(args)
	case "complete", "done":
		err = cmdComplete(args)
	case "delete", "rm":
		err = cmdDelete(args)
	case "clear":
		err = cmdClear(args)
	case "help", "-h", "--help":
		printUsage()
		return
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	usage := `
📝 Task CLI - A simple task management tool

USAGE:
  task-cli <command> [arguments]

COMMANDS:
  add <title> [flags]       Add a new task
  list, ls [flags]          List all tasks
  complete, done <id>       Mark task as completed
  delete, rm <id>           Delete a task
  clear                     Delete all tasks
  help                      Show this help message

FLAGS:
  -desc <description>       Task description
  -priority <level>         Task priority (low, medium, high)
  -status <status>          Filter by status (pending, completed)

EXAMPLES:
  # Add a task
  task-cli add "Buy groceries" -desc "Milk, eggs, bread" -priority high

  # List all tasks
  task-cli list

  # List only pending tasks
  task-cli list -status pending

  # List high priority tasks
  task-cli list -priority high

  # Complete a task
  task-cli complete 1234567890

  # Delete a task
  task-cli delete 1234567890

  # Clear all tasks
  task-cli clear
`

	fmt.Println(usage)
}

// ============================================================================
// EXAMPLE SESSION
// ============================================================================

/*
$ task-cli add "Buy groceries" -desc "Milk, eggs, bread" -priority high
✅ Task added: Buy groceries (ID: 1642234567890123456)

$ task-cli add "Finish project" -priority medium
✅ Task added: Finish project (ID: 1642234567890123457)

$ task-cli list

📋 Tasks (2):

⏳ [1642234567890123456] 🔴 Buy groceries
   Milk, eggs, bread
   Created: 2024-01-15 10:30

⏳ [1642234567890123457] 🟡 Finish project
   Created: 2024-01-15 10:30

$ task-cli complete 1642234567890123456
✅ Task completed: Buy groceries

$ task-cli list -status pending

📋 Tasks (1):

⏳ [1642234567890123457] 🟡 Finish project
   Created: 2024-01-15 10:30

$ task-cli delete 1642234567890123457
🗑️  Task deleted: Finish project

$ task-cli list
📝 No tasks found
*/

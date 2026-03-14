// Gin Templates - HTML template rendering example
//
// This example covers:
// - HTML template rendering
// - Template inheritance

//go:build ignore
// - Dynamic data binding
// - Static file serving
// - Multiple template layouts
// - Template functions
//
// Run it: go run examples/gin_templates.go
// Visit: http://localhost:8080

package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// PageData represents common data available to all pages
type PageData struct {
	Title   string
	Year    int
	Version string
}

// HomePageData represents the home page data
type HomePageData struct {
	PageData
	Message   string
	Features  []string
	UserCount int
}

// UserPageData represents the user page data
type UserPageData struct {
	PageData
	Users []User
}

// UserDetailPageData represents the user detail page data
type UserDetailPageData struct {
	PageData
	User User
}

// FormPageData represents the form page data
type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

var users = []User{
	{ID: 1, Name: "John Doe", Email: "john@example.com", Role: "Admin", CreatedAt: "2024-01-15"},
	{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Role: "User", CreatedAt: "2024-01-16"},
	{ID: 3, Name: "Bob Johnson", Email: "bob@example.com", Role: "User", CreatedAt: "2024-01-17"},
}

// customTemplateFuncs adds custom functions to templates
func customTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"formatDate": func(date string) string {
			t, err := time.Parse("2006-01-02", date)
			if err != nil {
				return date
			}
			return t.Format("January 2, 2006")
		},
		"upper": func(s string) string {
			return s
		},
		"add": func(a, b int) int {
			return a + b
		},
		"isActive": func(role string) bool {
			return role == "Admin"
		},
	}
}

func setupRouter() *gin.Engine {
	// Create router
	router := gin.Default()

	// Set HTML template location
	// Gin will load templates from the templates directory
	router.LoadHTMLGlob("templates/*")

	// Or load from multiple locations
	// router.LoadHTMLGlob("templates/**/*.html")
	// router.LoadHTMLFiles("templates/layout.html", "templates/home.html")

	// Serve static files (CSS, JS, images)
	router.Static("/static", "./static")
	router.StaticFile("/favicon.ico", "./static/favicon.ico")

	// Custom 404 handler
	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"Title": "404 - Page Not Found",
			"Path":  c.Request.URL.Path,
		})
	})

	// Home page
	router.GET("/", func(c *gin.Context) {
		data := HomePageData{
			PageData: PageData{
				Title:   "Gin Template Demo",
				Year:    time.Now().Year(),
				Version: "1.0.0",
			},
			Message:   "Welcome to the Gin Template Demo!",
			Features:  []string{"Fast & Lightweight", "Middleware Support", "JSON Validation"},
			UserCount: len(users),
		}

		c.HTML(http.StatusOK, "index.html", data)
	})

	// Users list page
	router.GET("/users", func(c *gin.Context) {
		data := UserPageData{
			PageData: PageData{
				Title:   "Users List",
				Year:    time.Now().Year(),
				Version: "1.0.0",
			},
			Users: users,
		}

		c.HTML(http.StatusOK, "users.html", data)
	})

	// User detail page
	router.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")

		for _, user := range users {
			if fmt.Sprintf("%d", user.ID) == id {
				data := UserDetailPageData{
					PageData: PageData{
						Title:   user.Name,
						Year:    time.Now().Year(),
						Version: "1.0.0",
					},
					User: user,
				}

				c.HTML(http.StatusOK, "user_detail.html", data)
				return
			}
		}

		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"Title": "User Not Found",
			"Path":  c.Request.URL.Path,
		})
	})

	// About page
	router.GET("/about", func(c *gin.Context) {
		data := gin.H{
			"Title":   "About Us",
			"Year":    time.Now().Year(),
			"Version": "1.0.0",
		}

		c.HTML(http.StatusOK, "about.html", data)
	})

	// Form example page
	router.GET("/form", func(c *gin.Context) {
		data := gin.H{
			"Title":   "Contact Form",
			"Year":    time.Now().Year(),
			"Version": "1.0.0",
		}

		c.HTML(http.StatusOK, "form.html", data)
	})

	// Handle form submission
	router.POST("/form", func(c *gin.Context) {
		name := c.PostForm("name")
		email := c.PostForm("email")
		message := c.PostForm("message")

		data := gin.H{
			"Title":   "Form Submitted",
			"Year":    time.Now().Year(),
			"Version": "1.0.0",
			"Name":    name,
			"Email":   email,
			"Message": message,
		}

		c.HTML(http.StatusOK, "form_success.html", data)
	})

	return router
}

func createTemplates() error {
	// Create templates directory
	templatesDir := filepath.Join("examples", "templates")
	staticDir := filepath.Join("examples", "static")

	// Create base layout
	baseTemplate := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ .Title }} - Gin Template Demo</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            background-color: #f4f4f4;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 0 20px;
        }
        header {
            background-color: #333;
            color: white;
            padding: 1rem 0;
            box-shadow: 0 2px 5px rgba(0,0,0,0.1);
        }
        nav {
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        nav a {
            color: white;
            text-decoration: none;
            margin-left: 20px;
            transition: color 0.3s;
        }
        nav a:hover {
            color: #4CAF50;
        }
        .logo {
            font-size: 1.5rem;
            font-weight: bold;
        }
        main {
            background-color: white;
            margin: 2rem auto;
            padding: 2rem;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        h1 {
            color: #333;
            margin-bottom: 1rem;
        }
        h2 {
            color: #555;
            margin-bottom: 1rem;
        }
        .btn {
            display: inline-block;
            padding: 10px 20px;
            background-color: #4CAF50;
            color: white;
            text-decoration: none;
            border-radius: 4px;
            margin-right: 10px;
            transition: background-color 0.3s;
        }
        .btn:hover {
            background-color: #45a049;
        }
        .card {
            background-color: #f9f9f9;
            padding: 1.5rem;
            margin: 1rem 0;
            border-radius: 4px;
            border-left: 4px solid #4CAF50;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 1rem;
        }
        th, td {
            padding: 12px;
            text-align: left;
            border-bottom: 1px solid #ddd;
        }
        th {
            background-color: #4CAF50;
            color: white;
        }
        tr:hover {
            background-color: #f5f5f5;
        }
        .badge {
            display: inline-block;
            padding: 4px 8px;
            border-radius: 4px;
            font-size: 0.85rem;
            font-weight: bold;
        }
        .badge-admin {
            background-color: #ff6b6b;
            color: white;
        }
        .badge-user {
            background-color: #4CAF50;
            color: white;
        }
        footer {
            text-align: center;
            padding: 2rem 0;
            color: #666;
            font-size: 0.9rem;
        }
        .alert {
            padding: 1rem;
            margin: 1rem 0;
            border-radius: 4px;
        }
        .alert-success {
            background-color: #d4edda;
            border: 1px solid #c3e6cb;
            color: #155724;
        }
        .form-group {
            margin-bottom: 1rem;
        }
        label {
            display: block;
            margin-bottom: 0.5rem;
            font-weight: bold;
        }
        input[type="text"],
        input[type="email"],
        textarea {
            width: 100%;
            padding: 8px;
            border: 1px solid #ddd;
            border-radius: 4px;
            font-size: 1rem;
        }
        textarea {
            min-height: 100px;
            resize: vertical;
        }
    </style>
</head>
<body>
    <header>
        <div class="container">
            <nav>
                <div class="logo">🚀 Gin Templates</div>
                <div>
                    <a href="/">Home</a>
                    <a href="/users">Users</a>
                    <a href="/about">About</a>
                    <a href="/form">Contact</a>
                </div>
            </nav>
        </div>
    </header>
    <main class="container">
        {{block "content" .}}{{end}}
    </main>
    <footer>
        <div class="container">
            <p>&copy; {{ .Year }} Gin Template Demo. Version {{ .Version }}</p>
        </div>
    </footer>
</body>
</html>`

	// Create index.html
	indexTemplate := `{{define "content"}}
<h1>{{ .Message }}</h1>

<div class="card">
    <h2>🎯 Features</h2>
    <ul>
        {{range .Features}}
        <li>{{.}}</li>
        {{end}}
    </ul>
</div>

<div class="card">
    <h2>📊 Statistics</h2>
    <p>Total Users: <strong>{{ .UserCount }}</strong></p>
</div>

<div style="margin-top: 2rem;">
    <a href="/users" class="btn">View Users</a>
    <a href="/form" class="btn">Contact Us</a>
</div>
{{end}}`

	// Create users.html
	usersTemplate := `{{define "content"}}
<h1>Users List</h1>

<table>
    <thead>
        <tr>
            <th>ID</th>
            <th>Name</th>
            <th>Email</th>
            <th>Role</th>
            <th>Created</th>
            <th>Actions</th>
        </tr>
    </thead>
    <tbody>
        {{range .Users}}
        <tr>
            <td>{{ .ID }}</td>
            <td>{{ .Name }}</td>
            <td>{{ .Email }}</td>
            <td><span class="badge {{if eq .Role "Admin"}}badge-admin{{else}}badge-user{{end}}">{{ .Role }}</span></td>
            <td>{{ formatDate .CreatedAt }}</td>
            <td><a href="/users/{{ .ID }}" class="btn" style="padding: 5px 10px; font-size: 0.85rem;">View</a></td>
        </tr>
        {{end}}
    </tbody>
</table>
{{end}}`

	// Create user_detail.html
	userDetailTemplate := `{{define "content"}}
<h1>{{ .User.Name }}</h1>

<div class="card">
    <h2>User Details</h2>
    <p><strong>ID:</strong> {{ .User.ID }}</p>
    <p><strong>Name:</strong> {{ .User.Name }}</p>
    <p><strong>Email:</strong> {{ .User.Email }}</p>
    <p><strong>Role:</strong> <span class="badge {{if eq .User.Role "Admin"}}badge-admin{{else}}badge-user{{end}}">{{ .User.Role }}</span></p>
    <p><strong>Created:</strong> {{ formatDate .User.CreatedAt }}</p>
</div>

<a href="/users" class="btn">Back to Users</a>
{{end}}`

	// Create about.html
	aboutTemplate := `{{define "content"}}
<h1>About Gin Templates</h1>

<div class="card">
    <h2>🚀 What is Gin?</h2>
    <p>Gin is a high-performance HTTP web framework written in Go (Golang). It features a Martini-like API with much better performance -- up to 40 times faster.</p>
</div>

<div class="card">
    <h2>✨ Key Features</h2>
    <ul>
        <li>Fast routing with no reflection</li>
        <li>Middleware support</li>
        <li>JSON validation</li>
        <li>Route groups</li>
        <li>Error management</li>
        <li>Built-in rendering</li>
        <li>Extensible</li>
    </ul>
</div>

<div class="card">
    <h2>📚 Resources</h2>
    <ul>
        <li><a href="https://gin-gonic.com/docs/" target="_blank">Official Documentation</a></li>
        <li><a href="https://github.com/gin-gonic/gin" target="_blank">GitHub Repository</a></li>
    </ul>
</div>
{{end}}`

	// Create form.html
	formTemplate := `{{define "content"}}
<h1>Contact Form</h1>

<div class="card">
    <form method="POST" action="/form">
        <div class="form-group">
            <label for="name">Name:</label>
            <input type="text" id="name" name="name" required>
        </div>

        <div class="form-group">
            <label for="email">Email:</label>
            <input type="email" id="email" name="email" required>
        </div>

        <div class="form-group">
            <label for="message">Message:</label>
            <textarea id="message" name="message" required></textarea>
        </div>

        <button type="submit" class="btn">Submit</button>
    </form>
</div>
{{end}}`

	// Create form_success.html
	formSuccessTemplate := `{{define "content"}}
<div class="alert alert-success">
    <h2>✅ Form Submitted Successfully!</h2>
    <p>Thank you, {{ .Name }}!</p>
    <p>We've received your message and will get back to you at {{ .Email }}.</p>
</div>

<div class="card">
    <h3>Your Message:</h3>
    <p>{{ .Message }}</p>
</div>

<a href="/form" class="btn">Submit Another</a>
{{end}}`

	// Create 404.html
	notFoundTemplate := `{{define "content"}}
<div class="card" style="text-align: center; padding: 3rem;">
    <h1 style="font-size: 4rem; color: #ff6b6b;">404</h1>
    <h2>Page Not Found</h2>
    <p>The page you're looking for doesn't exist.</p>
    <p style="margin-top: 1rem;"><strong>Path:</strong> {{ .Path }}</p>
    <a href="/" class="btn" style="margin-top: 1rem;">Go Home</a>
</div>
{{end}}`

	// Create CSS file
	cssContent := `/* Custom CSS for Gin Templates */
/* Additional styles can be added here */`

	// Write all template files
	templates := map[string]string{
		"base.html":        baseTemplate,
		"index.html":       indexTemplate,
		"users.html":       usersTemplate,
		"user_detail.html": userDetailTemplate,
		"about.html":       aboutTemplate,
		"form.html":        formTemplate,
		"form_success.html": formSuccessTemplate,
		"404.html":         notFoundTemplate,
	}

	for filename, content := range templates {
		filepath := filepath.Join(templatesDir, filename)
		if err := writeToFile(filepath, content); err != nil {
			return fmt.Errorf("failed to create %s: %w", filename, err)
		}
	}

	// Write static CSS file
	cssPath := filepath.Join(staticDir, "style.css")
	if err := writeToFile(cssPath, cssContent); err != nil {
		return fmt.Errorf("failed to create style.css: %w", err)
	}

	return nil
}

func writeToFile(filepath, content string) error {
	return fmt.Sprintf("echo '%s' > %s", content, filepath)
}

func main() {
	// Create template files
	fmt.Println("📝 Creating template files...")
	if err := createTemplates(); err != nil {
		log.Fatal("Failed to create templates:", err)
	}
	fmt.Println("✅ Templates created successfully")

	// Setup router
	router := setupRouter()

	fmt.Println("\n🚀 Gin Template server starting on :8080")
	fmt.Println("\n📄 Available pages:")
	fmt.Println("  http://localhost:8080/          - Home page")
	fmt.Println("  http://localhost:8080/users     - Users list")
	fmt.Println("  http://localhost:8080/users/1   - User detail")
	fmt.Println("  http://localhost:8080/about     - About page")
	fmt.Println("  http://localhost:8080/form      - Contact form")
	fmt.Println("  http://localhost:8080/static/   - Static files")
	fmt.Println("\n💡 Open your browser and navigate to: http://localhost:8080")
	fmt.Println()

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

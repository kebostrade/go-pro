package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Shared data models
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// In-memory store (shared across all API types)
type Store struct {
	mu     sync.RWMutex
	users  map[int]*User
	nextID int
}

func NewStore() *Store {
	return &Store{
		users:  make(map[int]*User),
		nextID: 1,
	}
}

func (s *Store) Create(username, email, role string) *User {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	user := &User{
		ID:        s.nextID,
		Username:  username,
		Email:     email,
		Role:      role,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.users[s.nextID] = user
	s.nextID++
	return user
}

func (s *Store) GetAll() []*User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]*User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users
}

func (s *Store) GetByID(id int) (*User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[id]
	return user, exists
}

// API Gateway
type Gateway struct {
	store      *Store
	restRouter *gin.Engine
	grpcClient interface{} // Would be actual gRPC client
}

func NewGateway() *Gateway {
	gin.SetMode(gin.ReleaseMode)
	
	g := &Gateway{
		store:      NewStore(),
		restRouter: gin.Default(),
	}
	
	g.setupRESTRoutes()
	return g
}

// REST API Routes
func (g *Gateway) setupRESTRoutes() {
	api := g.restRouter.Group("/api/rest")
	{
		api.GET("/users", g.handleRESTGetUsers)
		api.POST("/users", g.handleRESTCreateUser)
		api.GET("/users/:id", g.handleRESTGetUser)
	}
}

func (g *Gateway) handleRESTGetUsers(c *gin.Context) {
	users := g.store.GetAll()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    users,
		"source":  "REST API",
	})
}

func (g *Gateway) handleRESTCreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Role     string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := g.store.Create(req.Username, req.Email, req.Role)
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    user,
		"source":  "REST API",
	})
}

func (g *Gateway) handleRESTGetUser(c *gin.Context) {
	id := 1 // Simplified
	user, exists := g.store.GetByID(id)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
		"source":  "REST API",
	})
}

// GraphQL Handler
func (g *Gateway) handleGraphQL(c *gin.Context) {
	var req struct {
		Query string `json:"query"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Simplified GraphQL query handling
	if contains(req.Query, "users") {
		users := g.store.GetAll()
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"users": users,
			},
			"source": "GraphQL",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   nil,
		"source": "GraphQL",
	})
}

// gRPC Handler (simplified - would normally use actual gRPC server)
func (g *Gateway) handleGRPCProxy(c *gin.Context) {
	// In a real implementation, this would proxy to a gRPC server
	// For this tutorial, we'll simulate it
	
	users := g.store.GetAll()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    users,
		"source":  "gRPC (proxied)",
		"note":    "This is a simplified proxy. In production, use actual gRPC client.",
	})
}

// Helper function
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0
}

// Dashboard handler
func (g *Gateway) handleDashboard(c *gin.Context) {
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, dashboardHTML)
}

// API comparison endpoint
func (g *Gateway) handleComparison(c *gin.Context) {
	comparison := map[string]interface{}{
		"rest": map[string]interface{}{
			"pros": []string{
				"Simple and widely understood",
				"Stateless and cacheable",
				"Good tooling support",
				"Easy to debug",
			},
			"cons": []string{
				"Over-fetching or under-fetching data",
				"Multiple round trips for related data",
				"Versioning challenges",
			},
			"use_cases": []string{
				"Public APIs",
				"CRUD operations",
				"Mobile apps",
				"Web applications",
			},
		},
		"grpc": map[string]interface{}{
			"pros": []string{
				"High performance with Protocol Buffers",
				"Bi-directional streaming",
				"Strong typing",
				"Code generation",
			},
			"cons": []string{
				"Steeper learning curve",
				"Limited browser support",
				"Binary format (harder to debug)",
			},
			"use_cases": []string{
				"Microservices communication",
				"Real-time streaming",
				"Internal APIs",
				"High-performance systems",
			},
		},
		"graphql": map[string]interface{}{
			"pros": []string{
				"Flexible queries",
				"Single endpoint",
				"No over-fetching",
				"Strong typing with schema",
			},
			"cons": []string{
				"Complexity in implementation",
				"Caching challenges",
				"Query complexity management",
			},
			"use_cases": []string{
				"Complex data requirements",
				"Mobile apps with limited bandwidth",
				"Aggregating multiple data sources",
				"Rapid frontend development",
			},
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"comparison": comparison,
	})
}

func main() {
	gateway := NewGateway()

	// Seed data
	gateway.store.Create("alice", "alice@example.com", "admin")
	gateway.store.Create("bob", "bob@example.com", "user")
	gateway.store.Create("charlie", "charlie@example.com", "guest")

	// Setup routes
	gateway.restRouter.GET("/", gateway.handleDashboard)
	gateway.restRouter.POST("/api/graphql", gateway.handleGraphQL)
	gateway.restRouter.GET("/api/grpc/users", gateway.handleGRPCProxy)
	gateway.restRouter.GET("/api/comparison", gateway.handleComparison)

	port := ":8082"
	fmt.Println("🚀 API Gateway starting on http://localhost" + port)
	fmt.Println("=" + string(make([]byte, 60)))
	fmt.Println("\n📚 Available endpoints:")
	fmt.Println("\n🌐 Dashboard:")
	fmt.Println("  GET  /                        - Interactive dashboard")
	fmt.Println("\n🔵 REST API:")
	fmt.Println("  GET  /api/rest/users          - List users")
	fmt.Println("  POST /api/rest/users          - Create user")
	fmt.Println("  GET  /api/rest/users/:id      - Get user")
	fmt.Println("\n🟢 GraphQL:")
	fmt.Println("  POST /api/graphql             - GraphQL endpoint")
	fmt.Println("\n🟣 gRPC (proxied):")
	fmt.Println("  GET  /api/grpc/users          - List users via gRPC")
	fmt.Println("\n📊 Comparison:")
	fmt.Println("  GET  /api/comparison          - API comparison guide")
	fmt.Println("\n" + "=" + string(make([]byte, 60)))

	if err := gateway.restRouter.Run(port); err != nil {
		log.Fatal(err)
	}
}

const dashboardHTML = `<!DOCTYPE html>
<html>
<head>
    <title>API Gateway Dashboard</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background: #f5f5f5; }
        .header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 30px; text-align: center; }
        .container { max-width: 1200px; margin: 30px auto; padding: 0 20px; }
        .grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(350px, 1fr)); gap: 20px; margin-top: 30px; }
        .card { background: white; border-radius: 10px; padding: 25px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .card h2 { color: #333; margin-bottom: 15px; font-size: 1.5em; }
        .endpoint { background: #f8f9fa; padding: 12px; margin: 10px 0; border-radius: 5px; font-family: monospace; font-size: 0.9em; }
        .method { display: inline-block; padding: 4px 8px; border-radius: 3px; font-weight: bold; margin-right: 10px; }
        .get { background: #61affe; color: white; }
        .post { background: #49cc90; color: white; }
        button { background: #667eea; color: white; border: none; padding: 12px 24px; border-radius: 5px; cursor: pointer; font-size: 1em; margin: 5px; }
        button:hover { background: #5568d3; }
        #response { background: #f8f9fa; padding: 15px; border-radius: 5px; margin-top: 15px; white-space: pre-wrap; font-family: monospace; font-size: 0.9em; max-height: 400px; overflow-y: auto; }
    </style>
</head>
<body>
    <div class="header">
        <h1>🚀 API Gateway Dashboard</h1>
        <p>REST • gRPC • GraphQL - All in One Place</p>
    </div>
    
    <div class="container">
        <div class="grid">
            <div class="card">
                <h2>🔵 REST API</h2>
                <div class="endpoint"><span class="method get">GET</span>/api/rest/users</div>
                <div class="endpoint"><span class="method post">POST</span>/api/rest/users</div>
                <button onclick="testREST()">Test REST API</button>
            </div>
            
            <div class="card">
                <h2>🟢 GraphQL</h2>
                <div class="endpoint"><span class="method post">POST</span>/api/graphql</div>
                <p style="margin: 10px 0; color: #666;">Flexible queries with a single endpoint</p>
                <button onclick="testGraphQL()">Test GraphQL</button>
            </div>
            
            <div class="card">
                <h2>🟣 gRPC (Proxied)</h2>
                <div class="endpoint"><span class="method get">GET</span>/api/grpc/users</div>
                <p style="margin: 10px 0; color: #666;">High-performance RPC framework</p>
                <button onclick="testGRPC()">Test gRPC</button>
            </div>
        </div>
        
        <div class="card" style="margin-top: 20px;">
            <h2>📊 Response</h2>
            <div id="response">Click a button above to test an API...</div>
        </div>
        
        <div class="card" style="margin-top: 20px;">
            <h2>📚 API Comparison</h2>
            <button onclick="showComparison()">Show Comparison</button>
        </div>
    </div>
    
    <script>
        async function testREST() {
            const response = await fetch('/api/rest/users');
            const data = await response.json();
            document.getElementById('response').textContent = JSON.stringify(data, null, 2);
        }
        
        async function testGraphQL() {
            const response = await fetch('/api/graphql', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ query: '{ users { id username email } }' })
            });
            const data = await response.json();
            document.getElementById('response').textContent = JSON.stringify(data, null, 2);
        }
        
        async function testGRPC() {
            const response = await fetch('/api/grpc/users');
            const data = await response.json();
            document.getElementById('response').textContent = JSON.stringify(data, null, 2);
        }
        
        async function showComparison() {
            const response = await fetch('/api/comparison');
            const data = await response.json();
            document.getElementById('response').textContent = JSON.stringify(data, null, 2);
        }
    </script>
</body>
</html>`


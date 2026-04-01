package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Home renders the home page
func Home(c *gin.Context) {
	data := struct {
		Title string
		Year  int
	}{
		Title: "Gin Web App",
		Year:  time.Now().Year(),
	}
	c.HTML(http.StatusOK, "home.html", data)
}

// About renders the about page
func About(c *gin.Context) {
	data := struct {
		Title string
		Year  int
	}{
		Title: "About - Gin Web App",
		Year:  time.Now().Year(),
	}
	c.HTML(http.StatusOK, "about.html", data)
}

// HealthCheck returns API health status
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"version": "1.0.0",
	})
}

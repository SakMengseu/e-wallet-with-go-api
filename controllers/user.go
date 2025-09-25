package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	// Static test data
	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully!",
		"user": map[string]string{
			"name":  "Test User",
			"email": "test@example.com",
		},
	})
}

func GetUsers(c *gin.Context) {
	// Static list
	c.JSON(http.StatusOK, gin.H{
		"users": []map[string]string{
			{"name": "Alice", "email": "alice@example.com"},
			{"name": "Bob", "email": "bob@example.com"},
		},
	})
}

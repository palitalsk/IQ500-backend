package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from header
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID required"})
			c.Abort()
			return
		}

		// Set in context
		c.Set("user_id", userID)
		c.Next()
	}
}

func BusinessAccessMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = c.GetString("user_id")
		businessCode := c.Param("business_code")

		if businessCode == "" {
			// query params
			businessCode = c.Query("business_code")
		}

		if businessCode == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Business Code required"})
			c.Abort()
			return
		}

		c.Set("business_code", businessCode)
		c.Next()
	}
}

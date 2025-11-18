package Infrastructure

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Auth middleware expects Authorization: Bearer <token>
// places (username, role, user_id) into gin.Context keys
func AuthMiddleware(jwtSrv JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			return
		}
		parts := strings.SplitN(h, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid auth header"})
			return
		}
		token := parts[1]
		_, claims, err := jwtSrv.Validate(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token", "detail": err.Error()})
			return
		}
		// expected claims: username, role, sub
		if uname, ok := claims["username"].(string); ok {
			c.Set("username", uname)
		}
		if role, ok := claims["role"].(string); ok {
			c.Set("role", role)
		}
		if sub, ok := claims["sub"].(string); ok {
			c.Set("user_id", sub)
		}
		c.Next()
	}
}

// Admin middleware
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		rv, ok := c.Get("role")
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "role missing"})
			return
		}
		role, ok := rv.(string)
		if !ok || role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin role required"})
			return
		}
		c.Next()
	}
}

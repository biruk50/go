package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"task_manager_testify/Infrastructure"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddlewareValidAndInvalid(t *testing.T) {
	gin.SetMode(gin.TestMode)
	j := Infrastructure.NewJWTService("secret-test")
	token, _ := j.Generate("uid", "alice", "admin", time.Hour)

	r := gin.New()
	r.GET("/protected", Infrastructure.AuthMiddleware(j), func(c *gin.Context) {
		user := c.GetString("username")
		c.JSON(200, gin.H{"user": user})
	})

	// valid token
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// missing token
	req2 := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, 401, w2.Code)
}

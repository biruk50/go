package infrastructure

import (
	"testing"
	"time"
	"task_manager_testify/Infrastructure"
	"github.com/stretchr/testify/assert"
) 

func TestJWTGenerateAndValidate(t *testing.T) {
	secret := "test-secret-123"
	j := Infrastructure.NewJWTService(secret)

	token, err := j.Generate("uid123", "alice", "admin", time.Hour)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validate returns claims map in our implementation
	_, claims, err := j.Validate(token)
	assert.NoError(t, err)
	assert.Equal(t, "alice", claims["username"])
	assert.Equal(t, "admin", claims["role"])
	assert.Equal(t, "uid123", claims["sub"])
}

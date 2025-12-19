package infrastructure

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"task_manager_testify/Infrastructure"
)

func TestPasswordHashAndCompare(t *testing.T) {
	pw := Infrastructure.NewPasswordService()

	hash, err := pw.Hash("s3cret")
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	// correct password
	err = pw.Compare(hash, "s3cret")
	assert.NoError(t, err)

	// wrong password
	err = pw.Compare(hash, "wrong")
	assert.Error(t, err)
}

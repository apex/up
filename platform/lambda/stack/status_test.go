package stack

import (
	"testing"

	"github.com/tj/assert"
)

func TestStatus_String(t *testing.T) {
	assert.Equal(t, "Creating", Status("CREATE_IN_PROGRESS").String())
	assert.Equal(t, "Deleting", Status("DELETE_IN_PROGRESS").String())
	assert.Equal(t, "Failed to update", Status("UPDATE_FAILED").String())
}

func TestStatus_State(t *testing.T) {
	assert.Equal(t, Pending, Status("CREATE_IN_PROGRESS").State())
	assert.Equal(t, Pending, Status("UPDATE_IN_PROGRESS").State())
	assert.Equal(t, Success, Status("CREATE_COMPLETE").State())
	assert.Equal(t, Failure, Status("CREATE_FAILED").State())
}

func TestStatus_IsDone(t *testing.T) {
	assert.False(t, Status("CREATE_IN_PROGRESS").IsDone())
	assert.False(t, Status("UPDATE_IN_PROGRESS").IsDone())
	assert.True(t, Status("CREATE_COMPLETE").IsDone())
	assert.True(t, Status("UPDATE_COMPLETE").IsDone())
	assert.True(t, Status("DELETE_COMPLETE").IsDone())
	assert.True(t, Status("DELETE_FAILED").IsDone())
}

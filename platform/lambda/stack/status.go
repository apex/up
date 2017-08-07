package stack

import (
	"fmt"

	"github.com/apex/up/internal/colors"
)

// status map for humanization.
var statusMap = map[Status]string{
	CreateInProgress: "Creating",
	CreateFailed:     "Failed to create",
	CreateComplete:   "Created",

	DeleteInProgress: "Deleting",
	DeleteFailed:     "Failed to delete",
	DeleteComplete:   "Deleted",
	DeleteSkipped:    "Skipped",

	UpdateInProgress: "Updating",
	UpdateFailed:     "Failed to update",
	UpdateComplete:   "Updated",

	RollbackInProgress: "Rolling back",
	RollbackFailed:     "Failed to rollback",
	RollbackComplete:   "Rollback complete",
}

// State represents a generalized stack event state.
type State int

// States available.
const (
	Success State = iota
	Pending
	Failure
)

// Status represents a stack event status.
type Status string

// Statuses available.
const (
	Unknown Status = "unknown"

	CreateInProgress = "CREATE_IN_PROGRESS"
	CreateFailed     = "CREATE_FAILED"
	CreateComplete   = "CREATE_COMPLETE"

	DeleteInProgress = "DELETE_IN_PROGRESS"
	DeleteFailed     = "DELETE_FAILED"
	DeleteComplete   = "DELETE_COMPLETE"
	DeleteSkipped    = "DELETE_SKIPPED"

	UpdateInProgress = "UPDATE_IN_PROGRESS"
	UpdateFailed     = "UPDATE_FAILED"
	UpdateComplete   = "UPDATE_COMPLETE"

	RollbackInProgress = "ROLLBACK_IN_PROGRESS"
	RollbackFailed     = "ROLLBACK_FAILED"
	RollbackComplete   = "ROLLBACK_COMPLETE"
)

// String returns the human representation.
func (s Status) String() string {
	return statusMap[s]
}

// IsDone returns true when failed or complete.
func (s Status) IsDone() bool {
	return s.State() != Pending
}

// Color the given string based on the status.
func (s Status) Color(v string) string {
	switch s.State() {
	case Success:
		return colors.Blue(v)
	case Pending:
		return colors.Yellow(v)
	case Failure:
		return colors.Red(v)
	default:
		return v
	}
}

// State returns a generalized state.
func (s Status) State() State {
	switch s {
	case CreateFailed, UpdateFailed, DeleteFailed, RollbackFailed:
		return Failure
	case CreateInProgress, UpdateInProgress, DeleteInProgress, RollbackInProgress:
		return Pending
	case CreateComplete, UpdateComplete, DeleteComplete, DeleteSkipped, RollbackComplete:
		return Success
	default:
		panic(fmt.Sprintf("unhandled state %q", s))
	}
}

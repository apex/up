package stack

import (
	"fmt"

	"github.com/apex/up/internal/colors"
)

// status map for humanization.
var statusMap = map[Status]string{
	Unknown: "Unknown",

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

	UpdateCompleteCleanup:         "Update complete cleanup in progress",
	UpdateRollbackCompleteCleanup: "Update rollback complete cleanup in progress",
	UpdateRollbackInProgress:      "Update rollback in progress",
	UpdateRollbackComplete:        "Update rollback complete",

	RollbackInProgress: "Rolling back",
	RollbackFailed:     "Failed to rollback",
	RollbackComplete:   "Rollback complete",

	CreatePending: "Create pending",
	Failed:        "Failed",
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
	Unknown Status = ""

	CreateInProgress = "CREATE_IN_PROGRESS"
	CreateFailed     = "CREATE_FAILED"
	CreateComplete   = "CREATE_COMPLETE"
	CreatePending    = "CREATE_PENDING"

	DeleteInProgress = "DELETE_IN_PROGRESS"
	DeleteFailed     = "DELETE_FAILED"
	DeleteComplete   = "DELETE_COMPLETE"
	DeleteSkipped    = "DELETE_SKIPPED"

	UpdateInProgress = "UPDATE_IN_PROGRESS"
	UpdateFailed     = "UPDATE_FAILED"
	UpdateComplete   = "UPDATE_COMPLETE"

	UpdateRollbackInProgress      = "UPDATE_ROLLBACK_IN_PROGRESS"
	UpdateRollbackComplete        = "UPDATE_ROLLBACK_COMPLETE"
	UpdateRollbackCompleteCleanup = "UPDATE_ROLLBACK_COMPLETE_CLEANUP_IN_PROGRESS"
	UpdateCompleteCleanup         = "UPDATE_COMPLETE_CLEANUP_IN_PROGRESS"

	RollbackInProgress = "ROLLBACK_IN_PROGRESS"
	RollbackFailed     = "ROLLBACK_FAILED"
	RollbackComplete   = "ROLLBACK_COMPLETE"

	Failed = "FAILED"
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
	case CreateFailed, UpdateFailed, DeleteFailed, RollbackFailed, Failed:
		return Failure
	case CreateInProgress, UpdateInProgress, DeleteInProgress, RollbackInProgress, CreatePending, UpdateRollbackInProgress:
		return Pending
	case CreateComplete, UpdateComplete, DeleteComplete, DeleteSkipped, RollbackComplete, UpdateCompleteCleanup, UpdateRollbackCompleteCleanup, UpdateRollbackComplete:
		return Success
	default:
		panic(fmt.Sprintf("unhandled state %q", string(s)))
	}
}

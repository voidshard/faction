package structs

// Event is something we want to report to the caller
type Event struct {
	// Key is some user defined key in an Action config (see pkg/config/action.go)
	Key string

	// TODO: expand on event
	Message string
}

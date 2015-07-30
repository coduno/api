package model

// Endpoints encapsulates two possible
// ways to deliver the outcome of trying
// to fulfill an assignment.
type Endpoints struct {
	// Name of the WebComponent used to render
	// the assignment accordingly.
	WebInterface string

	// URL of the remote that should be pushed
	// to.
	//
	// NOTE(flowlo): No backwards-compatibility
	// guarantee on this.
	//
	// TODO(flowlo): Investigate why we can not use
	// url.URL
	GitRepository string
}

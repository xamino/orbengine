package orbengine

import "errors"

/*
Engine errors.
*/
var (
	ErrMissingComponents = errors.New("missing required components")
	ErrEntityIDExists    = errors.New("entity ID already in use")
)

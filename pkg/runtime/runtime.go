package runtime

import "github.com/kataras/muxie"

// Controller constraints all concrete types provide a function,
// that can process routes register.
type Controller interface {
	// RegisterRoute sets route rules for the concrete controller.
	RegisterRoute(muxie.SubMux)
}

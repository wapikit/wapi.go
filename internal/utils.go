package internal

import (
	"github.com/go-playground/validator/v10"
)

// validate is shared: a *validator.Validate is safe for concurrent use and
// caches per-struct reflection metadata, so building one per call is wasteful.
var validate = validator.New()

// GetValidator returns the shared validator instance.
func GetValidator() *validator.Validate {
	return validate
}

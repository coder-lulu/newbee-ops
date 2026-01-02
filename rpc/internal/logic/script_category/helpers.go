package script_category

import "fmt"

// newInvalidArgumentError creates a simple error for invalid arguments
func newInvalidArgumentError(field string) error {
	return fmt.Errorf("invalid argument: %s is required", field)
}

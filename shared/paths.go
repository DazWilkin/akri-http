package shared

import (
	"strings"
)

// Paths is an alias to use repeated flags with flag
type Paths []string

// String is a method required by flag.Value interface
func (e *Paths) String() string {
	result := strings.Join(*e, "\n")
	return result
}

// Set is a method required by flag.Value interface
func (e *Paths) Set(value string) error {
	*e = append(*e, value)
	return nil
}

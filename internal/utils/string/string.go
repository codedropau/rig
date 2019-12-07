package string

import (
	"fmt"
	"strings"
)

// Contains a string in the slice.
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}

// SplitBySeparator is used for marshalling Docker Compose configuration.
func SplitBySeparator(value, separator string) (string, string, error) {
	sl := strings.SplitN(value, separator, 2)

	if len(sl) != 2 {
		return "", "", fmt.Errorf("unable to split: %s", value)
	}

	return sl[0], sl[1], nil
}

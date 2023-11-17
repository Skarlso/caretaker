package slash

import (
	"fmt"
	"strings"
)

const keyValueLength = 2

// ConvertArgs converts a list of arguments into a key=value map.
func ConvertArgs(args ...string) (map[string]string, error) {
	result := make(map[string]string)

	for _, arg := range args {
		split := strings.Split(arg, "=")
		if len(split) != keyValueLength {
			return nil, fmt.Errorf("invalid format for argument, wanted k=v but was: %s", arg)
		}

		result[split[0]] = split[1]
	}

	return result, nil
}

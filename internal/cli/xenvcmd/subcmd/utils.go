package subcmd

import "fmt"

// parseNameVersion parses a string in the format name:version into its components
func parseNameVersion(input string) (name, version string, err error) {
	parts := splitByLast(input, ":")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid format, expected name:version, got %s", input)
	}
	return parts[0], parts[1], nil
}

// splitByLast splits a string by the last occurrence of the separator
func splitByLast(s, sep string) []string {
	lastIndex := -1
	for i := len(s) - 1; i >= 0; i-- {
		if string(s[i]) == sep {
			lastIndex = i
			break
		}
	}

	if lastIndex == -1 {
		return []string{s}
	}

	return []string{s[:lastIndex], s[lastIndex+1:]}
}
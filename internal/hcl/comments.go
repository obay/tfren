package hcl

import "strings"

func findPrecedingComments(lines []string, blockStartLine int) []string {
	var comments []string

	for i := blockStartLine - 2; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])

		if strings.HasPrefix(line, "#") ||
			strings.HasPrefix(line, "//") ||
			strings.HasPrefix(line, "/*") ||
			strings.HasSuffix(line, "*/") {
			comments = append([]string{lines[i]}, comments...)
		} else if line == "" {
			break
		} else {
			break
		}
	}

	return comments
}

func extractLines(lines []string, start, end int) string {
	if start < 0 {
		start = 0
	}
	if end > len(lines) {
		end = len(lines)
	}

	var result []string
	for i := start; i < end; i++ {
		result = append(result, lines[i])
	}
	return strings.Join(result, "\n")
}

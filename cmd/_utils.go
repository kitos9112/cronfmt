package cmd

import (
	"regexp"
)

func isSubstring(str string, substr string) bool {
	re := regexp.MustCompile(substr)
	if re.MatchString(str) {
		return true
	}
	return false
}

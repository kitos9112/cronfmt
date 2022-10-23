package cmd

import (
	"fmt"
	"regexp"
	"strings"
)

// isSubstring returns true whether a substring is within a string, otherwise false
func isSubstring(str string, substr string) bool {
	re := regexp.MustCompile(substr)
	if re.MatchString(str) {
		return true
	}
	return false
}

// getMatchesLength returns the number of subStr matches found in a string
func getMatchesLength(str string, substr string) int {
	re := regexp.MustCompile(substr)
	return len(re.FindAllStringIndex(str, -1))
}

// createStringRange returns a string range from start to end with a step to calculate offset.
func createStringRange(minRange int, maxRange int, step int) string {
	var sliceRange []int

	if minRange == 0 && (maxRange+1)%step == 0 {
		minRange = 1
		sliceRange = append(sliceRange, 0)
	} else {
		minRange = 1
	}

	for i := minRange; i <= maxRange; i++ {
		if i%step == 0 {
			sliceRange = append(sliceRange, i)
		}
	}
	sliceStr := strings.Fields(fmt.Sprint(sliceRange))
	sliceStrSpaces := strings.Join(sliceStr, " ")
	return strings.Trim(sliceStrSpaces, "[]")
}

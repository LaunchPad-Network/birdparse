package birdparse

import (
	"regexp"
	"strconv"
	"strings"
)

func parseOptionalInt(s string) int {
	if strings.Contains(s, "-") && !regexp.MustCompile(`^\d+$`).MatchString(s) {
		return 0
	}
	return atoi(s)
}

func atoi(s string) int {
	i, _ := strconv.Atoi(strings.Split(s, ".")[0])
	return i
}

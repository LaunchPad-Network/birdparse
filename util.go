package birdparse

import (
	"regexp"
	"strconv"
	"strings"
)

func parseOptionalIntAsString(s string) string {
	if strings.Contains(s, "-") && !regexp.MustCompile(`^\d+$`).MatchString(s) {
		return "0"
	}
	return strconv.Itoa(atoi(s))
}

func atoi(s string) int {
	i, _ := strconv.Atoi(strings.Split(s, ".")[0])
	return i
}

package util

import (
	"regexp"
)

var ansiEscapeCodes = `\x1B\[[0-?]*[ -/]*[@-~]`
var ansiRegex = regexp.MustCompile(ansiEscapeCodes)

func RemoveAnsiEscapeCodes(input string) string {
	return ansiRegex.ReplaceAllString(input, "")
}

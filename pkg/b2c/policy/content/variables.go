package content

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	variableRegEx         = regexp.MustCompile("{Settings:(.+)}")
	variableReplaceFormat = "{Settings:%s}"
)

// SetVariableIdentifierRegexFromString creates a regex from the given string and updates the
// - variableRegEx which is used to find and extract variable names from content
// - variableReplaceFormat which is used to replace variable identifiers with a concrete value
//
// Example:
//
//	SetVariableIdentifierRegexFromString("{Settings:(.+)}")
func SetVariableIdentifierRegexFromString(s string) {
	variableRegEx = regexp.MustCompile(s)
	SetVariableIdentifierRegex(variableRegEx)
}

// SetVariableIdentifierRegex uses the given regex to update the
// - variableRegEx which is used to find and extract variable names from content
// - variableReplaceFormat which is used to replace variable identifiers with a concrete value
func SetVariableIdentifierRegex(r *regexp.Regexp) {
	variableRegEx = r
	setVariableReplaceFormatString()
}

// setVariableReplaceFormatString extracts delimiters and variable identifier from the variableRegEx
// and builds a format string that can be used to replace variable identifiers with a concrete value
//
// Example: "{Settings:(.+)}" will be transferred to "{Setting:%s}"
func setVariableReplaceFormatString() {
	s := variableRegEx.String()
	p := strings.Split(s[1:len(s)-1], ":")
	variableReplaceFormat = fmt.Sprintf("%s%s:%s%s", string(s[0]), p[0], "%s", s[len(s)-1:])
}

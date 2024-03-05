package content

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/go-multierror"
)

// Raw is a type alias for a byte slice used to represent raw policy content.
type Raw []byte

// String returns the string representation of the Raw.
// Example usage:
//
//	c := Raw([]byte{'h', 'e', 'l', 'l', 'o'})
//	sourceRoot := c.String()  // sourceRoot = "hello"
func (c *Raw) String() string {
	return string(*c)
}

// ReplaceVariables replaces variables in the Raw with provided values.
// It takes a map of variable-value pairs as input.
// If a variable is not found in the map, it returns an error.
// If the provided value for a variable is empty or equal to "null", it returns an error.
// This method uses regular expressions to find and replace variables in the Raw.
// It modifies the Raw in-place.
//
// Example usage:
//
//	c := Raw("{Settings:foo}")
//	vs := map[string]string{"foo": "bar"}
//	err := c.ReplaceVariables(vs)
//	if err != nil {
//	    fmt.Println(err)
//	}
//	sourceRoot := c.String()  // sourceRoot = "bar"
func (c *Raw) ReplaceVariables(vs map[string]string, f string) error {
	var err error
	for _, v := range c.ExtractVariables() {
		val, ok := vs[v]
		if !ok {
			err = multierror.Append(err, fmt.Errorf("variable '%s' requested in %s but not provided", v, f))
			continue
		}
		if val == "" || val == "null" {
			err = multierror.Append(err, fmt.Errorf("provided value for '%s' must not be empty", v))
			continue
		}
		re := regexp.MustCompile(fmt.Sprintf(variableReplaceFormat, v))
		*c = []byte(re.ReplaceAllString(c.String(), val))
	}

	if err != nil {
		return err
	}

	return nil
}

// ExtractVariables extracts unique variables from the Raw.
//
// It uses a regular expression to find all matches of variables in the Raw.
// The found variables are stored in a slice without duplication.
//
// Returns:
// - cm: a slice containing all unique variables found in the Raw
func (c *Raw) ExtractVariables() []string {
	m := variableRegEx.FindAllStringSubmatch(c.String(), -1)
	var cm []string
	seen := make(map[string]bool, len(m))
	for _, match := range m {
		if !seen[match[1]] {
			cm = append(cm, match[1])
			seen[match[1]] = true
		}
	}

	return cm
}

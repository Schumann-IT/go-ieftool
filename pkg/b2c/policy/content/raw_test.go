package content

import (
	"errors"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
)

func Test_ExtractVariables(t *testing.T) {
	expected := []string{
		"Name",
	}

	c := Raw("{Settings:Name}")
	actual := c.ExtractVariables()

	assert.Equal(t, expected, actual)
}

func Test_ReplaceVariables(t *testing.T) {
	expected := "Var1\nVar2"

	c := Raw("{Settings:Var1}\n{Settings:Var2}")
	_ = c.ReplaceVariables(map[string]string{
		"Var1": "Var1",
		"Var2": "Var2",
	}, "virtual")
	actual := c.String()

	assert.Equal(t, expected, actual)
}

func Test_ReplaceVariablesFailsForNonExistentVariable(t *testing.T) {
	expected := 2

	c := Raw("{Settings:NotExistent}\n{Settings:NotExistent2}")
	err := c.ReplaceVariables(map[string]string{}, "virtual")
	var mer *multierror.Error
	errors.As(err, &mer)
	actual := mer.Len()

	assert.Equal(t, expected, actual)
}

func Test_ReplaceVariablesFailsForEmptyVariable(t *testing.T) {
	expected := 2

	c := Raw("{Settings:Empty}\n{Settings:Null}")
	err := c.ReplaceVariables(map[string]string{
		"Empty": "",
		"Null":  "null",
	}, "virtual")
	var mer *multierror.Error
	errors.As(err, &mer)
	actual := mer.Len()

	assert.Equal(t, expected, actual)
}

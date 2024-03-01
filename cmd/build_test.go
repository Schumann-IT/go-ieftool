package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ExecuteBuild(t *testing.T) {
	actual := new(bytes.Buffer)
	cmd := rootCmd
	cmd.SetOut(actual)
	cmd.SetArgs([]string{"build", "--config", "../test/fixtures/config.yaml", "--source", "../test/fixtures/base", "--destination", "../build"})
	cmd.Execute()

	assert.Equal(t, "", actual.String(), "error is expected")
}

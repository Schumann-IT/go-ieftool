package content

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CustomVariableRegex(t *testing.T) {
	assert.Equal(t, "{Settings:(.+)}", variableRegEx.String())
	assert.Equal(t, "{Settings:%s}", variableReplaceFormat)

	rs := "|Variable:(.*)|"
	SetVariableIdentifierRegexFromString(rs)

	assert.Equal(t, rs, variableRegEx.String())
	assert.Equal(t, "|Variable:%s|", variableReplaceFormat)
}

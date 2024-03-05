package content

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SourceMapLength(t *testing.T) {
	expected := 0
	m := Source{}
	actual := m.Len()
	assert.Equal(t, expected, actual)

	expected = 1
	m = Source{
		"e": Raw("example"),
	}
	actual = m.Len()
	assert.Equal(t, expected, actual)
}

func Test_ProcessedMapLength(t *testing.T) {
	expected := 0
	m := Processed{}
	actual := m.Len()
	assert.Equal(t, expected, actual)

	expected = 1
	m = Processed{
		"e": []byte("example"),
	}
	actual = m.Len()
	assert.Equal(t, expected, actual)
}

package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Branch(t *testing.T) {
	b := NewBranch("foo")
	assert.Equal(t, "foo", b.Data(), "root data does not match")

	b.AddChild(NewBranch("bar"))
	assert.Equal(t, 1, len(b.Children()), "child count for root does not match")

	for _, c := range b.Children() {
		assert.Equal(t, 0, len(c.Children()), "child count for child does not match")
		assert.Equal(t, "bar", c.Data(), "child data does not match")
	}
}

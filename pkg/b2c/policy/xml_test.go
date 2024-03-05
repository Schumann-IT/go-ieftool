package policy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ReadRootXmlPolicy(t *testing.T) {
	p, err := New("../../../test/fixtures/build/base/base.xml")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	assert.Equal(t, false, p.HasParent())
}

func Test_ReadChildXmlPolicy(t *testing.T) {
	p, err := New("../../../test/fixtures/build/child/1.xml")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	assert.Equal(t, true, p.HasParent())
}

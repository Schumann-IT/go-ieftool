package policy

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ReadTreeFromDir(t *testing.T) {
	expected := []int{1, 2}

	d, _ := filepath.Abs("../../../test/fixtures/build")
	tt := &Tree{}
	err := tt.Read(d)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	bs := tt.Batches()

	assert.Equal(t, len(expected), len(bs))

	for i, b := range bs {
		assert.Equal(t, expected[i], len(b))
	}
}

func Test_ExtractBatchesFromPolicyTree(t *testing.T) {
	expected := []int{1, 3, 2}

	tt := &Tree{
		&XmlPolicy{
			PolicyId:      "r",
			BaseXmlPolicy: []*BaseXmlPolicy{},
			s:             nil,
		},
		&XmlPolicy{
			PolicyId:      "b101",
			BaseXmlPolicy: testCreateBasePolicy("r"),
			s:             nil,
		},
		&XmlPolicy{
			PolicyId:      "b102",
			BaseXmlPolicy: testCreateBasePolicy("r"),
			s:             nil,
		},
		&XmlPolicy{
			PolicyId:      "b103",
			BaseXmlPolicy: testCreateBasePolicy("r"),
			s:             nil,
		},
		&XmlPolicy{
			PolicyId:      "b201",
			BaseXmlPolicy: testCreateBasePolicy("b101"),
			s:             nil,
		},
		&XmlPolicy{
			PolicyId:      "b202",
			BaseXmlPolicy: testCreateBasePolicy("b101"),
			s:             nil,
		},
	}
	bs := tt.Batches()

	assert.Equal(t, len(expected), len(bs))

	for i, b := range bs {
		assert.Equal(t, expected[i], len(b))
	}
}

func testCreateBasePolicy(id string) []*BaseXmlPolicy {
	var r []*BaseXmlPolicy
	b := BaseXmlPolicy{
		PolicyId: []PolicyId{
			{Value: id},
		},
	}
	r = append(r, &b)

	return r
}

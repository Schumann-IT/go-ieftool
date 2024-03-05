package policy

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"com.schumann-it.go-ieftool/pkg/b2c/policy/content"
	"github.com/stretchr/testify/assert"
)

func Test_CreateSourceFromDir(t *testing.T) {
	expected := 1

	b := NewBuilder()

	d, _ := filepath.Abs("../../../test/fixtures/base")
	err := b.Read(d)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	actual := b.Len()

	assert.Equal(t, expected, actual)
}

func Test_Process(t *testing.T) {
	tests := map[string][]string{
		"foo.txt": {"FOO", "foo value"},
		"bar.txt": {"{BAR}", "bar value"},
	}

	for k, v := range tests {
		expected := v[1]
		s := &Builder{
			s: content.Source{
				k: content.Raw(fmt.Sprintf("{Settings:%s}", v[0])),
			},
			p: content.Processed{},
		}
		err := s.Process(map[string]string{
			v[0]: expected,
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		r := s.Result()
		actual := string(r[k])

		assert.Equal(t, expected, actual)
	}
}

func Test_Write(t *testing.T) {
	expected := []byte("content")

	b := Builder{
		p: map[string][]byte{
			"test.txt": expected,
		},
	}
	err := b.Write("/tmp/b2cpolicytest/build")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	actual, err := os.ReadFile("/tmp/b2cpolicytest/build/test.txt")
	if err != nil {
		t.Fatalf("unextected error: %s", err)
	}

	assert.Equal(t, expected, actual)
}

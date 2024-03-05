package converter

import (
	"testing"
)

func TestConvert(t *testing.T) {
	t.Run("convert slice or map", func(t *testing.T) {
		from1 := []string{"test"}
		to1 := &[]string{}
		err := Convert(from1, to1)
		if err != nil {
			t.Fatalf("from value %v should not return an error", from1)
		}
		result1 := *to1
		if result1[0] != from1[0] {
			t.Fatalf("from value %v converts to %v", from1, to1)
		}

		from2 := map[string]string{"foo": "bar"}
		to2 := &map[string]string{}
		err = Convert(from2, to2)
		if err != nil {
			t.Fatalf("empty %v should return an error", from2)
		}
		result2 := *to2
		if result2["foo"] != from2["foo"] {
			t.Fatalf("from value %v converts to %v", from2, to2)
		}
	})

	t.Run("convert empty struct type should not return an error", func(t *testing.T) {
		from1 := struct {
			Foo string `json:"foo,omitempty"`
		}{}
		from2 := struct {
			Foo string
		}{}
		to1 := &map[string]interface{}{}

		err := Convert(from1, to1)
		if err != nil {
			t.Fatalf("empty %v should not return an error", from1)
		}
		if len(*to1) != 0 {
			t.Fatalf("empty %v should produce an empty value: %v", from1, to1)
		}

		err = Convert(from2, to1)
		if err != nil {
			t.Fatalf("empty %v should not return an error", from2)
		}
		if len(*to1) != 1 {
			t.Fatalf("empty %v should produce a non empty value: %v", from2, to1)
		}
	})
}

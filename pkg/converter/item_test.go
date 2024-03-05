package converter

import (
	"testing"
)

func TestItemInitialization(t *testing.T) {
	t.Run("convert type that has length with len() = 0 returns an error", func(t *testing.T) {
		to := &map[string]string{}

		_, err := newItem([]string{}, to)
		if err == nil {
			t.Fatalf("empty %v should return an error", []string{})
		}
		_, err = newItem(map[string]string{}, to)
		if err == nil {
			t.Fatalf("empty %v should return an error", map[string]string{})
		}
	})

	t.Run("convert empty struct type should not return an error", func(t *testing.T) {
		to := &map[string]interface{}{}

		from1 := struct {
			Foo string `json:"foo,omitempty"`
		}{}
		_, err := newItem(from1, to)
		if err != nil {
			t.Fatalf("empty %v should not return an error", from1)
		}

		from2 := struct {
			Foo string
		}{}
		_, err = newItem(from2, to)
		if err != nil {
			t.Fatalf("empty %v should not return an error", from2)
		}
	})
}

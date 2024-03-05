package content

// Source is a type alias for a map that associates string keys with values of type Raw.
type Source map[string]Raw

// Len returns the length of the Source map.
func (c *Source) Len() int {
	return len(*c)
}

// Processed is a type alias for a map that maps strings to slices of bytes.
type Processed map[string][]byte

// Len returns the number of elements in the Processed map.
func (c *Processed) Len() int {
	return len(*c)
}

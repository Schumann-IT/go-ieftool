package converter

type Create[T any] interface {
	From(interface{}) error
}

// Convert uses json.Marshal and json.Unmarshal
// to convert any value to its representation
// defined in the to parameter
func Convert(from, to interface{}) error {
	i, err := newItem(from, to)
	if err != nil {
		return err
	}
	return i.convert()
}

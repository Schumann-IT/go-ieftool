package policy

// Policy is an interface that represents a policy.
type Policy interface {
	Id() string
	HasParent() bool
	Parent() Policy
	File() string
	Byte() []byte
}

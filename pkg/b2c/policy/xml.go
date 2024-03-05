package policy

import (
	"encoding/xml"
	"fmt"
	"os"
)

// New reads the XML file specified by the string f and returns a new Policy object
// using xml.Unmarshal. As a basic check, New ensures that the id of the policy is not empty.
func New(f string) (Policy, error) {
	b, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}
	p := &XmlPolicy{
		s: &f,
		b: b,
	}
	err = xml.Unmarshal(b, p)
	if err != nil {
		return nil, err
	}
	if !p.Valid() {
		return nil, fmt.Errorf("file %s is not a valid policy", f)
	}

	return p, nil
}

// XmlPolicy represents the value of a policy in XML, used to be unmarshalled by xml.Unmarshal
type XmlPolicy struct {
	PolicyId      string           `xml:"PolicyId,attr"`
	BaseXmlPolicy []*BaseXmlPolicy `xml:"BasePolicy"`
	s             *string
	b             []byte
}

// Id @see Policy.Id
func (p *XmlPolicy) Id() string {
	return p.PolicyId
}

// HasParent @see Policy.HasParent
func (p *XmlPolicy) HasParent() bool {
	return len(p.BaseXmlPolicy) > 0
}

// Parent @see Policy.Parent
func (p *XmlPolicy) Parent() Policy {
	return p.BaseXmlPolicy[0]
}

// File @see Policy.File
func (p *XmlPolicy) File() string {
	return *p.s
}

// Byte @see Policy.Byte
func (p *XmlPolicy) Byte() []byte {
	return p.b
}

func (p *XmlPolicy) Valid() bool {
	return p.Id() != ""
}

// BaseXmlPolicy represents the value of a base policy in XML.
type BaseXmlPolicy struct {
	PolicyId []PolicyId `xml:"PolicyId"`
}

// Id @see Policy.Id
func (p *BaseXmlPolicy) Id() string {
	return p.PolicyId[0].Value
}

// HasParent @see Policy.HasParent
func (p *BaseXmlPolicy) HasParent() bool {
	return false
}

// Parent @see Policy.Parent
func (p *BaseXmlPolicy) Parent() Policy {
	return nil
}

// File @see Policy.File
func (p *BaseXmlPolicy) File() string {
	return ""
}

// Byte @see Policy.Byte
func (p *BaseXmlPolicy) Byte() []byte {
	return nil
}

// PolicyId represents the value of a policy ID in XML
type PolicyId struct {
	Value string `xml:",chardata"`
}

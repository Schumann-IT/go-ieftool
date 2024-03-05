package keyset

import (
	"strings"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

type KeySet models.TrustFrameworkKeySet

type RsaKeySet struct {
	KeySet
	k *models.TrustFrameworkKey
}

type CertificateKeySet struct {
	KeySet
	c *Certificate
}

func NewKeySet(id string) *models.TrustFrameworkKeySet {
	ks := models.NewTrustFrameworkKeySet()
	ks.SetId(to.StringPtr(id))

	return ks
}

func NewCertificateKeySet(id, body, password string) *CertificateKeySet {
	ks := NewKeySet(id)
	return &CertificateKeySet{
		KeySet: KeySet(*ks),
		c: &Certificate{
			Body:     body,
			Password: password,
		},
	}
}

func NewRsaKeySet(id, use string) *RsaKeySet {
	k := models.NewTrustFrameworkKey()
	k.SetUse(to.StringPtr(use))
	k.SetKty(to.StringPtr("RSA"))

	ks := NewKeySet(id)
	return &RsaKeySet{
		KeySet: KeySet(*ks),
		k:      k,
	}
}

type KeySets struct {
	IDs []string
}

func New(ids []string) *KeySets {
	return &KeySets{IDs: ids}
}

func (ks *KeySets) Add(s string) {
	ks.IDs = append(ks.IDs, s)
}

func (ks *KeySets) String() string {
	return strings.Join(ks.IDs, ", ")
}

func (ks *KeySets) Remove(s *string) {
	for i, e := range ks.IDs {
		if e == *s {
			ids := append(ks.IDs[:i], ks.IDs[i+1:]...)
			ks.IDs = ids
			break
		}
	}
}

func (ks *KeySets) Has(s string) bool {
	for _, e := range ks.IDs {
		if e == s {
			return true
		}
	}

	return false
}

package b2c

import (
	"fmt"

	"com.schumann-it.go-ieftool/pkg/b2c/keyset"
)

func (s *Service) CreateKeySets(en string, cert *keyset.Certificate) error {
	n := keyset.New([]string{"B2C_1A_TokenSigningKeyContainer", "B2C_1A_TokenEncryptionKeyContainer"})
	if cert != nil {
		n.Add("B2C_1A_SamlIdpCert")
	}

	log.Debugf("processing key sets: %s", n.String())

	c, err := s.createGraphClient(en)
	if err != nil {
		return fmt.Errorf("failed to create graph client: %s", err)
	}

	ks, err := c.GetKeySets()
	for _, k := range ks {
		n.Remove(k.GetId())
	}

	for _, id := range n.IDs {
		k, err := c.CreateKeySet(id)
		if err != nil {
			return err
		}
		switch id {
		case "B2C_1A_TokenSigningKeyContainer":
			err = c.GenerateRsaKeyFor(k, "sig")
			if err != nil {
				return err
			}
			break
		case "B2C_1A_TokenEncryptionKeyContainer":
			err = c.GenerateRsaKeyFor(k, "enc")
			if err != nil {
				return err
			}
			break
		case "B2C_1A_SamlIdpCert":
			err = c.UploadPkcs12For(k, *cert)
			if err != nil {
				return err
			}
			break
		default:
			log.Fatalf("Key Set %s not recognized", id)
		}
	}

	return nil
}

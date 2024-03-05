package msgraph

import (
	"context"

	"com.schumann-it.go-ieftool/pkg/b2c/keyset"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/trustframework"
)

func (s *ServiceClient) GetKeySets() ([]*models.TrustFrameworkKeySet, error) {
	resp, err := s.gc.TrustFramework().KeySets().Get(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	var ks []*models.TrustFrameworkKeySet
	for _, sa := range resp.GetValue() {
		k := models.NewTrustFrameworkKeySet()
		k.SetId(sa.GetId())
		ks = append(ks, k)
	}

	return ks, nil
}

func (s *ServiceClient) DeleteKeySet(ks *models.TrustFrameworkKeySet) error {
	return s.gc.TrustFramework().KeySets().ByTrustFrameworkKeySetId(to.String(ks.GetId())).Delete(context.Background(), nil)
}

func (s *ServiceClient) CreateKeySet(id string) (*models.TrustFrameworkKeySet, error) {
	ks := models.NewTrustFrameworkKeySet()
	ks.SetId(to.StringPtr(id))
	_, err := s.gc.TrustFramework().KeySets().Post(context.Background(), ks, nil)
	if err != nil {
		return nil, err
	}

	return ks, nil
}

func (s *ServiceClient) GenerateRsaKeyFor(ks *models.TrustFrameworkKeySet, use string) error {
	k := trustframework.NewKeySetsItemGenerateKeyPostRequestBody()
	k.SetUse(to.StringPtr(use))
	k.SetKty(to.StringPtr("RSA"))
	_, err := s.gc.TrustFramework().KeySets().ByTrustFrameworkKeySetId(to.String(ks.GetId())).GenerateKey().Post(context.Background(), k, nil)

	return err
}

func (s *ServiceClient) UploadPkcs12For(ks *models.TrustFrameworkKeySet, cert keyset.Certificate) error {
	b := trustframework.NewKeySetsItemUploadPkcs12PostRequestBody()
	b.SetKey(to.StringPtr(cert.Body))
	b.SetPassword(to.StringPtr(cert.Password))
	_, err := s.gc.TrustFramework().KeySets().ByTrustFrameworkKeySetId(to.String(ks.GetId())).UploadPkcs12().Post(context.Background(), b, nil)

	return err
}

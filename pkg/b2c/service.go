package b2c

import (
	"fmt"
	"path/filepath"

	"com.schumann-it.go-ieftool/pkg/b2c/environment"
	"com.schumann-it.go-ieftool/pkg/b2c/msgraph"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

type Service struct {
	es []environment.Config
	sd string
	td string
}

func NewService(cp string, sd string, td string) (*Service, error) {
	c, err := environment.NewConfigFromFile(cp)
	if err != nil {
		return nil, err
	}

	asd, err := filepath.Abs(sd)
	if err != nil {
		return nil, err
	}

	tsd, err := filepath.Abs(td)
	if err != nil {
		return nil, err
	}

	return &Service{
		es: *c,
		sd: asd,
		td: tsd,
	}, nil
}

func (s *Service) findConfig(n string) (*environment.Config, error) {
	for _, e := range s.es {
		if e.Name == n {
			return &e, nil
		}
	}

	return nil, fmt.Errorf("environment %s not found", n)
}

func (s *Service) createGraphClient(en string) (*msgraph.ServiceClient, error) {
	cred, err := azidentity.NewEnvironmentCredential(nil)
	if err != nil {
		return nil, err
	}

	return msgraph.NewClientWithCredential(cred)
}

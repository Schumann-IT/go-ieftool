package b2c

import (
	"fmt"
	"path"
	"path/filepath"

	"com.schumann-it.go-ieftool/pkg/b2c/environment"
	"com.schumann-it.go-ieftool/pkg/b2c/policy"
	"com.schumann-it.go-ieftool/pkg/b2c/secret"
)

type Api struct {
	es []environment.Config
	sd string
	td string
}

type Policies struct {
	e environment.Config
	b [][]policy.Policy
}

func NewApi(cp string, sd string, td string) (*Api, error) {
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

	return &Api{
		es: *c,
		sd: asd,
		td: tsd,
	}, nil
}

func (a *Api) FindConfig(n string) *environment.Config {
	for _, e := range a.es {
		if e.Name == n {
			return &e
		}
	}

	return nil
}

func (a *Api) BuildPolicies(en string) error {
	e := a.FindConfig(en)
	if e == nil {
		return fmt.Errorf("environment %s not found", en)
	}

	b := policy.NewBuilder()
	err := b.Read(a.sd)
	if err != nil {
		return fmt.Errorf("failed to read from %s/%s: %s", a.sd, e.Name, err)
	}
	err = b.Process(e.Settings)
	if err != nil {
		return fmt.Errorf("failed to process %s/%s: %s", a.sd, e.Name, err)
	}
	err = b.Write(path.Join(a.td, e.Name))
	if err != nil {
		return fmt.Errorf("failed to write to %s/%s: %s", a.td, e.Name, err)
	}

	return nil
}

func (a *Api) Batch(en string) (*Policies, error) {
	e := a.FindConfig(en)
	if e == nil {
		return nil, fmt.Errorf("environment %s not found", en)
	}

	sd := path.Join(a.td, e.Name)
	t := &policy.Tree{}
	err := t.Read(sd)
	if err != nil {
		if err != nil {
			return nil, fmt.Errorf("failed to read from %s/%s, did you run build?: %s", a.td, e.Name, err)
		}
	}

	return &Policies{
		b: t.Batches(),
		e: *e,
	}, nil
}

func DeletePolicies(e environment.Config) {}

func (a *Api) ListPolicies(e environment.Config) {

}

func FixAppRegistration(e environment.Config) {}

func CreateKeySets(e environment.Config) {}

func DeleteKeySets(e environment.Config) {}

type GraphApiCred struct {
	TenantId *string `json:"tenant_id,omitempty"`
	ClientId string  `json:"client_id"`
	SecretId string  `json:"secret_id"`
}

type SamlCert struct {
	Cert         string `json:"cert"`
	CertPassword string `json:"cert_password"`
}

func SecretProviderFromEnvironment(c environment.Config) secret.Provider {
	sp := secret.ProviderRegistry[c.SecretProviderConfig.Name]
	o := c.SecretProviderConfig.Options
	err := sp.Setup(o)
	if err != nil {
		log.Fatalf("failed to setup secret provider: %v", err)
	}
	return sp
}

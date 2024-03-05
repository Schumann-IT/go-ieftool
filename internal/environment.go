package internal

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"com.schumann-it.go-ieftool/internal/msgraph"
	"com.schumann-it.go-ieftool/internal/msgraph/trustframework"
	"com.schumann-it.go-ieftool/internal/vault"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

//go:embed msgraph/trustframework/ApplicationPatchIdentityFramework.json
var iefApplicationPatch []byte

//go:embed msgraph/trustframework/ApplicationPatchSaml.json
var samlApplicationPatch []byte

type EnvironmentSaml struct {
	AppObjectId *string `yaml:"appObjectId,omitempty"`
	MetadataUrl *string `yaml:"metadataUrl,omitempty"`
}

type Environment struct {
	Name                                   string                 `yaml:"name"`
	SecretPath                             string                 `yaml:"secretPath"`
	IsProduction                           bool                   `yaml:"isProduction"`
	Tenant                                 string                 `yaml:"tenant"`
	TenantId                               string                 `yaml:"tenantId"`
	IdentityExperienceFrameworkAppObjectId *string                `yaml:"identityExperienceFrameworkAppObjectId,omitempty"`
	Saml                                   *EnvironmentSaml       `yaml:"saml,omitempty"`
	Settings                               map[string]interface{} `yaml:"settings"`
	Secret                                 *vault.Secret
	GraphClient                            *msgraph.Client
}

func (env Environment) Build(s string, d string) error {
	var errs Errors
	root := s
	err := filepath.WalkDir(s, func(p string, e fs.DirEntry, err error) error {
		log.Debugf("processing dir: %s", s)
		if err != nil {
			return err
		}
		if e.IsDir() {
			return nil
		}
		if filepath.Ext(e.Name()) == ".xml" {
			t := path.Join(d, strings.ReplaceAll(p, root, env.Name))
			c, ve := env.replaceVariables(p)
			if ve != nil {
				errs = append(errs, ve)
				return nil
			}
			ve = os.MkdirAll(filepath.Dir(t), os.ModePerm)
			if ve != nil {
				errs = append(errs, ve)
				return nil
			}
			log.Debugf("compiled: %s", t)
			if env.IsProduction {
				// @TODO remove debug code
				log.Debug("removed debug parameters as this is a prod environment.")
			}
			ve = os.WriteFile(t, c, os.ModePerm)
			if ve != nil {
				errs = append(errs, ve)
				return nil
			}
			log.Debugf("processed file: %s", t)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if errs.HasErrors() {
		return errs.Format()
	}

	return nil
}

func (env Environment) replaceVariables(p string) ([]byte, error) {
	log.Debugf("replacing variables for: %s", p)

	content, err := os.ReadFile(p)
	policy := string(content)
	if err != nil {
		return nil, err
	}

	var errs Errors
	for _, v := range env.variables(policy) {
		val, ve := env.value(v)
		if ve != nil {
			errs = append(errs, fmt.Errorf("%s: but required in source %s", ve.Error(), p))
			continue
		}
		if val == "" || val == "null" {
			errs = append(errs, fmt.Errorf("variable '%s' must not be empty: source %s", v, p))
			continue
		}
		re := regexp.MustCompile(fmt.Sprintf("{Settings:%s}", v))
		policy = re.ReplaceAllString(policy, val)
	}

	if errs.HasErrors() {
		return nil, errs.Format()
	}

	return []byte(policy), nil
}

func (env Environment) variables(c string) []string {
	re := regexp.MustCompile(`{Settings:(.+)}`)
	m := re.FindAllStringSubmatch(c, -1)
	var cm []string
	seen := make(map[string]bool, len(m))
	for _, match := range m {
		if !seen[match[1]] {
			cm = append(cm, match[1])
			seen[match[1]] = true
		}
	}

	return cm
}

func (env Environment) value(n string) (string, error) {
	switch n {
	case "Tenant":
		return env.Tenant, nil
	default:
		if env.Settings[n] == nil {
			return "", fmt.Errorf("variable '%s' is not provided in settings", n)
		}

		return env.Settings[n].(string), nil
	}
}

func (env Environment) Deploy(d string) error {
	ps, err := trustframework.NewPoliciesFromDir(path.Join(d, env.Name))
	if err != nil {
		return err
	}
	bs := ps.GetBatch()

	for i, b := range bs {
		log.Debugf("processing batch %d", i)
		env.GraphClient.UploadPolicies(b)
	}

	return nil
}

func (env Environment) ListRemotePolicies() ([]string, error) {
	return env.GraphClient.ListPolicies()
}

func (env Environment) DeleteRemotePolicies() error {
	return env.GraphClient.DeletePolicies()
}

func (env Environment) FixAppRegistrations() error {
	if env.IdentityExperienceFrameworkAppObjectId == nil {
		return fmt.Errorf("please specify identityExperienceFrameworkObjectId in envirnment")
	}
	err := env.GraphClient.FixAppRegistration(*env.IdentityExperienceFrameworkAppObjectId, iefApplicationPatch)
	if err != nil {
		log.Fatal(err)
	}

	if env.Saml != nil && env.Saml.AppObjectId != nil {
		var p map[string]interface{}
		err = json.Unmarshal(samlApplicationPatch, &p)
		if env.Saml.MetadataUrl != nil {
			p["samlMetadataUrl"] = env.Saml.MetadataUrl
		}
		patch, err := json.Marshal(p)
		if err != nil {
			log.Fatal(err)
		}
		err = env.GraphClient.FixAppRegistration(*env.Saml.AppObjectId, patch)
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func (env Environment) CreateKeySets() error {
	return env.GraphClient.CreateKeySets(env.Secret)
}

func (env Environment) DeleteKeySets() interface{} {
	return env.GraphClient.DeleteKeySets()
}

type Environments struct {
	e []Environment
	s string
	d string
}

func MustNewEnvironmentsFromFlags(f *pflag.FlagSet) *Environments {
	cf, err := f.GetString("config")
	if err != nil {
		log.Fatalf("could not parse flag 'config': \n%s", err.Error())
	}

	en, err := f.GetString("environment")
	if err != nil {
		log.Fatalf("could not parse flag 'environment': \n%s", err.Error())
	}

	e, err := NewEnvironmentsFromConfig(cf, en)
	if err != nil {
		log.Fatalf("could not read environments config: \n%s", err.Error())
	}

	return e
}

func NewEnvironmentsFromConfig(p string, n string) (*Environments, error) {
	var e []Environment

	c, err := os.ReadFile(p)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not read %s: %s", p, err.Error()))
	}

	err = yaml.Unmarshal(c, &e)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not unmarshal config from %s: %s", p, err.Error()))
	}

	es := Environments{
		e: e,
	}
	es.e = e
	es.filter(n)

	return &es, nil
}

func (es *Environments) Len() int {
	return len(es.e)
}

func (es *Environments) String() string {
	var s []string
	for _, e := range es.e {
		s = append(s, e.Name)
	}
	return strings.Join(s, ", ")
}

func (es *Environments) Build(s string, d string) error {
	es.s = s
	es.d = d

	for _, e := range es.e {
		err := e.Build(es.s, es.d)
		if err != nil {
			return err
		}
	}

	return nil
}

func (es *Environments) FetchSecrets() error {
	vc := vault.NewClient()

	for i, _ := range es.e {
		if es.e[i].Secret == nil {
			sp := fmt.Sprintf("%s/%s", es.e[i].SecretPath, es.e[i].Name)
			s, err := vc.GetSecret(sp)
			if err != nil {
				return errors.New(fmt.Sprintf("could not find secret %s: %s", sp, err.Error()))
			}
			es.e[i].Secret = s
		}
	}

	return nil
}

func (es *Environments) CreateGraphClients() error {
	err := es.FetchSecrets()
	if err != nil {
		return err
	}

	for i, _ := range es.e {
		if es.e[i].GraphClient == nil {
			c, err := msgraph.NewClient(es.e[i].TenantId, es.e[i].Secret.ClientId, es.e[i].Secret.ClientSecret)
			if err != nil {
				return fmt.Errorf("could not create graph client: %s", err.Error())
			}
			es.e[i].GraphClient = c
		}
	}

	return nil
}

func (es *Environments) Deploy(d string) error {
	err := es.CreateGraphClients()
	if err != nil {
		return err
	}

	es.d = d
	for _, e := range es.e {
		err := e.Deploy(es.d)
		if err != nil {
			return err
		}
	}

	return nil
}

func (es *Environments) filter(n string) {
	var ne []Environment

	for _, e := range es.e {
		if n == "" || n == e.Name {
			ne = append(ne, e)
		}
	}

	es.e = ne
}

func (es *Environments) FetchRemotePolicies() (map[string][]string, error) {
	err := es.CreateGraphClients()
	if err != nil {
		return nil, err
	}

	var errs Errors
	r := map[string][]string{}
	for _, e := range es.e {
		ps, err := e.ListRemotePolicies()
		if err != nil {
			errs = append(errs, err)
		}
		r[e.Name] = ps
	}

	if errs.HasErrors() {
		return nil, errs.Format()
	}

	return r, nil
}

func (es *Environments) DeleteRemotePolicies() error {
	err := es.CreateGraphClients()
	if err != nil {
		return err
	}

	var errs Errors
	for _, e := range es.e {
		err := e.DeleteRemotePolicies()
		if err != nil {
			errs = append(errs, errors.New(fmt.Sprintf("Failed to delete policy from environment %s: %s", e.Name, err)))
		}
	}

	if errs.HasErrors() {
		return errs.Format()
	}

	return nil
}

func (es *Environments) FixAppRegistrations() error {
	err := es.CreateGraphClients()
	if err != nil {
		return err
	}

	var errs Errors
	for _, e := range es.e {
		err := e.FixAppRegistrations()
		if err != nil {
			errs = append(errs, errors.New(fmt.Sprintf("Failed to fix app registrations from environment %s: %s", e.Name, err)))
		}
	}

	if errs.HasErrors() {
		return errs.Format()
	}

	return nil
}

func (es *Environments) CreateKeySets() error {
	err := es.CreateGraphClients()
	if err != nil {
		return err
	}

	var errs Errors
	for _, e := range es.e {
		err := e.CreateKeySets()
		if err != nil {
			errs = append(errs, errors.New(fmt.Sprintf("Failed to key sets from environment %s: %s", e.Name, err)))
		}
	}

	if errs.HasErrors() {
		return errs.Format()
	}

	return nil
}

func (es *Environments) DeleteKeySets() error {
	err := es.CreateGraphClients()
	if err != nil {
		return err
	}

	var errs Errors
	for _, e := range es.e {
		err := e.DeleteKeySets()
		if err != nil {
			errs = append(errs, errors.New(fmt.Sprintf("Failed to key sets from environment %s: %s", e.Name, err)))
		}
	}

	if errs.HasErrors() {
		return errs.Format()
	}

	return nil
}

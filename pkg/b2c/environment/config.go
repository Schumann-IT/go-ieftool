package environment

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Name                 string                `yaml:"name"`
	Tenant               string                `yaml:"tenant"`
	IsProduction         bool                  `yaml:"isProduction"`
	Applications         []ApplicationConfig   `yaml:"applications,omitempty"`
	SecretProviderConfig *SecretProviderConfig `yaml:"secretProviderConfig,omitempty"`
	Settings             map[string]string     `yaml:"settings,omitempty"`
}

type ApplicationConfig struct {
	Name     string                 `yaml:"name"`
	ObjectId string                 `yaml:"objectId"`
	Patch    map[string]interface{} `yaml:"patch,omitempty"`
}

type SecretProviderConfig struct {
	Name    string                  `yaml:"name"`
	Options *map[string]interface{} `yaml:"options,omitempty"`
}

func NewConfig(b []byte) (*Config, error) {
	c := &Config{}
	err := yaml.Unmarshal(b, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func NewConfigFromFile(f string) (*[]Config, error) {
	c := &[]Config{}
	b, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(b, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

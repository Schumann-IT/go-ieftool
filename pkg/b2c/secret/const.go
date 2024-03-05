package secret

import "gopkg.in/op/go-logging.v1"

var (
	ProviderRegistry = map[string]Provider{
		"env":   &EnvProvider{},
		"vault": &VaultProvider{},
	}

	log = logging.MustGetLogger("b2c/secret")
)

package secret

import (
	"context"
	"fmt"
	"strings"
	"time"

	"com.schumann-it.go-ieftool/pkg/converter"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

type VaultProviderOptions struct {
	Address   *string `json:"address,omitempty"`
	RoleId    string  `json:"roleId"`
	SecretId  string  `json:"secretId"`
	MountPath string  `json:"mountPath"`
}

func (o *VaultProviderOptions) From(i *map[string]interface{}) error {
	return converter.Convert(i, o)
}

type VaultProvider struct {
	s *vault.Secrets
	m string
}

func (p *VaultProvider) Setup(o *map[string]interface{}) error {
	var err error
	if o == nil {
		// read options from environment
		ep := ProviderRegistry["env"]
		op, err := ep.Secret(Request{
			Data: map[string]interface{}{
				"vault_addr":              "address",
				"vault_role_id":           "roleId",
				"vault_secret_id":         "secretId",
				"vault_secret_mount_path": "mountPath",
			},
		})
		if err != nil {
			return err
		}
		o = &op.r
	}

	po := &VaultProviderOptions{}
	err = po.From(o)
	if err != nil {
		return err
	}

	c, err := vault.New(
		vault.WithAddress(to.String(po.Address)),
		vault.WithRequestTimeout(30*time.Second),
	)
	if err != nil {
		return err
	}

	s, err := c.Auth.AppRoleLogin(
		context.Background(),
		schema.AppRoleLoginRequest{
			RoleId:   po.RoleId,
			SecretId: po.SecretId,
		},
	)
	if err != nil {
		return err
	}

	err = c.SetToken(s.Auth.ClientToken)
	if err != nil {
		return err
	}

	p.s = &c.Secrets
	p.m = po.MountPath

	return nil
}

func (p *VaultProvider) Secret(r Request) (*Response, error) {
	var parts []string
	path, ok := r.Data["path"]
	if !ok {
		return nil, fmt.Errorf("secret request must contain path parameter")
	}
	parts = append(parts, strings.ToLower(path.(string)))
	if r.Name != nil {
		parts = append(parts, strings.ToLower(to.String(r.Name)))
	}
	if r.Environment != nil {
		parts = append(parts, strings.ToLower(to.String(r.Environment)))
	}

	d, err := p.s.KvV2Read(context.Background(), strings.Join(parts, "/"), vault.WithMountPath(p.m))
	if err != nil {
		return nil, err
	}

	return &Response{
		r: d.Data.Data,
	}, nil
}

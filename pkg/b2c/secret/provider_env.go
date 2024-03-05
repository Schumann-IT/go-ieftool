package secret

import (
	"fmt"
	"os"
	"strings"

	"github.com/Azure/go-autorest/autorest/to"
)

type EnvProvider struct{}

func (p *EnvProvider) Setup(_ *map[string]interface{}) error {
	return nil
}

func (p *EnvProvider) Secret(r Request) (*Response, error) {
	d := map[string]interface{}{}

	var parts []string
	if r.Name != nil {
		parts = append(parts, strings.ToUpper(to.String(r.Name)))
	}
	if r.Environment != nil {
		parts = append(parts, strings.ToUpper(to.String(r.Environment)))
	}
	prefix := strings.Join(parts, "_")

	for k, tk := range r.Data {
		vn := fmt.Sprintf("%s_%s", prefix, strings.ToUpper(k))
		v := os.Getenv(vn)
		if v == "" {
			return nil, fmt.Errorf("please set env var %s", vn)
		}
		d[tk.(string)] = v
	}

	return &Response{
		r: d,
	}, nil
}

package secret

import "com.schumann-it.go-ieftool/pkg/converter"

type Provider interface {
	Setup(*map[string]interface{}) error
	Secret(Request) (*Response, error)
}

type Request struct {
	Name        *string
	Environment *string
	Data        map[string]interface{}
}

type Response struct {
	r map[string]interface{}
}

func (r Response) Convert(to interface{}) error {
	return converter.Convert(r.r, to)
}

func (r Response) Get(f string) string {
	return r.r[f].(string)
}
